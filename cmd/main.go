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

	archs := []string{"aarch64", "armh", "i586", "noarch", "ppc64le", "x86_64", "x86_64-i586"}

	// cmd := exec.Command("clear")
	// clearConsoleCommand, _ := cmd.Output()

	// selectedArchs := chooseArch(archs, clearConsoleCommand)
	// fmt.Println(string(clearConsoleCommand))
	// fmt.Printf("Selected archs <%+v> f	// archs := []string{"aarch64", "noarch", "i586", "x86_64", "x86_64-i586"}or branches: <%s>, <%s>\n", selectedArchs, branch1, branch2)

	// showMenu(branch1, branch2)

	// branch1 := "p9"
	// branch2 := "p10"

	// получаем списки пакетов

	t := time.Now()
	packageListsBranch1 := map[string][]models.Package{}
	packageListsBranch2 := map[string][]models.Package{}

	var wg sync.WaitGroup
	branches := [2]string{os.Args[1], os.Args[2]}

	for i, branch := range branches {

		for _, arch := range archs {
			time.Sleep(300 * time.Millisecond)
			wg.Add(1)
			go func(branch, arch string, i int) {
				defer wg.Done()
				packages, err := comparison.FetchPackages(branch, arch)
				if err != nil {
					log.Fatalf("Error getting packages for %s on %s: %s", branch, arch, err.Error())
				}
				if i == 0 {
					packageListsBranch1[arch+"_"+branch] = packages
				} else {
					packageListsBranch2[arch+"_"+branch] = packages
				}
				fmt.Printf("%s: %v\n", branch, packages[0].Arch)
			}(branch, arch, i)

		}
	}
	wg.Wait()
	fmt.Println("Success")

	fmt.Println(time.Now().Sub(t).Seconds())

	// //выводим результаты сравнения
	// for _, arch := range archs {

	// fmt.Printf("Packages only present in %s on %s : ", branch1, archs[0])
	// comparison.PrintDiffPackages(packageListsBranch1[archs[0]+"_"+branch1], packageListsBranch2[archs[0]+"_"+branch2])
	// // 	fmt.Println()
	// // }

	// fmt.Println("Packages with greater versions in", branch1, ":")
	comparison.GetGreaterPackagesVersions(packageListsBranch1["x86_64_"+branches[0]], packageListsBranch2["x86_64_"+branches[1]])
	fmt.Println()

	// fmt.Println("All packages in", branch1, ":")
	// fmt.Printf("HuHHIH: %v", packageListsBranch1["aarch64_p10"][0].Name)
	// printAllPackages(packageListsBranch1["aarch64_p10"])

}
