package config

type AppConfig struct {
	Datastore Datastore `json:"datastore"`
	Jwt       Jwt       `json:"jwt"`
}

type Datastore struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Jwt struct {
	Secret string `json:"secret"`
}
