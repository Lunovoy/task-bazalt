package cli

import (
	"fmt"
	"os"
)

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
