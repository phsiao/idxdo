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
	Short: "GitCoin passport utilities",
	Long: `
GitCoin passport utility commands.
`,
}

type Stamp struct {
	Provider   string `json:"provider"`
	Credential json.RawMessage
}

type GitCoinPassport struct {
	ExpiryDate   string  `json:"expiryDate"`
	IssuanceDate string  `json:"issuanceDate"`
	Stamps       []Stamp `json:"stamps"`
}

// gpDumpCmd represents the 'gp backup' command
var gpDumpCmd = &cobra.Command{
	Use:   "dump [flags] account",
	Short: "GitCoin passport utilities",
	Long: `
GitCoin passport utility commands.
`,
	ArgAliases: []string{"account"},
	Args:       cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idxStreamid := idx.StreamIDFromPKH(1, args[0])
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
				// fmt.Printf("=> Found GitCoin Passport record at %s\n", record)
				content := response.State.Content
				if response.State.Next != nil {
					content = *response.State.Next.Content
				}
				passport := GitCoinPassport{}
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
}
