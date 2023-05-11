package config

type Config struct {
	GinMode       string `env:"GIN_MODE" envDefault:"debug"`
	AdminPassword string `env:"ADMIN_PASSWORD" envDefault:"xxx"`
	Gtag          string `env:"GTAG" envDefault:"g-tag"`
	BaseURL       string `env:"BASE_URL" envDefault:"https://seashells.io/v/"`
	NetCatBinding string `env:"NETCAT_BINDING" envDefault:":1337"`
	WebAppBinding string `env:"WEBAPP_BINDING" envDefault:":8080"`
}
