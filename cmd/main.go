package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type BranchBinaryPackages struct {
	RequestArgs struct {
		Arch any `json:"arch"`
	} `json:"request_args"`
	Length   int       `json:"length"`
	Packages []Package `json:"packages"`
}

type TotalBranchPackagesResults struct {
	UniqueInFirst         []BranchPackages `json:"uniqueInFirst"`
	UniqueInSecond        []BranchPackages `json:"uniqueInSecond"`
	GreaterPackageVersion []BranchPackages `json:"greaterVersionsInFirst"`
}

type BranchPackages struct {
	Arch     string    `json:"arch"`
	Packages []Package `json:"packages"`
}

type Package struct {
	Name      string `json:"name"`
	Epoch     int    `json:"epoch"`
	Version   string `json:"version"`
	Release   string `json:"release"`
	Arch      string `json:"arch"`
	Disttag   string `json:"disttag"`
	Buildtime int    `json:"buildtime"`
	Source    string `json:"source"`
}

// получает список пакетов для заданной ветки и архитектуры
func getPackages(branch, arch string) ([]Package, error) {
	url := fmt.Sprintf("https://rdb.altlinux.org/api/export/branch_binary_packages/%s?arch=%s", branch, arch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "AltLinux Package BranchPackages CLI")

	client := &http.Client{}

	respChan := make(chan *http.Response)
	errChan := make(chan error)

	go func() {
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}

		respChan <- resp
	}()

	select {
	case resp := <-respChan:
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		branchBinaryPackages := BranchBinaryPackages{}
		if err := json.Unmarshal(body, &branchBinaryPackages); err != nil {
			return nil, err
		}

		return branchBinaryPackages.Packages, nil

	case err := <-errChan:
		return nil, err

	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("timed out while getting packages for %s on ", branch)
	}
}

// выводит список всех пакетов
func printAllPackages(packages []Package) {
	for _, pkg := range packages {
		fmt.Printf("%s-%s-%s.%s\n", pkg.Name, pkg.Version, pkg.Release, pkg.Arch)
	}
}

// выводит список пакетов, которые присутствуют в первом списке, но отсутствуют во втором списке
func printDiffPackages(packages1 []Package, packages2 []Package) {
	diff := map[string]bool{}

	for _, pkg := range packages1 {
		diff[pkg.Name] = true
	}

	for _, pkg := range packages2 {
		delete(diff, pkg.Name)
	}

	for name := range diff {
		fmt.Printf("%s\n", name)
	}
}

// выводит список пакетов, version-release которых больше в первом списке, чем во втором списке
func printGreaterVersions(packages1 []Package, packages2 []Package) {
	// greaterVersions := map[Package]bool{}
	packagesWithGreaterVersions := make([]Package, 0, len(packages1))
	for _, pkg1 := range packages1 {
		for _, pkg2 := range packages2 {
			if pkg1.Name == pkg2.Name {
				v1 := fmt.Sprintf("%s-%s", pkg1.Version, pkg1.Release)
				v2 := fmt.Sprintf("%s-%s", pkg2.Version, pkg2.Release)

				if compareVersions(v1, v2) == 1 {
					packagesWithGreaterVersions = append(packagesWithGreaterVersions, pkg1)
				}
			}
		}
	}

	BranchPackagesVersions := BranchPackages{
		Arch:     packagesWithGreaterVersions[0].Arch,
		Packages: packagesWithGreaterVersions,
	}

	jsonEncoded, err := json.MarshalIndent(BranchPackagesVersions, "", " ")
	if err != nil {

	}
	fmt.Println(string(jsonEncoded))

}

//	for name, version := range greaterVersions {
//		fmt.Printf("%s-%s\n", name, version)
//	}
//
// сравнивает версии двух пакетов в формате 'version-release' и возвращает
//
//	1, если версия первого пакета больше
//	0, если версии равны
//	-1, если версия первого пакета меньше
func compareVersions(v1 string, v2 string) int {
	parts1 := strings.Split(v1, "-")
	parts2 := strings.Split(v2, "-")

	// сравниваем версии

	if parts1[0] > parts2[0] {
		return 1
	} else if parts1[0] < parts2[0] {
		return -1
	}

	// если версии равны, то сравниваем releases
	if parts1[1] > parts2[1] {
		return 1
	} else if parts1[1] < parts2[1] {
		return -1
	}

	return 0
}

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
