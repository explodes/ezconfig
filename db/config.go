package db

type DbConfig struct {
	Database Host
}

type Host struct {
	Type     string // only sqlite3, postgres is supported
	Host     string // file, or :memory:, for sqlite3
	Port     int
	DbName   string `toml:"dbname"`
	User     string
	Password string
	Ssl      string
}

func (config *DbConfig) GetDbConfig() *DbConfig {
	return config
}
