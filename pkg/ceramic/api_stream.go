package ceramic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	COMMIT_GENESIS uint = iota
	COMMIT_SIGNED
	COMMIT_ANCHOR
)

const (
	SIGNATURE_STATUS_GENESIS uint = iota
	SIGNATURE_STATUS_PARTIAL
	SIGNATURE_STATUS_SIGNED
)

type CeramicAnchorProof struct {
	BlockNumber    uint   `json:"blockNumber"`
	BlockTimestamp uint   `json:"blockTimestamp"`
	ChainId        string `json:"chainId"`
	Root           string `json:"root"`
	TxHash         string `json:"txHash"`
}

type CeramicLogEntry struct {
	CID       string `json:"cid"`
	Timestamp *uint  `json:"timestamp,omitempty"`
	Type      uint   `json:"type"`
}

type CeramicStreamMetadata struct {
	Controllers            []string `json:"controllers"`
	Family                 *string  `json:"family,omitempty"`
	ForbidControllerChange *bool    `json:"forbidControllerChange,omitempty"`
	Schema                 *string  `json:"schema,omitempty"`
	Tags                   []string `json:"tags,omitempty"`
}

type CeramicStreamNext struct {
	Content     *json.RawMessage       `json:"content,omitempty"`
	Controllers []string               `json:"controllers,omitempty"`
	Metadata    *CeramicStreamMetadata `json:"metadata,omitempty"`
}

// https://developers.ceramic.network/reference/typescript/interfaces/_ceramicnetwork_common.streamstate-1.html
type CeramicStreamState struct {
	AnchorProof *CeramicAnchorProof `json:"anchorProof,omitempty"`

	// AnchorStatus not supoorted yet

	Content         json.RawMessage       `json:"content"`
	Log             []CeramicLogEntry     `json:"log"`
	Metadata        CeramicStreamMetadata `json:"metadata"`
	Next            *CeramicStreamNext    `json:"next,omitempty"`
	SignatureStatus uint                  `json:"signature"`
	Type            int                   `json:"type"`
}

type CeramicGetResponse struct {
	StreamID string             `json:"streamId"`
	State    CeramicStreamState `json:"state"`
}

func (a *API) GetStream(streamid string) (*CeramicGetResponse, error) {
	handle := url.URL{
		Host:   a.host,
		Scheme: a.scheme,
		Path:   fmt.Sprintf("%s/api/v0/streams/%s", a.root, streamid),
	}
	resp, err := http.Get(handle.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := CeramicGetResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
