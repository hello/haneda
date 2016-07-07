package core

type ServerConfig struct {
	ExternalHost string `toml:"external_host"`
	InternalHost string `toml:"internal_host"`
}

type ProxyConfig struct {
	Endpoint string
}

type RedisConfig struct {
	Host   string
	PubSub string
}

type AwsConfig struct {
	Region            string
	Key               string
	Secret            string
	KeyStoreTableName string `toml:"keystore_table"`
}

type GraphiteConfig struct {
	Host   string
	Prefix string
}

type HelloConfig struct {
	Server   *ServerConfig
	Proxy    *ProxyConfig
	Redis    *RedisConfig
	Aws      *AwsConfig
	Graphite *GraphiteConfig
}
