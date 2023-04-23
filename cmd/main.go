package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Lunovoy/task-bazalt/pkg/comparison"
	"github.com/Lunovoy/task-bazalt/pkg/models"
)

func CreateFile(filename, data string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("error occured while creating file: %s", err.Error())
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintln(data))
	if err != nil {
		log.Printf("error occured while writing to file: %s", err.Error())
		return
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <branch1> <branch2>\n", os.Args[0])
		os.Exit(1)
	}

	branches := [2]string{os.Args[1], os.Args[2]}
	arches := []string{"aarch64", "armh", "i586", "noarch", "ppc64le", "x86_64", "x86_64-i586"}

	// cmd := exec.Command("clear")
	// clearConsoleCommand, _ := cmd.Output()

	// selectedArches := chooseArch(archs, clearConsoleCommand)
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
	fmt.Println("Fetching packages from server")

	tokens := make(chan struct{}, 8)
	for i, branch := range branches {

		for _, arch := range arches {
			wg.Add(1)
			go func(branch, arch string, i int) {
				defer wg.Done()
				tokens <- struct{}{}
				packages, err := comparison.FetchPackages(branch, arch)
				<-tokens
				if err != nil {
					log.Fatalf("Error getting packages for %s on %s: %s", branch, arch, err.Error())
				}
				if i == 0 {
					packageListsBranch1[branch+"_"+arch] = packages
				} else {
					packageListsBranch2[branch+"_"+arch] = packages
				}
				fmt.Printf("%s: %v\n", branch, packages[0].Arch)
			}(branch, arch, i)

		}
	}
	wg.Wait()

	fmt.Printf("Fetched packages in %v\n", time.Now().Sub(t).Seconds())

	size := len(packageListsBranch1)
	uniqueInFirst := make([]models.BranchPackages, 0, size)
	uniqueInSecond := make([]models.BranchPackages, 0, size)
	greaterPackagesVersions := make([]models.BranchPackages, 0, size)

	t2 := time.Now()
	fmt.Println("Processing packages")

	for _, arch := range arches {

		wg.Add(1)
		go func(arch string, branches [2]string) {
			defer wg.Done()
			fmt.Printf("Getting packages only present in %s on %s\n", branches[0], arch)
			inFirst := comparison.GetDiffPackages(packageListsBranch1[branches[0]+"_"+arch], packageListsBranch2[branches[1]+"_"+arch])
			uniqueInFirst = append(uniqueInFirst, inFirst)
		}(arch, branches)

		wg.Add(1)
		go func(arch string, branches [2]string) {
			defer wg.Done()
			fmt.Printf("Getting packages only present in %s on %s\n", branches[1], arch)
			inSecond := comparison.GetDiffPackages(packageListsBranch2[branches[1]+"_"+arch], packageListsBranch1[branches[0]+"_"+arch])
			uniqueInSecond = append(uniqueInSecond, inSecond)
		}(arch, branches)

		wg.Add(1)
		go func(arch string, branches [2]string) {
			defer wg.Done()
			fmt.Printf("Getting packages with greater versions in %s on %s\n", branches[0], arch)
			greaterPkgVersions := comparison.GetGreaterPackagesVersions(packageListsBranch1[branches[0]+"_"+arch], packageListsBranch2[branches[1]+"_"+arch])
			greaterPackagesVersions = append(greaterPackagesVersions, greaterPkgVersions)
		}(arch, branches)

	}

	wg.Wait()
	fmt.Printf("Packages processed in %v\n", time.Now().Sub(t2).Seconds())
	fmt.Println()
	fmt.Println("Encoding json")
	jsonEncoded, err := json.MarshalIndent(models.TotalBranchPackagesResult{
		UniqueInFirst:         uniqueInFirst,
		UniqueInSecond:        uniqueInSecond,
		GreaterPackageVersion: greaterPackagesVersions,
	}, "", " ")
	if err != nil {
		log.Fatalf("Error encoding json: %s", err.Error())
	}
	fmt.Println("Json Encoded")

	filename := "out.txt"

	CreateFile(filename, string(jsonEncoded))

	fmt.Println("Json written in ", filename)

}
