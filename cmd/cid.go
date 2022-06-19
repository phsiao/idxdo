package cmd

import (
	"fmt"

	"github.com/phsiao/idxdo/pkg/cidinspect"
	"github.com/spf13/cobra"
)

// cidCmd represents the cid command
var cidCmd = &cobra.Command{
	Use:   "cid",
	Short: "cid command category",
	Long: `
Subcommands in this category help interacting with cid.
`,
}

// cidInspectCmd represents the 'cid inspect' command
var cidInspectCmd = &cobra.Command{
	Use:   "inspect [flags] <cid-to-inspect>",
	Short: "Decode cid into a more user-friendly form",
	Long: `
CID is very powerful and this command helps you understanding it.
`,
	ArgAliases: []string{"cid"},
	Args:       cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := cidinspect.ToHumanReadable(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Println(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cidCmd)
	cidCmd.AddCommand(cidInspectCmd)
}
