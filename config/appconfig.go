package config

type AppConfig struct {
	Datastore Datastore `json:"datastore"`
}

type Datastore struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}
