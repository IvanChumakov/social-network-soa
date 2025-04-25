package config

type Config struct {
	ServerConfig
	PostgresConfig
	KafkaConfig
}

type ServerConfig struct {
	ServAddr string
}

type PostgresConfig struct {
	PostgresDb       string
	PostgresUser     string
	PostgresPassword string
	PostgresPort     int
}

type KafkaConfig struct {
	KafkaUrl  string
	KafkaPort string
}

func NewConfig() *Config {
	return &Config{
		ServerConfig: ServerConfig{
			ServAddr: ":50051",
		},
		PostgresConfig: PostgresConfig{
			PostgresUser:     "user",
			PostgresPassword: "password",
			PostgresPort:     5432,
			PostgresDb:       "posts-db",
		},
		KafkaConfig: KafkaConfig{
			KafkaUrl:  "kafka",
			KafkaPort: ":19092",
		},
	}
}
