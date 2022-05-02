package configs

type Config struct {
	requireConfirmation bool
}

var config Config

func LoadConfig() {
	config = Config{
		requireConfirmation: true,
	}
}

func GetConfig() Config {
	return config
}
