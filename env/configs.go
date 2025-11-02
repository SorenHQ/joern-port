package env


import (



	"os"


	"github.com/joho/godotenv"
)


func GetPort() (port string) {
	port = get_env_value("PORT")
	if port == "" {
		return ":9330"
	}
	return
}

func get_env_value(key string) string {
	env := ".env"
	if appEnv := os.Getenv("ENV"); appEnv != "" {
		env = env + "." + appEnv
	}

	err := godotenv.Load(env)
	if err != nil {
		return ""
	}
	return os.Getenv(key)
}

func GetJoernUrl() (url string) {
	url = get_env_value("JOERN_URL")
	if url == "" {
		return "http://localhost:8080"
	}
	return url
}