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
IDX is the directory that identity records such as Gitcoin Passport are stored,
and you need it to find your Passport and other identity documents, so you
should be able to inspect and validate your IDX document.
	
Subcommands in this category help you inspect and validate your IDX document.`,
}

// idxIdCmd represents the 'idx id' command
var idxIdCmd = &cobra.Command{
	Use:   "id",
	Short: "Get your IDX document StreamID from your DID",
	Long: `
Each Ceramic Stream has a StreamID. Your IDX index StreamID is computed
deterministically from your DID.  Different DID would result in different
IDX index StreamID.
`,
}

var (
	pkhChainId uint
	pkhAccount string
)

// idxIdPkhCmd represents the 'idx id pkh' command
var idxIdPkhCmd = &cobra.Command{
	Use:   "pkh",
	Short: "Get IDX document StreamID using pkh DID method",
	Long: `
Compute your IDX index StreamID from pkh DID method.

The output is your StreamID. You can use the StreamID with 'streamid state' or
'idx record' to inspect your records.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if pkhAccount == "" {
			return fmt.Errorf("please provides an account using --account")
		}
		if !IsEthereumAccount(pkhAccount) {
			return fmt.Errorf("argument %s is not a valid account", pkhAccount)
		}
		streamid := idx.StreamIDFromPKH(pkhChainId, pkhAccount)
		fmt.Println(streamid)
		return nil
	},
}

const (
	GITCOIN_PASSPORT_DEFINITION = "kjzl6cwe1jw148h1e14jb5fkf55xmqhmyorp29r9cq356c7ou74ulowf8czjlzs"
)

// idxRecordCmd represents the 'idx record' command
var idxRecordCmd = &cobra.Command{
	Use:   "record [flags] <streamid>",
	Short: "Get your IDX idenity records",
	Long: `
Your Gitcoin Passport and other identity documents are stored in the IDX index
as record.  You should be able to see what records you have in your IDX.
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

		content := map[string]string{}
		err = json.Unmarshal(response.State.Content, &content)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Showing %d available record(s)\n", len(content))
		for definition, record := range content {
			u, err := url.Parse(record)
			if err != nil {
				panic(err)
			}
			response, err := api.GetStream(u.Host)
			if err != nil {
				panic(err)
			}
			switch definition {
			case GITCOIN_PASSPORT_DEFINITION:
				fmt.Printf("=> Found Gitcoin Passport record at %s\n", record)
				content := response.State.Content
				if response.State.Next != nil {
					content = *response.State.Next.Content
				}
				out, err := colorPrettyJson(content)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(out))
			default:
				fmt.Printf("=> Found unknown record: %s\n", record)
				out, err := colorPrettyJson(response.State.Content)
				if err != nil {
					panic(err)
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
	idxCmd.AddCommand(idxRecordCmd)

	idxIdPkhCmd.Flags().UintVar(&pkhChainId, "chainid", 1, "EIP-155 Chain ID to use for your identity")
	idxIdPkhCmd.Flags().StringVar(&pkhAccount, "account", "", "Account to use for your identity")
}
