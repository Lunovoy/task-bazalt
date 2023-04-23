package comparison

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Lunovoy/task-bazalt/pkg/models"
)

// получает список пакетов для заданной ветки и архитектуры
func FetchPackages(branch, arch string) ([]models.Package, error) {
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

		branchBinaryPackages := models.BranchBinaryPackages{}
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

// выводит список пакетов, которые присутствуют в первом списке, но отсутствуют во втором списке
func GetDiffPackages(packages1 []models.Package, packages2 []models.Package) models.BranchPackages {
	t := time.Now()
	diffPackages := make([]models.Package, 0, len(packages1))
	diff := map[string]bool{}

	for _, pkg := range packages1 {
		diff[pkg.Name] = true
	}

	for _, pkg := range packages2 {
		delete(diff, pkg.Name)
	}

	for pkgName := range diff {
		for _, v := range packages1 {
			if v.Name == pkgName {

				diffPackages = append(diffPackages, v)
			}
		}
	}

	fmt.Println("Getting unique packages: ", len(diffPackages), time.Now().Sub(t).Seconds())
	return models.BranchPackages{
		Arch:     diffPackages[0].Arch,
		Packages: diffPackages,
	}

}

// выводит список пакетов, version-release которых больше в первом списке, чем во втором списке
func GetGreaterPackagesVersions(packages1 []models.Package, packages2 []models.Package) models.BranchPackages {
	packagesWithGreaterVersions := make([]models.Package, 0, len(packages1))
	t := time.Now()
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
	fmt.Println("Getting greater version: ", time.Now().Sub(t).Seconds())
	return models.BranchPackages{
		Arch:     packagesWithGreaterVersions[0].Arch,
		Packages: packagesWithGreaterVersions,
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
