package config

type Config struct {
	ServAddr         string
	PostgresDb       string
	PostgresUser     string
	PostgresPassword string
	PostgresPort     int
}

func NewConfig() *Config {
	return &Config{
		ServAddr:         ":50051",
		PostgresUser:     "user",
		PostgresPassword: "password",
		PostgresPort:     5432,
		PostgresDb:       "posts-db",
	}
}
