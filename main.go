package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func promptUser(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func createDirectoryStructure(programType, programName, domain string) {
	directoryPath := fmt.Sprintf("%s/%s/%s", programType, programName, domain)
	os.MkdirAll(directoryPath, 0755)

	command := exec.Command("sh", "-c", fmt.Sprintf("subfinder -d %s --all | tee domains.txt", domain))
	command.Dir = directoryPath
	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("%s\n", out.String())
}

func addProgram(programType, programName string) {
	domainsInput := promptUser("Enter domain names (space-separated): ")
	domains := strings.Fields(domainsInput)

	for _, domain := range domains {
		createDirectoryStructure(programType, programName, domain)
	}
}

func discoveredDomains() error {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "domains.txt" {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Println(string(fileContent))
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	var rootCmd = &cobra.Command{Use: "bbounty"}
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a program with domain names",
		Run: func(cmd *cobra.Command, args []string) {
			programType, _ := cmd.Flags().GetString("program")
			programName, _ := cmd.Flags().GetString("name")

			if programType == "" {
				programType = promptUser("Enter program type (vdp or bbp): ")
			}
			if programName == "" {
				programName = promptUser("Enter program name: ")
			}
			addProgram(programType, programName)
		},
	}
	addCmd.Flags().StringP("program", "p", "", "Specify the program type (vdp or bbp)")
	addCmd.Flags().StringP("name", "n", "", "Specify the program name")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List discovered domains",
		Run: func(cmd *cobra.Command, args []string) {
			discoveredDomains()
		},
	}

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
