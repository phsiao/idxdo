/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

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

The output is your StreamID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if pkhAccount == "" {
			return fmt.Errorf("please provides an account using --account")
		}
		streamid := idx.StreamIDFromPKH(pkhChainId, pkhAccount)
		fmt.Println(streamid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(idxCmd)
	idxCmd.AddCommand(idxIdCmd)
	idxIdCmd.AddCommand(idxIdPkhCmd)

	idxIdPkhCmd.Flags().UintVar(&pkhChainId, "chainid", 1, "EIP-155 Chain ID to use for your identity")
	idxIdPkhCmd.Flags().StringVar(&pkhAccount, "account", "", "Account to use for your identity")
}
