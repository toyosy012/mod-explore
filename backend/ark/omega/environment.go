package omega

import "github.com/kelseyhightower/envconfig"

type Environments struct {
	DBConfig
	ServerConfig
}

func LoadConfig() (*Environments, error) {
	var cfg Environments
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type DBConfig struct {
	DBUsername   string `envconfig:"DB_USERNAME" required:"true"`
	DBPassword   string `envconfig:"DB_PASSWORD" required:"true"`
	DatabaseName string `envconfig:"DB_DATABASE_NAME" required:"true"`
	Port         uint16 `envconfig:"DB_PORT" default:"42731"`
	DatabaseURL  string `envconfig:"DB_URL"`
}

type ServerConfig struct {
	Address string `envconfig:"ADDRESS" required:"true"`
}
