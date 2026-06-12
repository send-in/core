package config

type ServerConfig struct {
	Port 	string
	Passkey string
}

type DatabaseConfig struct {
	SSL 	 string
	Name 	 string
	Host 	 string
	Port 	 string
	Password string
	Username string
}

type AuthConfig struct {
	ClientID     	 string
	CallbackURL  	 string
	ClientSecret 	 string
	SessionSecret 	 string
	EncryptionSecret string
}

type Config struct {
	Auth 	 AuthConfig
	Server   ServerConfig
	Database DatabaseConfig
}