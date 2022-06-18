/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/phsiao/idxdo/pkg/ceramic"
	"github.com/phsiao/idxdo/pkg/idx"
	"github.com/spf13/cobra"
)

// idxCmd represents the idx command
var idxCmd = &cobra.Command{
	Use:   "idx",
	Short: "Inspect and validate your IDX document",
	Long: `
IDX is the directory that identity records such as GitCoin Passport are stored,
and you need it to find your Passport and other identity documents, so you
should be able to inspect and validate your IDX document.
	
Subcommands in this category help you inspect and validate your IDX document.`,
}

// idxIdCmd represents the 'idx id' command
var idxIdCmd = &cobra.Command{
	Use:   "id",
	Short: "Get your IDX document StreamID",
	Long: `
Each Ceramic Stream has a StreamID, and the StreamID is computed from your
DID such as pkh.

The output is your StreamID.`,
}

var (
	pkhChainId uint
	pkhAccount string
)

// idxIdPkhCmd represents the 'idx id pkh' command
var idxIdPkhCmd = &cobra.Command{
	Use:   "pkh",
	Short: "Get your IDX document StreamID using PKH",
	Long: `
Each Ceramic Stream has a StreamID, and the StreamID is computed from your
DID such as pkh.

The output is your StreamID. You can use the StreamID with 'idx state' or
'idx record' to inspect your records.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if pkhAccount == "" {
			return fmt.Errorf("please provides an account using --account")
		}
		streamid := idx.StreamIDFromPKH(pkhChainId, pkhAccount)
		fmt.Println(streamid)
		return nil
	},
}

// idxStateCmd represents the 'idx state' command
var idxStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Get your IDX document state",
	Long: `
Each Ceramic Stream has a state.  The output is the state for your IDX document.
`,
	ArgAliases: []string{"streamid"},
	Args:       cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		streamid := args[0]
		api := ceramic.NewAPI()
		response, err := api.GetStream(streamid)
		if err != nil {
			return err
		}

		out, err := colorPrettyJson(response.State)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

const (
	GITCOIN_PASSPORT_DEFINITION = "kjzl6cwe1jw148h1e14jb5fkf55xmqhmyorp29r9cq356c7ou74ulowf8czjlzs"
)

// idxRecordCmd represents the 'idx record' command
var idxRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Get your IDX idenity records",
	Long: `
Your GitCoin passport and other identity documents are stored in the IDX index
as record.  You should be able to see what records you have in your IDX.
`,
	ArgAliases: []string{"streamid"},
	Args:       cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		streamid := args[0]
		api := ceramic.NewAPI()
		response, err := api.GetStream(streamid)
		if err != nil {
			return err
		}

		content := map[string]string{}
		err = json.Unmarshal(response.State.Content, &content)
		if err != nil {
			return err
		}

		fmt.Printf("Showing %d available record(s)\n", len(content))
		for definition, record := range content {
			fmt.Printf("=> Inspecting %s -> %s\n", definition, record)
			switch definition {
			case GITCOIN_PASSPORT_DEFINITION:
				fmt.Println("   Found GitCoin Passport record")
				u, err := url.Parse(record)
				if err != nil {
					return err
				}
				response, err := api.GetStream(u.Host)
				if err != nil {
					return err
				}
				out, err := colorPrettyJson(response.State)
				if err != nil {
					return err
				}
				fmt.Println(string(out))
			default:
				fmt.Printf("   Unknown record definition: %s\n", definition)
				response, err := api.GetStream(definition)
				if err != nil {
					return err
				}
				out, err := colorPrettyJson(response.State)
				if err != nil {
					return err
				}
				fmt.Println(string(out))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(idxCmd)
	idxCmd.AddCommand(idxIdCmd)
	idxIdCmd.AddCommand(idxIdPkhCmd)
	idxCmd.AddCommand(idxStateCmd)
	idxCmd.AddCommand(idxRecordCmd)

	idxIdPkhCmd.Flags().UintVar(&pkhChainId, "chainid", 1, "EIP-155 Chain ID to use for your identity")
	idxIdPkhCmd.Flags().StringVar(&pkhAccount, "account", "", "Account to use for your identity")
}
