package ceramic

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
)

// StreamID supports common operations with StreamID
// https://github.com/ceramicnetwork/CIP/blob/main/CIPs/CIP-59/CIP-59.md
type Common struct {
	StreamID string
	Encoding multibase.Encoding
	Type     uint64
}

type RecordEntry struct {
	CIDVersion  uint64
	ContentType multicodec.Code
	ContentID   string
}

type StreamID struct {
	Entry RecordEntry
	Common
}

type CommitID struct {
	Entries []RecordEntry
	Common
}

// Decode can return either StreamID or CommitID
func Decode(id string) (interface{}, error) {
	sid := Common{
		StreamID: id,
	}

	encoding, streamCodecRemainder, err := multibase.Decode(id)
	if err != nil {
		return nil, err
	}
	sid.Encoding = encoding

	// check <multicodec-streamCodec>
	streamCodec, streamCodecRemainder, _, err := GetUVarInt(streamCodecRemainder)
	if err != nil {
		return nil, err
	}
	if multicodec.Code(streamCodec) != multicodec.Streamid {
		return nil, fmt.Errorf("unexpected multicodec %x != 0xce", streamCodecRemainder[0])
	}

	// check <stream-type>
	streamtype, streamTypeRemainder, _, err := GetUVarInt(streamCodecRemainder)
	if err != nil {
		return nil, err
	}
	sid.Type = streamtype

	entries := []RecordEntry{}

	for len(streamTypeRemainder) > 1 {
		entry := RecordEntry{}

		// check cid version
		cidversion, cidVersionRemainder, _, err := GetUVarInt(streamTypeRemainder)
		if err != nil {
			return nil, err
		}
		entry.CIDVersion = cidversion

		// check codec
		codec, codecRemainder, _, err := GetUVarInt(cidVersionRemainder)
		if err != nil {
			return nil, err
		}
		entry.ContentType = multicodec.Code(codec)

		_, mhCodecRemainder, mhCodecLength, err := GetUVarInt(codecRemainder)
		if err != nil {
			return nil, err
		}

		mhLength, _, mhLengthLength, err := GetUVarInt(mhCodecRemainder)
		if err != nil {
			return nil, err
		}

		idx := mhCodecLength + mhLengthLength + int(mhLength)
		multihashBytes := codecRemainder[:idx]
		multihashBytesRemainder := codecRemainder[idx:]
		_, err = multihash.Decode(multihashBytes)
		if err != nil {
			return nil, err
		}

		var c cid.Cid
		if entry.CIDVersion == 0 {
			c = cid.NewCidV0(multihashBytes)
		} else {
			c = cid.NewCidV1(codec, multihashBytes)
		}
		entry.ContentID = c.String()
		entries = append(entries, entry)

		streamTypeRemainder = multihashBytesRemainder
	}

	if len(entries) == 1 {
		return StreamID{
			Entry:  entries[0],
			Common: sid,
		}, nil
	} else {
		return CommitID{
			Entries: entries,
			Common:  sid,
		}, nil
	}
}

// IsStreamID returns whether the given input is a valid StreamID
func IsStreamID(input string) bool {
	_, err := Decode(input)
	return err == nil
}
