/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/phsiao/idxdo/pkg/ceramic"
	"github.com/phsiao/idxdo/pkg/idx"
	"github.com/spf13/cobra"
)

// gpCmd represents the gp command
var gpCmd = &cobra.Command{
	Use:   "gp",
	Short: "Gitcoin passport utility command category",
	Long: `
Gitcoin Passport utility command category.
`,
}

type Stamp struct {
	Provider   string `json:"provider"`
	Credential json.RawMessage
}

type GitcoinPassport struct {
	ExpiryDate   string  `json:"expiryDate"`
	IssuanceDate string  `json:"issuanceDate"`
	Stamps       []Stamp `json:"stamps"`
}

// gpDumpCmd represents the 'gp backup' command
var gpDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump the content of the account's Gitcoin Passport",
	Long: `
Dump the content of the account's Gitcoin Passport.  The output of
this command is a JSON document contains stamps in the format of
Verifiable Credential that are signed by Gitcoin Passport.

Gitcoin stamps contains a hash of your verified identity with a
secret key, so anyone can see if two stamps are issued to the same
verified identity but no one can see what is the verified identity.
If you know the Gitcoin Passport private key and the identity you
can check if it matches the hash.

Gitcoin stamps are also signed using Gitcoin Passport's key
(did:key:z6MkghvGHLobLEdj1bgRLhS4LPGJAvbMA1tn2zcRyqmYU5LC) and can
be validated by anyone.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if pkhAccount == "" {
			return fmt.Errorf("please provides an account using --account")
		}
		if !IsEthereumAccount(pkhAccount) {
			return fmt.Errorf("argument %s is not a valid account", pkhAccount)
		}
		idxStreamid := idx.StreamIDFromPKH(pkhChainId, pkhAccount)
		api := ceramic.NewAPI(ceramic.WithHost(ceramic.GITCOIN_PASSPORT_CERAMIC_ENDPOINT))
		response, err := api.GetStream(idxStreamid)
		if err != nil {
			panic(err)
		}

		content := map[string]string{}
		err = json.Unmarshal(response.State.Content, &content)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("Showing %d available record(s)\n", len(content))
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
				// fmt.Printf("=> Found Gitcoin Passport record at %s\n", record)
				content := response.State.Content
				if response.State.Next != nil {
					content = *response.State.Next.Content
				}
				passport := GitcoinPassport{}
				err = json.Unmarshal(content, &passport)
				if err != nil {
					panic(err)
				}
				for i, stamp := range passport.Stamps {
					tokens := strings.Split(string(stamp.Credential), "://")
					cred := tokens[1][:len(tokens[1])-1] // get rid of trailing "
					response, err := api.GetStream(cred)
					if err != nil {
						panic(err)
					}
					content := response.State.Content

					passport.Stamps[i].Credential = content
				}
				out, err := colorPrettyJson(passport)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%v\n", string(out))
			default:
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gpCmd)
	gpCmd.AddCommand(gpDumpCmd)

	gpDumpCmd.Flags().UintVar(&pkhChainId, "chainid", 1, "EIP-155 Chain ID to use for your identity")
	gpDumpCmd.Flags().StringVar(&pkhAccount, "account", "", "Account to use for your identity")
}
