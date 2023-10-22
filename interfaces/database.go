package interfaces

type DbConfig struct {
	Address      string `json:"address"`
	User         string `json:"user"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
}
