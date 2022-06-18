package ceramic

const (
	GITCOIN_PASSPORT_CERAMIC_ENDPOINT  = "ceramic.passport-iam.gitcoin.co"
	MAINNET_COMMUNITY_CERAMIC_ENDPOINT = "gateway.ceramic.network"
	CLAY_COMMUNITY_CERAMIC_ENDPOINT    = "gateway-clay.ceramic.network"
)

// API interacts with Ceramic API gateway
type API struct {
	host   string
	scheme string
	root   string
}

func NewAPI() API {
	return API{
		host:   MAINNET_COMMUNITY_CERAMIC_ENDPOINT,
		scheme: "https",
		root:   "",
	}
}
