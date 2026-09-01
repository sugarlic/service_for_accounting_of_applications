package config

type Config struct {
	// Server
	HTTPAddr string `env:"HTTP_ADDR" envDefault:":8080"`
	AppEnv   string `env:"APP_ENV"   envDefault:"production"`

	// Database
	DatabaseURL string `env:"DATABASE_URL"`
	CRMAPIKey   string `env:"CRM_API_KEY"`
}
