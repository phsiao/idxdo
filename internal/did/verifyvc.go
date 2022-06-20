package did

import (
	"fmt"
	"strings"

	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	ldstore "github.com/hyperledger/aries-framework-go/pkg/store/ld"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
	"github.com/phsiao/idxdo/pkg/ceramic"
)

// VerifyVC verifies a limited set of Verifiable Credential.
// It returns error when the input can't be verified.
//
// Supported VC must be signed using did:key: method, and
// with Ed25519 signature scheme
func VerifyVC(input []byte) (*verifiable.Credential, error) {
	// build loader
	loader, err := getJSONLDLoader()
	if err != nil {
		return nil, err
	}

	// parse the document without proof check first, so we can
	// extract the issuer and verification method
	nvc, err := verifiable.ParseCredential(input,
		verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDDocumentLoader(loader))
	if err != nil {
		return nil, err
	}
	issuer := nvc.Issuer.ID
	pkDid, err := did.Parse(issuer)
	if err != nil {
		return nil, err
	}
	if pkDid.Method != "key" {
		return nil, fmt.Errorf("unsupport DID method for issuer")
	}

	// extract a public for verification
	var pubkey []byte
	for _, proof := range nvc.Proofs {
		if v, found := proof["verificationMethod"]; found {
			switch val := v.(type) {
			case string:
				// the verificationMethod must match the issuer
				if strings.Index(string(val), issuer) == 0 {
					_, output, err := multibase.Decode(pkDid.MethodSpecificID)
					if err != nil {
						return nil, err
					}
					mc, pk, _, err := ceramic.GetUVarInt(output)
					if err != nil {
						return nil, err
					}
					if multicodec.Code(mc) != multicodec.Ed25519Pub {
						return nil, fmt.Errorf("unsupported signing key type %s", multicodec.Code(mc).String())
					}
					pubkey = pk
				}
			default:
				continue
			}
		}
	}
	if len(pubkey) == 0 {
		return nil, fmt.Errorf("cannot find matching public key for verification")
	}

	// verify the proof of the document
	vc, err := verifiable.ParseCredential(input,
		verifiable.WithPublicKeyFetcher(verifiable.SingleKey(pubkey, kms.ED25519)),
		verifiable.WithJSONLDDocumentLoader(loader))
	if err != nil {
		return nil, err
	}

	return vc, nil
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
