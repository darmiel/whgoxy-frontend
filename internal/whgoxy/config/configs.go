package config

var ConfigOAuth OAuthConfig

type OAuthConfig struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
	Scopes       []string

	// Endpoint
	AuthURL  string
	TokenURL string

	// Cookie
	CookieHost   string
	CookieSecret string
	CookieName   string
}

var ConfigWeb WebConfig

type WebConfigTLS struct {
	AddrTLS  string
	CertFile string
	KeyFile  string
}

type WebConfig struct {
	WebDir string
	Addr   string
	TLS    WebConfigTLS
}
