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
	"gopkg.in/yaml.v3"
)

// execCommand executes a shell command in the specified directory with the given domain and command template.
func execCommand(directory string, domain string, commandTemplate string, verbose bool) {
	commandString := fmt.Sprintf(commandTemplate, domain)
	if verbose {
		fmt.Println("Executing: " + commandString)
	}
	command := exec.Command("sh", "-c", commandString)
	command.Dir = directory
	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("%s\n", out.String())
}

// promptUser prompts the user for input and returns the response.
func promptUser(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// createDirectoryStructure creates the directory structure for the program.
func createDirectoryStructure(programType, programName, domain string) string {
	directoryPath := fmt.Sprintf("%s/%s/%s", programType, programName, domain)
	os.MkdirAll(directoryPath, 0755)
	return directoryPath
}

// addProgram creates directories for each domain and returns a slice of created directories.
func addProgram(programType, programName string, domains []string) []string {
	directories := []string{}

	for _, domain := range domains {
		dir := createDirectoryStructure(programType, programName, domain)
		directories = append(directories, dir) // Only append the directory
	}

	return directories
}

// discoveredDomains lists discovered domains in the current directory.
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

	return err
}

// CommandConfig holds the command configuration.
type CommandConfig struct {
	Command string `yaml:"command"`
}

// loadCommandConfig loads the command configuration from a YAML file.
func loadCommandConfig(filename string) (string, error) {
	var config CommandConfig
	file, err := os.ReadFile(filename)

	if err != nil {
		return "", err // Return error if the file cannot be read
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return "", err // Return error if unmarshaling fails
	}
	return config.Command, nil // Return the command string
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	filePath := filepath.Join(homeDir, ".config/bbounty/config.yml")

	commandTemplate, err := loadCommandConfig(filePath)
	if err != nil {
		fmt.Println("Error loading command configuration:", err)
		return
	}

	// Default to using subfinder if no command is specified
	if commandTemplate == "" {
		commandTemplate = "subfinder -d %s --all | tee domains.txt"
	}

	var rootCmd = &cobra.Command{Use: "bbounty"}
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a program with domain names",
		Run: func(cmd *cobra.Command, args []string) {
			programType, _ := cmd.Flags().GetString("program")
			programName, _ := cmd.Flags().GetString("name")
			verbose, _ := cmd.Flags().GetBool("verbose")

			// Prompt for program type if not provided
			if programType == "" {
				programType = promptUser("Enter program type (vdp or bbp): ")
			}

			// Prompt for program name if not provided
			if programName == "" {
				programName = promptUser("Enter program name: ")
			}

			var domains []string
			scanner := bufio.NewScanner(os.Stdin)

			fmt.Println("Enter domain names (press Enter to finish):")

			// Read domains from standard input
			for scanner.Scan() {
				line := scanner.Text() // Get the current line
				if line == "" {
					break // Stop reading on empty line
				}
				domains = append(domains, strings.Fields(line)...)
			}

			// If no domains were read from stdin, prompt the user for them
			if len(domains) == 0 {
				domainsInput := promptUser("Enter domain names (space-separated): ")
				domains = strings.Fields(domainsInput)
			}

			// Create directories for all domains
			programs := addProgram(programType, programName, domains)

			// Execute the command for each domain
			for i, domain := range domains {
				execCommand(programs[i], domain, commandTemplate, verbose)
			}
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
