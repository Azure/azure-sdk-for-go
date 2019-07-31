package cmd

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/repo"
	"github.com/spf13/cobra"
)

var updateSpecsCmd = &cobra.Command{
	Use:   "updateSpec",
	Short: "Update the specs repo on master branch",
	Long: `This command will change the working directory to the specs folder,
	checkout to master branch and update it`,
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		specs := args[0]
		err := theUpdateSpecsCommand(specs)
		return err
	},
}

func init() {
	rootCmd.AddCommand(updateSpecsCmd)
}

func theUpdateSpecsCommand(spec string) error {
	vprintln("Updating specs repo...")
	err := changeDir(spec)
	if err != nil {
		return fmt.Errorf("failed to change directory to %s: %v", spec, err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get the current working directory: %v", err)
	}
	wt, err := repo.Get(cwd)
	if err != nil {
		return fmt.Errorf("failed to get the working tree: %v", err)
	}
	vprintf("Checking out to master branch in %s\n", cwd)
	err = wt.Checkout(master)
	if err != nil {
		return fmt.Errorf("checkout failed: %v", err)
	}
	err = wt.Pull(specUpstream, master)
	if err != nil {
		return fmt.Errorf("pull failed: %v", err)
	}
	return nil
}
