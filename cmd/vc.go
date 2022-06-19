package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	ldstore "github.com/hyperledger/aries-framework-go/pkg/store/ld"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
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

type provider struct {
	ContextStore        ldstore.ContextStore
	RemoteProviderStore ldstore.RemoteProviderStore
}

func (p *provider) JSONLDContextStore() ldstore.ContextStore {
	return p.ContextStore
}

func (p *provider) JSONLDRemoteProviderStore() ldstore.RemoteProviderStore {
	return p.RemoteProviderStore
}

func getJSONLDLoader() (*ld.DocumentLoader, error) {
	contextStore, err := ldstore.NewContextStore(mem.NewProvider())
	if err != nil {
		return nil, fmt.Errorf("create JSON-LD context store: %v", err)
	}

	remoteProviderStore, err := ldstore.NewRemoteProviderStore(mem.NewProvider())
	if err != nil {
		return nil, fmt.Errorf("create remote JSON-LD context provider store: %v", err)
	}

	p := &provider{
		ContextStore:        contextStore,
		RemoteProviderStore: remoteProviderStore,
	}
	loader, err := ld.NewDocumentLoader(p)
	if err != nil {
		return nil, err
	}
	return loader, nil
}

// vcVerifyCmd represents the 'vc verify' command
var vcVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify the Verifiable Credential in stdin or the specified StreamID",
	Long: `
A StreamID may contain a Verifiable Credential and subcommands in this
category can interact with them.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte
		if len(args) > 0 {
			streamid := args[0]
			api := ceramic.NewAPI()

			// get document by StreamID
			response, err := api.GetStream(streamid)
			if err != nil {
				return err
			}
			input = response.State.Content
		} else {
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			input = bytes
		}

		// build loader
		loader, err := getJSONLDLoader()
		if err != nil {
			return err
		}

		// parse the document without proof check first, so we can
		// extract the issuer and verification method
		nvc, err := verifiable.ParseCredential(input,
			verifiable.WithDisabledProofCheck(),
			verifiable.WithJSONLDDocumentLoader(loader))
		if err != nil {
			return err
		}
		issuer := nvc.Issuer.ID
		didTokens := strings.Split(issuer, ":")
		if len(didTokens) < 3 || didTokens[0] != "did" || didTokens[1] != "key" {
			return fmt.Errorf("unsupport DID method for issuer")
		}

		// extract a public for verification
		var pubkey []byte
		for _, proof := range nvc.Proofs {
			if v, found := proof["verificationMethod"]; found {
				switch val := v.(type) {
				case string:
					if strings.Index(string(val), issuer) == 0 {
						_, output, err := multibase.Decode(didTokens[2])
						if err != nil {
							return err
						}
						mc, pk, _, err := ceramic.GetUVarInt(output)
						if err != nil {
							return err
						}
						if multicodec.Code(mc) != multicodec.Ed25519Pub {
							return fmt.Errorf("unsupported signing key type %s", multicodec.Code(mc).String())
						}
						pubkey = pk
					}
				default:
					continue
				}
			}
		}

		// verify the proof of the document
		vc, err := verifiable.ParseCredential(input,
			verifiable.WithPublicKeyFetcher(verifiable.SingleKey(pubkey, kms.ED25519)),
			verifiable.WithJSONLDDocumentLoader(loader))
		if err != nil {
			return err
		}

		out, err := colorPrettyJson(vc)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(vcCmd)
	vcCmd.AddCommand(vcVerifyCmd)
}
