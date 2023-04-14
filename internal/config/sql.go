package config

type SqlConfig struct {
	Driver   string
	Host     string
	User     string
	Password string
	DbName   string
	SslMode  string
}
