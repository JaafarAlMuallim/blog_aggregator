package config

type Config struct {
	DbUrl string `json:"db_url"`
	User  string `json:"current_user_name"`
}
