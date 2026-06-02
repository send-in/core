package config

type ServerConfig struct {
	Port string
	Passkey string
}

type DatabaseConfig struct {
	Name string
	Username string
	Password string
	Host string
	Port string
	SSL string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}