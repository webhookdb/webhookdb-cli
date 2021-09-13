package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lithictech/go-aperitif/convext"
	"os"
)

var BuildTime = "btime"
var BuildSha = "bsha"

type Config struct {
	ApiHost   string
	Debug     bool
	LogFile   string
	LogFormat string
	LogLevel  string
}

func LoadConfig(filenames ...string) Config {
	_ = godotenv.Overload(filenames...)
	cfg := Config{
		ApiHost:   os.Getenv("API_HOST"),
		Debug:     os.Getenv("DEBUG") != "",
		LogFile:   os.Getenv("LOG_FILE"),
		LogFormat: os.Getenv("LOG_FORMAT"),
		LogLevel:  MustEnvStr("LOG_LEVEL"),
	}
	return cfg
}

func MustEnvStr(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(fmt.Sprintf("'%s' should have had a default set, something weird happened", k))
	}
	return v
}

func MustSetEnv(k string, v interface{}) {
	if _, ok := os.LookupEnv(k); !ok {
		convext.Must(os.Setenv(k, fmt.Sprintf("%v", v)))
	}
}

func init() {
	MustSetEnv("LOG_LEVEL", "warn")
	MustSetEnv("API_HOST", "https://webhookdb-api-production.herokuapp.com/")
}
