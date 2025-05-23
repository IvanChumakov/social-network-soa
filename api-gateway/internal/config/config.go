package config

type Config struct {
	Port          string
	GrpcAddr      string
	StatsGrpcAddr string
}

func NewConfig() *Config {
	return &Config{
		Port:          ":8080",
		GrpcAddr:      ":50051",
		StatsGrpcAddr: ":50052",
	}
}
