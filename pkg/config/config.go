package config

type appConfig struct {
	ListenHttp  string `toml:"listen_http" yaml:"Listenhttp" json:",default=127.0.0.1:80"`
	ListenHttps string `toml:"listen_https" yaml:"ListenHttps" json:","`
}

type logConfig struct {
	Level string `toml:"level" yaml:"Level" json:",default=info"`
	Size  string `toml:"size" yaml:"Size" json:",default=256m"`
	Dir   string `toml:"dir" yaml:"Dir" json:",default=logs"`
	Gin   string `toml:"gin" yaml:"Gin" json:",optional"`
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
	MySQL      *MySQLConfig      `toml:"mysql"`
	ClickHouse *ClickHouseConfig `toml:"clickhouse"`
}
