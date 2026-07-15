package env

import "github.com/caarlos0/env/v11"

type proxyEnvConfig struct {
	ProxyUrl string `env:"PROXY_URL,required"`
}

type proxyConfig struct {
	raw proxyEnvConfig
}

func NewProxyEnvConfig() (*proxyConfig, error) {
	var raw proxyEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &proxyConfig{
		raw: raw,
	}, nil
}

func (conf *proxyConfig) URL() string {
	return conf.raw.ProxyUrl
}
