package cmd

import (
	"fmt"

	"github.com/phsiao/idxdo/pkg/cidinspect"
	"github.com/spf13/cobra"
)

// cidCmd represents the cid command
var cidCmd = &cobra.Command{
	Use:   "cid",
	Short: "CID utility command category",
	Long: `
CID is very powerful and this command helps you understanding it.
Commands in this category help interacting with cid.
`,
}

// cidInspectCmd represents the 'cid inspect' command
var cidInspectCmd = &cobra.Command{
	Use:   "inspect [flags] <cid-to-inspect>",
	Short: "Decode CID into a more user-friendly form",
	Long: `
Decode CID into a more user-friendly form.
`,
	ArgAliases: []string{"cid"},
	Args:       cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cidinspect.IsCID(args[0]) {
			return fmt.Errorf("argument %s is not a valid CID", args[0])
		}
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
