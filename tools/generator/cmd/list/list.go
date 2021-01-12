package list

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	listCmd := &cobra.Command{
		Use: "list <package searching root directory>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := commandContext{
				root: args[0],
			}
			return ctx.execute()
		},
	}

	return listCmd
}

type commandContext struct {
	root string
}

func (c *commandContext) execute() error {

	return nil
}
