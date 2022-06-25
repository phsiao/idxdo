package cmd

import (
	"fmt"

	"github.com/phsiao/idxdo/pkg/ceramic"
	"github.com/phsiao/idxdo/pkg/cidinspect"
	"github.com/spf13/cobra"
)

// streamidCmd represents the streamid command
var streamidCmd = &cobra.Command{
	Use:   "streamid",
	Short: "Help understand StreamID",
	Long: `
StreamID is very opaque and this helps with interacting with it
`,
}

// streamidInspectCmd represents the streamid inspect command
var streamidInspectCmd = &cobra.Command{
	Use:   "inspect [flags] <streamid>",
	Short: "Decode StreamID into a more user-friendly form",
	Long: `
StreamID is very opaque and this helps with decoding it
`,
	ArgAliases: []string{"streamid"},
	Args:       cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !ceramic.IsStreamID(args[0]) {
			return fmt.Errorf("argument %s is not a valid StreamID", args[0])
		}
		streamid := args[0]

		id, err := ceramic.Decode(streamid)
		if err != nil {
			return err
		}
		switch obj := id.(type) {
		case ceramic.StreamID:
			fmt.Println("Type: StreamID")
			out, err := cidinspect.ToHumanReadable(obj.Entry.ContentID)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", out)
		case ceramic.CommitID:
			fmt.Println("Type: CommitID")
			for idx, entry := range obj.Entries {
				out, err := cidinspect.ToHumanReadable(entry.ContentID)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%d: %s\n", idx, out)
			}
		}

		return nil
	},
}

// streamidStateCmd represents the streamid state command
var streamidStateCmd = &cobra.Command{
	Use:   "state [flags] <streamid>",
	Short: "Show StreamID's entire state",
	Long: `
Download and pretty print the state of a StreamID
`,
	ArgAliases: []string{"streamid"},
	Args:       cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !ceramic.IsStreamID(args[0]) {
			return fmt.Errorf("argument %s is not a valid StreamID", args[0])
		}
		streamid := args[0]
		api := ceramic.NewAPI()
		response, err := api.GetStream(streamid)
		if err != nil {
			panic(err)
		}

		out, err := colorPrettyJson(response.State)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
		return nil
	},
}

// streamidContentCmd represents the streamid content command
var streamidContentCmd = &cobra.Command{
	Use:   "content [flags] <streamid>",
	Short: "Show StreamID's content",
	Long: `
Download and pretty print the content of a StreamID
`,
	ArgAliases: []string{"streamid"},
	Args:       cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !ceramic.IsStreamID(args[0]) {
			return fmt.Errorf("argument %s is not a valid StreamID", args[0])
		}
		streamid := args[0]
		api := ceramic.NewAPI()
		response, err := api.GetStream(streamid)
		if err != nil {
			panic(err)
		}

		out, err := colorPrettyJson(response.State.Content)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(streamidCmd)
	streamidCmd.AddCommand(streamidInspectCmd)
	streamidCmd.AddCommand(streamidStateCmd)
	streamidCmd.AddCommand(streamidContentCmd)
}
