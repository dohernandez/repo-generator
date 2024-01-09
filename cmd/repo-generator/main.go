// Package main cli.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	generator "github.com/dohernandez/repo-generator"
)

var (
	srcpkg string
	output string
	model  string
	create bool
	insert bool
	update bool
	del    bool
)

// ERR Unable to find 'foo' in any go files under this path dry-run=false version=v2.36.0

var rootCmd = &cobra.Command{
	Use:   "repo-generator",
	Short: "A CLI repo-generator",
	Long:  "Generate repo objects for your Golang model struct\n",
	//nolint:godox
	Version: "v0.2.1", // TODO: read from version.txt
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting %s version=%s\n", cmd.Use, cmd.Version)

		var (
			options  []generator.Option
			optFlags []string

			qualifiedName string
		)
		// Your logic for handling the selected functionalities
		if create {
			options = append(options, generator.WithCreateFunc())
			optFlags = append(optFlags, "--create")
		}

		if insert {
			options = append(options, generator.WithInsertFunc())
			optFlags = append(optFlags, "--insert")
		}

		if update {
			options = append(options, generator.WithUpdateFunc())
			optFlags = append(optFlags, "--update")
		}

		if del {
			options = append(options, generator.WithDeleteFunc())
			optFlags = append(optFlags, "--delete")
		}

		if model == "" {
			fmt.Println("\033[91mERR:\033[0m Use --model to specify the name of the model struct use to generate repo for")

			os.Exit(1)
		}

		var err error

		if srcpkg == "" {
			srcpkg, err = os.Getwd()
			if err != nil {
				panic(err)
			}
		}

		flags := strings.Join(optFlags, " ")

		qualifiedName, err = generator.FindModelSourcePath(model, srcpkg)
		if err != nil {
			fmt.Printf("\033[91mERR:\033[0m %s\n", err.Error())

			os.Exit(1)
		}

		// Additional logic for your application based on the provided arguments
		fmt.Printf("Generating mock model=%s qualified-name=%s %s version=%s\n",
			model,
			qualifiedName,
			flags,
			cmd.Version,
		)

		if output == "" {
			output = generateGenPath(qualifiedName)
		}

		if err := generator.Generate(
			qualifiedName,
			output,
			model,
			options...,
		); err != nil {
			fmt.Printf("\033[91mERR:\033[0m %s\n", err.Error())

			os.Exit(1)
		}

		fmt.Println("repo file generated!")
	},
}

//nolint:gochecknoinits
func init() {
	rootCmd.PersistentFlags().StringVar(&srcpkg, "srcpkg", "", "Source pkg to search for struct (required)")
	rootCmd.PersistentFlags().StringVar(&output, "output", "", "Path of the output file (required)")
	rootCmd.PersistentFlags().StringVar(&model, "model", "", "Name of the struct to generate repo for (required)")
	rootCmd.PersistentFlags().BoolVar(&create, "create", false, "Enable repo create functionality")
	rootCmd.PersistentFlags().BoolVar(&insert, "insert", false, "Enable repo insert functionality")
	rootCmd.PersistentFlags().BoolVar(&update, "update", false, "Enable repo update functionality")
	rootCmd.PersistentFlags().BoolVar(&del, "delete", false, "Enable repo delete functionality")

	rootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate the autocompletion script for the specified shell",
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation for the completion command
		fmt.Println("Generating autocompletion script...")
	},
}

func generateGenPath(originalPath string) string {
	// Extract the directory from the original path
	path := filepath.Dir(originalPath)

	// Extract the extension from the original path
	ext := filepath.Ext(originalPath)

	// Extract the filename without the extension
	filename := strings.TrimSuffix(filepath.Base(originalPath), ext)

	genFilename := fmt.Sprintf("%s_gen%s", filename, ext)

	// Append "_gen" to the filename and join with the directory
	return filepath.Join(path, genFilename)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
