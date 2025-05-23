package config

type Config struct {
	ClickHouseConfig
	ServerConfig
}

type ClickHouseConfig struct {
	Addr     string
	Port     int
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	ServerAddr string
}

func NewConfig() Config {
	return Config{
		ClickHouseConfig: ClickHouseConfig{
			Addr:     "clickhouse",
			Port:     9000,
			User:     "user",
			Password: "password",
			DBName:   "default",
		},
		ServerConfig: ServerConfig{
			ServerAddr: ":50052",
		},
	}
}
