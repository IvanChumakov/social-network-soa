package config

type Config struct {
	ServerAddr       string
	PostgresDb       string
	PostgresUser     string
	PostgresPassword string
	PostgresPort     int
}

func NewConfig() *Config {
	return &Config{
		ServerAddr:       ":8081",
		PostgresUser:     "user",
		PostgresPassword: "password",
		PostgresPort:     5432,
		PostgresDb:       "social-network",
	}
}
