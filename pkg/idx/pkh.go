package idx

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
	"github.com/phsiao/idxdo/pkg/ceramic"
)

type genesisHeader struct {
	controllers []string
	family      *string
}

type genesisCommit struct {
	header genesisHeader
}

// StreamIDFromPKH calculate StreamIDFromPKH from chainid and account
// The process is described here:
// https://github.com/ceramicnetwork/ceramic/blob/main/SPECIFICATION.md#genesis-commit
//
// The IDX genesis commit has no data and usually is done using dag-cbor
// https://github.com/ceramicnetwork/CIP/blob/main/CIPs/CIP-8/CIP-8.md
//  "The genesis commit can also be stored without a signature using dag-cbor.
//   However in this case the data property MUST be set to null."
//
func StreamIDFromPKH(chainid uint, account string) string {
	did := fmt.Sprintf("did:pkh:eip155:%d:%s", chainid, strings.ToLower(account))
	family := "IDX"

	// these are the fields used by IDX
	gc, err := ipld.LoadSchemaBytes([]byte(`
	  type GenesisHeader struct {
		controllers [String]
		family optional String
	  }
	  type GenesisCommit struct {
		header GenesisHeader
	  }
	`))
	if err != nil {
		panic(err)
	}
	gct := gc.TypeByName("GenesisCommit")

	ghv := genesisHeader{
		family:      &family,
		controllers: []string{did},
	}
	gcv := &genesisCommit{
		header: ghv,
	}

	node := bindnode.Wrap(gcv, gct)
	nodeRepr := node.Representation()
	cidBuffer := bytes.Buffer{}
	dagcbor.Encode(nodeRepr, &cidBuffer)

	prefix := cid.Prefix{
		Version:  1,
		Codec:    uint64(multicodec.DagCbor),
		MhType:   multihash.SHA2_256,
		MhLength: -1,
	}
	c, err := prefix.Sum(cidBuffer.Bytes())
	if err != nil {
		panic(err)
	}

	typeBuffer := bytes.Buffer{}
	typeBuffer.Write(ceramic.PutUVarInt(0))
	typeBuffer.Write(c.Bytes())

	codecBuffer := bytes.Buffer{}
	codecBuffer.Write(ceramic.PutUVarInt(uint64(multicodec.Streamid)))
	codecBuffer.Write(typeBuffer.Bytes())

	result, err := multibase.Encode(multibase.Base36, codecBuffer.Bytes())
	if err != nil {
		panic(err)
	}
	return result
}
