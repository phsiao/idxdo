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

type Opt func(api *API) *API

func NewAPI(opts ...Opt) API {
	api := API{
		host:   MAINNET_COMMUNITY_CERAMIC_ENDPOINT,
		scheme: "https",
		root:   "",
	}
	for _, opt := range opts {
		opt(&api)
	}

	return api
}

func WithHost(host string) Opt {
	return func(api *API) *API {
		api.host = host
		return api
	}
}
