package config

type Config struct {
	ServiceName string
	Server      Server
	KeyCloak    KeyCloak
	Debug       bool
	RunK8s      bool
}

type Server struct {
	HTTPAddr string
	WinSrv   bool
	Debug    bool
}
type KeyCloak struct {
	Realm                     string
	BaseUrl                   string
	KeyCloakClient            KeyCloakClient
	RealmMasterRS256PublicKey string
}

type KeyCloakClient struct {
	ClientId     string
	ClientSecret string
}
