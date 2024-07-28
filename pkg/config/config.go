package config

type appConfig struct {
	ListenHttp  string `toml:"listen_http"`
	ListenHttps string `toml:"listen_https"`
}

type logConfig struct {
	Level string `toml:"level"`
	Size  string `toml:"size"`
	Dir   string `toml:"dir"`
}

type MySQLConfig struct {
	Debug bool   `toml:"debug"`
	Log   string `toml:"log"`
	Link  string `toml:"link"`
}

type ClickHouseConfig struct {
	Debug bool   `toml:"debug"`
	Log   string `toml:"log"`
	Link  string `toml:"link"`
}

type DatabaseConfig struct {
	MySQL      []MySQLConfig      `toml:"mysql"`
	ClickHouse []ClickHouseConfig `toml:"clickhouse"`
}
