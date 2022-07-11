package ntp

// Config NTP server config
type Config struct {
	Listener string
}

// ElasticSearchConfig elastic seach cluster config
type ElasticSearchConfig struct {
	Endpoints  []string
	IndexAlias string
	IndexSplit bool
}
