package ceramic

const (
	GITCOIN_PASSPORT_CERAMIC_ENDPOINT = "ceramic.passport-iam.gitcoin.co"
)

// API interacts with Ceramic API gateway
type API struct {
	host   string
	scheme string
	root   string
}

func NewAPI() API {
	return API{
		host:   GITCOIN_PASSPORT_CERAMIC_ENDPOINT,
		scheme: "https",
		root:   "",
	}
}
