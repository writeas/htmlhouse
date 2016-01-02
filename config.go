package htmlhouse

import (
	"github.com/danryan/env"
)

type config struct {
	StaticDir string `env:"key=STATIC_DIR default=static"`

	DBName     string `env:"key=DB_DB required=true"`
	DBUser     string `env:"key=DB_USER required=true"`
	DBPassword string `env:"key=DB_PASSWORD required=true"`
	DBHost     string `env:"key=DB_HOST default=localhost"`

	PrivateKey string `env:"key=PRIVATE_KEY require=true"`
	PublicKey  string `env:"key=PUBLIC_KEY require=true"`

	ServerPort int `env:"key=PORT default=8080"`

	AutoApprove bool   `env:"key=AUTO_APPROVE default=false"`
	AdminPass   string `env:"key=ADMIN_PASS default=uhoh"`
}

func newConfig() (*config, error) {
	cfg := &config{}
	if err := env.Process(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
