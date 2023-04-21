package models

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
