package config

type ServerConfig struct {
	Port               int `env:"TODO_PORT"`
	DatabaseDriverName string
	DatabaseFile       string `env:"TODO_DBFILE"`
	RootPassword       string `env:"TODO_PASSWORD"`
	LoggerLvl          string
}
