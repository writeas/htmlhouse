package htmlhouse

import (
	"github.com/danryan/env"
)

type config struct {
	StaticDir string `env:"key=STATIC_DIR default=static"`

	//	PrivateKey string `env:"key=PRIVATE_KEY require=true"`
	//	PublicKey  string `env:"key=PUBLIC_KEY require=true"`

	ServerPort int `env:"key=PORT default=8080"`
}

func newConfig() (*config, error) {
	cfg := &config{}
	if err := env.Process(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
