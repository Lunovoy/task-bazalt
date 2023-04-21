package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Lunovoy/task-bazalt/pkg/comparison"
	"github.com/Lunovoy/task-bazalt/pkg/models"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <branch1> <branch2>\n", os.Args[0])
		os.Exit(1)
	}

	// archs := []string{"aarch64", "armh", "i586", "noarch", "ppc64le", "x86_64", "x86_64-i586"}
	archs := []string{"aarch64", "noarch", "i586", "x86_64", "x86_64-i586"}
	branch1 := os.Args[1]
	branch2 := os.Args[2]

	// cmd := exec.Command("clear")
	// clearConsoleCommand, _ := cmd.Output()

	// selectedArchs := chooseArch(archs, clearConsoleCommand)
	// fmt.Println(string(clearConsoleCommand))
	// fmt.Printf("Selected archs <%+v> for branches: <%s>, <%s>\n", selectedArchs, branch1, branch2)

	// showMenu(branch1, branch2)

	// branch1 := "p9"
	// branch2 := "p10"

	// получаем списки пакетов

	t := time.Now()
	packageListsBranch1 := map[string][]models.Package{}
	packageListsBranch2 := map[string][]models.Package{}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		var wg2 sync.WaitGroup
		for _, arch := range archs {

			wg2.Add(1)
			go func(branch1, arch string) {
				defer wg2.Done()
				packages1, err := comparison.FetchPackages(branch1, arch)
				if err != nil {
					log.Fatalf("Error getting packages for %s on %s: %s", branch1, arch, err.Error())
				}
				// fmt.Println("Dlina: ", len(packages1))
				packageListsBranch1[arch+"_"+branch1] = packages1
				fmt.Printf("%s: %v\n", branch1, packages1[0].Arch)
			}(branch1, arch)

			wg2.Add(1)
			go func(branch2, arch string) {
				defer wg2.Done()
				packages2, err := comparison.FetchPackages(branch2, arch)
				if err != nil {
					log.Fatalf("Error getting packages for %s on %s: %s", branch2, arch, err.Error())
				}
				packageListsBranch2[arch+"_"+branch2] = packages2

				fmt.Printf("%s: %v\n", branch2, packages2[0].Arch)
			}(branch2, arch)
		}
		wg2.Wait()
		fmt.Println("Succes")
		defer wg.Done()
	}()
	wg.Wait()

	fmt.Println(time.Now().Sub(t).Seconds())

	// //выводим результаты сравнения
	// for _, arch := range archs {

	fmt.Printf("Packages only present in %s on %s : ", branch1, archs[0])
	comparison.PrintDiffPackages(packageListsBranch1[archs[0]+"_"+branch1], packageListsBranch2[archs[0]+"_"+branch2])
	// 	fmt.Println()
	// }

	fmt.Println("Packages with greater versions in", branch1, ":")
	comparison.PrintGreaterVersions(packageListsBranch1["x86_64_"+branch1], packageListsBranch2["x86_64_"+branch2])
	fmt.Println()

	// fmt.Println("All packages in", branch1, ":")
	// fmt.Printf("HuHHIH: %v", packageListsBranch1["aarch64_p10"][0].Name)
	// printAllPackages(packageListsBranch1["aarch64_p10"])

}
