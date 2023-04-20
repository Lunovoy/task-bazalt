package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type BranchBinaryPackages struct {
	RequestArgs struct {
		Arch any `json:"arch"`
	} `json:"request_args"`
	Length   int       `json:"length"`
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
		// fmt.Println(string(body))
		var branchBinaryPackages BranchBinaryPackages
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
	greaterVersions := map[string]string{}

	for _, pkg1 := range packages1 {
		for _, pkg2 := range packages2 {
			if pkg1.Name == pkg2.Name {
				v1 := fmt.Sprintf("%s-%s", pkg1.Version, pkg1.Release)
				v2 := fmt.Sprintf("%s-%s", pkg2.Version, pkg2.Release)

				if compareVersions(v1, v2) == 1 {
					greaterVersions[pkg1.Name] = v1
				}
			}
		}
	}

	for name, version := range greaterVersions {
		fmt.Printf("%s-%s\n", name, version)
	}
}

// сравнивает версии двух пакетов в формате 'version-release' и возвращает
//
//	1, если версия первого пакета больше
//	0, если версии равны
//	-1, если версия первого пакета меньше
func compareVersions(v1 string, v2 string) int {
	parts1 := strings.Split(v1, "-")
	parts2 := strings.Split(v2, "-")

	// сравниваем версии
	for i := 0; i < len(parts1)-1 && i < len(parts2)-1; i++ {
		if parts1[i] > parts2[i] {
			return 1
		} else if parts1[i] < parts2[i] {
			return -1
		}
	}

	// если версии равны, то сравниваем releases
	if parts1[len(parts1)-1] > parts2[len(parts2)-1] {
		return 1
	} else if parts1[len(parts1)-1] < parts2[len(parts2)-1] {
		return -1
	}

	return 0
}

func showMenu() {
	fmt.Println("1) Lists of binary packages of 2 branches")
	fmt.Println("2) All packages that are in the 1st but not in the 2nd")
	fmt.Println("3) All packages that are in the 2nd but not in the 1st")
	fmt.Println("4) All packages whose version-release is greater in the 1st than in the 2nd")
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <branch1> <branch2>\n", os.Args[0])
		os.Exit(1)
	}

	archs := [7]string{"aarch64", "armh", "i586", "noarch", "ppc64le", "x86_64", "x86_64-i586"}

	branch1 := os.Args[1]
	branch2 := os.Args[2]

	// branch1 := "p9"
	// branch2 := "p10"

	// получаем списки пакетов
	packageListsFirstArg := map[string][]Package{}
	packageListsSecondArg := map[string][]Package{}
	for _, arch := range archs {
		packages1, err := getPackages(branch1, arch)
		if err != nil {
			log.Fatalf("Error getting packages for %s on %s: %s", branch1, arch, err)
		}

		packageListsFirstArg[arch+"_"+branch1] = packages1

		packages2, err := getPackages(branch2, arch)
		if err != nil {
			log.Fatalf("Error getting packages for %s on %s: %s", branch2, arch, err)
		}

		packageListsSecondArg[arch+"_"+branch2] = packages2
	}

	//выводим результаты сравнения
	for _, arch := range archs {

		fmt.Printf("Packages only present in %s on %s : ", branch1, arch)
		printDiffPackages(packageListsFirstArg["i586_"+branch1], packageListsSecondArg["i586_"+branch2])
		fmt.Println()
	}

	// fmt.Println("Packages with greater versions in", branch1, ":")
	// printGreaterVersions(packageLists["x86_64_"+branch1], packageLists["x86_64_"+branch2])
	// fmt.Println()

	// fmt.Println("All packages in", branch1, ":")
	// printAllPackages(packageLists["aarch64_"+branch1])

}
