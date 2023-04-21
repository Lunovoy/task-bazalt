package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// получает список пакетов для заданной ветки и архитектуры

//	for name, version := range greaterVersions {
//		fmt.Printf("%s-%s\n", name, version)
//	}
//

func chooseArch(availableArchs []string, clearConsoleCommand []byte) []string {

	archs := make([]string, 0, len(availableArchs))
	selectedArchs := map[string]bool{}
	var choice string

	for choice != "9" {
		fmt.Println("============")
		fmt.Println(string(clearConsoleCommand))
		if len(selectedArchs) != 0 {
			fmt.Print("Selected architectures: ")
			fmt.Print("[")
			for arch := range selectedArchs {
				fmt.Printf("%v ", arch)

			}
			fmt.Println("]")
			fmt.Println("Add architectures:")
		} else {
			fmt.Println("Select architectures:")
		}

		fmt.Println("1) All")
		for i, arch := range availableArchs {
			fmt.Printf("%v) %s\n", i+2, arch)
		}
		fmt.Println("9) Continue")
		fmt.Println("10) Exit")
		fmt.Println("============")
		fmt.Print("Enter number [1-10]: ")
		fmt.Scan(&choice)

		switch choice {
		case "1":
			return availableArchs
		case "2":
			selectedArchs[availableArchs[0]] = true
		case "3":
			selectedArchs[availableArchs[1]] = true
		case "4":
			selectedArchs[availableArchs[2]] = true
		case "5":
			selectedArchs[availableArchs[3]] = true
		case "6":
			selectedArchs[availableArchs[4]] = true
		case "7":
			selectedArchs[availableArchs[5]] = true
		case "8":
			selectedArchs[availableArchs[6]] = true
		case "9":
			for key := range selectedArchs {
				archs = append(archs, key)
			}
			return archs
		case "10":
			os.Exit(0)
		default:
			fmt.Println("Invalid input")
		}
	}
	return nil
}

func showMenu(branch1, branch2 string) {

	fmt.Println("1) Lists of binary packages of 2 branches")
	fmt.Printf("2) All packages that are in the <%s> but not in the <%s>\n", branch1, branch2)
	fmt.Printf("3) All packages that are in the <%s> but not in the <%s>\n", branch2, branch1)
	fmt.Printf("4) All packages whose version-release is greater in the <%s> than in the <%s>\n", branch1, branch2)
}

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
	packageListsBranch1 := map[string][]Package{}
	packageListsBranch2 := map[string][]Package{}

	t := time.Now()
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup
	wg.Add(1)
	go func() {

		for _, arch := range archs {

			wg2.Add(1)
			go func(branch1, arch string) {
				defer wg2.Done()
				packages1, err := getPackages(branch1, arch)
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
				packages2, err := getPackages(branch2, arch)
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
	printDiffPackages(packageListsBranch1[archs[0]+"_"+branch1], packageListsBranch2[archs[0]+"_"+branch2])
	// 	fmt.Println()
	// }

	fmt.Println("Packages with greater versions in", branch1, ":")
	printGreaterVersions(packageListsBranch1["x86_64_"+branch1], packageListsBranch2["x86_64_"+branch2])
	fmt.Println()

	// fmt.Println("All packages in", branch1, ":")
	// fmt.Printf("HuHHIH: %v", packageListsBranch1["aarch64_p10"][0].Name)
	// printAllPackages(packageListsBranch1["aarch64_p10"])

}
