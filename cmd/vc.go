package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/phsiao/idxdo/internal/did"
	"github.com/phsiao/idxdo/pkg/ceramic"
	"github.com/spf13/cobra"
)

// vcCmd represents the vc command
var vcCmd = &cobra.Command{
	Use:   "vc",
	Short: "Perform Verifiable Credential operations on StreamID",
	Long: `
A StreamID may contain a Verifiable Credential and subcommands in this
category can interact with them.
`,
}

// vcVerifyCmd represents the 'vc verify' command
var vcVerifyCmd = &cobra.Command{
	Use:   "verify [flags] <streamid>",
	Short: "Verify the Verifiable Credential in stdin or the specified StreamID",
	Long: `
A StreamID may contain a Verifiable Credential and subcommands in this
category can interact with them.

This command reads the Verifiable Credential in StreamID if it is specified via
argument, or from stdin if no StreamID is given.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte
		if len(args) > 0 {
			streamid := args[0]
			api := ceramic.NewAPI()

			// get document by StreamID
			response, err := api.GetStream(streamid)
			if err != nil {
				panic(err)
			}
			input = response.State.Content
		} else {
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
			input = bytes
		}

		vc, err := did.VerifyVC(input)
		if err != nil {
			panic(err)
		}

		out, err := colorPrettyJson(vc)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(vcCmd)
	vcCmd.AddCommand(vcVerifyCmd)
}
