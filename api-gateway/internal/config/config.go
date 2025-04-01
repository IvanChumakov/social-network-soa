package config

type Config struct {
	Port     string
	GrpcAddr string
}

func NewConfig() *Config {
	return &Config{
		Port:     ":8080",
		GrpcAddr: ":50051",
	}
}
