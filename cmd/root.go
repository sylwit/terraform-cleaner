package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sylwit/terraform-cleaner/terraform"
)

var rootCmd = &cobra.Command{
	Use:   "terraform-cleaner <path>",
	Short: "List variables and locals usage",
	RunE:  rootCmdExec,
}

func rootCmdExec(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	cmd.SilenceUsage = true

	fUnusedOnly, _ := cmd.Flags().GetBool("unused-only")
	fVariables, _ := cmd.Flags().GetBool("variables")
	fLocals, _ := cmd.Flags().GetBool("locals")

	dType := terraform.All
	if fVariables && !fLocals {
		dType = terraform.Variables
	} else if !fVariables && fLocals {
		dType = terraform.Locals
	}

	modules, err := terraform.ListTfModules(dir)
	if err != nil {
		return err
	}

	for path := range modules {
		moduleUsage, err := terraform.NewModuleUsage(path)
		if err != nil {
			return err
		}
		err = moduleUsage.Display(dType, fUnusedOnly)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\n%d modules processed", len(modules))

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().Bool("unused-only", false, "Display only unused values")
	rootCmd.Flags().Bool("variables", false, "Display only variables")
	rootCmd.Flags().Bool("locals", false, "Display only locals")
}
