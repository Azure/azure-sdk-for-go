package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "Calls dep command to execute dep ensure -update",
	Long:  "This command will invoke the dep ensure -update command",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := theDepCommand()
		return err
	},
}

func init() {
	rootCmd.AddCommand(depCmd)
}

func theDepCommand() error {
	vprintln("Executing dep...")
	depArgs := "ensure -update"
	if verboseFlag {
		depArgs += "-v"
	}
	c := exec.Command("dep", strings.Split(depArgs, " ")...)
	err := startCmd(c)
	if err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}
	return c.Wait()
}
