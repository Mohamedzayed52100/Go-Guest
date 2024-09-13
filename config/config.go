package config

type Config struct {
	Environment        string `env:"ENVIRONMENT,required"`
	DbPostgresHost     string `env:"DB_POSTGRES_HOST,required"`
	DbPostgresUser     string `env:"DB_POSTGRES_USER,required"`
	DbPostgresPassword string `env:"DB_POSTGRES_PASSWORD,required"`
	DbPostgresDb       string `env:"DB_POSTGRES_DB"`
	DbSSLMode          string `env:"DB_SSL_MODE"`
}
