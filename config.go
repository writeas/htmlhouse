package htmlhouse

import (
	"github.com/danryan/env"
	"regexp"
	"strings"
)

type config struct {
	StaticDir string `env:"key=STATIC_DIR default=static"`

	DBName     string `env:"key=DB_DB required=true"`
	DBUser     string `env:"key=DB_USER required=true"`
	DBPassword string `env:"key=DB_PASSWORD required=true"`
	DBHost     string `env:"key=DB_HOST default=localhost"`

	PrivateKey string `env:"key=PRIVATE_KEY require=true"`
	PublicKey  string `env:"key=PUBLIC_KEY require=true"`

	HostName   string `env:"key=HOST default=https://html.house"`
	ServerPort int    `env:"key=PORT default=8080"`

	AutoApprove  bool   `env:"key=AUTO_APPROVE default=false"`
	PreviewsHost string `env:"key=PREVIEWS_HOST`
	AdminPass    string `env:"key=ADMIN_PASS default=uhoh"`
	BrowseItems  int    `env:"key=BROWSE_ITEMS default=10"`

	BlacklistTerms string `env:"key=BLACKLIST_TERMS"`
	BlacklistReg   *regexp.Regexp

	// Twitter configuration
	TwitterConsumerKey    string `env:"key=TWITTER_KEY default=notreal"`
	TwitterConsumerSecret string `env:"key=TWITTER_SECRET default=notreal"`
	TwitterToken          string `env:"key=TWITTER_TOKEN default=notreal"`
	TwitterTokenSecret    string `env:"key=TWITTER_TOKEN_SECRET default=notreal"`
}

func newConfig() (*config, error) {
	cfg := &config{}
	if err := env.Process(cfg); err != nil {
		return cfg, err
	}

	// Process anything
	termsReg := `(?i)\b` + cfg.BlacklistTerms + `\b`
	termsReg = strings.Replace(termsReg, ",", `\b|\b`, -1)
	cfg.BlacklistReg = regexp.MustCompile(termsReg)

	// Return result
	return cfg, nil
}
