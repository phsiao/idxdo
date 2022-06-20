package cidinspect

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
)

// ToHumanReadable() parses a CID into its components in a human readable string
// https://cid.ipfs.io/ is a better tool but this is good enough for debugging
func ToHumanReadable(input string) (string, error) {
	parsed, err := cid.Parse(input)
	if err != nil {
		return "", err
	}

	decoded, err := multihash.Decode(parsed.Hash())
	if err != nil {
		return "", err
	}

	prefix := parsed.Prefix()
	return fmt.Sprintf("cid([v%x %s(0x%x) %s(0x%x) len=%d] %x) = %s",
		prefix.Version, multicodec.Code(prefix.Codec).String(), prefix.Codec,
		multicodec.Code(prefix.MhType).String(), prefix.MhType, prefix.MhLength,
		decoded.Digest, input), nil
}
