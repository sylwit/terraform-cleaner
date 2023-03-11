package cmd

import (
	"fmt"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var rootCmd = &cobra.Command{
	Use:   "terraform-cleaner <path>",
	Short: "Remove unused variables",
	RunE:  rootCmdExec,
}

func rootCmdExec(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	cmd.SilenceUsage = true

	fUnusedOnly, _ := cmd.Flags().GetBool("unused-only")

	modules, err := findTfModules(dir)
	if err != nil {
		return err
	}

	for path := range modules {
		stats, err := findVariablesUsage(path)
		if err != nil {
			return err
		}

		if fUnusedOnly {
			for name, count := range stats.variables {
				if count > 0 {
					delete(stats.variables, name)
				}
			}
			if len(stats.variables) == 0 {
				continue
			}
		}

		err = displayModule(path, &stats)
		if err != nil {
			return err
		}

	}

	fmt.Printf("%d modules processed", len(modules))

	return nil
}

func displayModule(path string, stats *VariablesUsage) error {
	fmt.Printf("Module: %s (%d variables found)\n", path, len(stats.variables))

	for name, count := range stats.variables {
		fmt.Printf("%s : %d\n", name, count)
	}
	fmt.Println("")

	return nil
}

type VariablesUsage struct {
	variables map[string]int
}

func findVariablesUsage(path string) (VariablesUsage, error) {
	out := VariablesUsage{variables: map[string]int{}}

	module, diagnostics := tfconfig.LoadModule(path)
	if diagnostics.HasErrors() {
		return out, diagnostics.Err()
	}

	result, err := countVariables(path, module)
	if err != nil {
		return out, err
	}

	out.variables = result

	return out, nil
}

func countVariables(path string, tfconfig *tfconfig.Module) (map[string]int, error) {
	out := map[string]int{}

	files, err := os.ReadDir(path)
	if err != nil {
		return out, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tf" {
			data, err := os.ReadFile(filepath.Join(path, file.Name()))
			if err != nil {
				return out, err
			}

			content := string(data)

			for variable := range tfconfig.Variables {
				regex := regexp.MustCompile(fmt.Sprintf(`var\.%s\W`, variable))
				matches := regex.FindAllStringIndex(content, -1)

				out[variable] += len(matches)
			}

		}
	}

	return out, err
}

func findTfModules(path string) (map[string]bool, error) {
	var directories = make(map[string]bool)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".tf" {
			module := filepath.Dir(path)
			log.Debugf("Visited: %s\n", module)
			if _, ok := directories[module]; !ok {
				directories[module] = true
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return directories, nil
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().Bool("unused-only", false, "Display only unused variables")
}
