package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lithictech/go-aperitif/convext"
	"os"
)

var BuildTime = "1970-01-01T00:00:00Z"
var BuildSha = "0000000000000000000000000000000000000000"

const Version = "0.1.2"
const Repo = "lithictech/webhookdb-cli"

type Config struct {
	// ApiHost is the URL of the service, like
	// https://api.production.webhookdb.com, or http://localhost:1234.
	ApiHost string
	// Use a non-empty environment variable value to enable debug mode,
	// which uses debug-level logging and may change other behaviors.
	Debug     bool
	LogFile   string
	LogFormat string
	LogLevel  string
	// PrefsNamespace is used to namespace different environments
	// in the .webhookdb prefs file.
	// It defaults to API_HOST but you can set it to something else
	// so multiple api hosts can use the same prefs,
	// like if they are backed by the same DB.
	PrefsNamespace string
	SentryDsn      string
	WebsiteHost    string
}

func LoadConfig(filenames ...string) Config {
	_ = godotenv.Overload(filenames...)
	cfg := Config{
		ApiHost:        MustEnvStr("WEBHOOKDB_API_HOST"),
		Debug:          os.Getenv("WEBHOOKDB_DEBUG") != "",
		LogFile:        os.Getenv("WEBHOOKDB_LOG_FILE"),
		LogFormat:      os.Getenv("WEBHOOKDB_LOG_FORMAT"),
		LogLevel:       MustEnvStr("WEBHOOKDB_LOG_LEVEL"),
		PrefsNamespace: os.Getenv("WEBHOOKDB_PREFS_NAMESPACE"),
		SentryDsn:      os.Getenv("WEBHOOKDB_SENTRY_DSN"),
		WebsiteHost:    MustEnvStr("WEBHOOKDB_WEBSITE_HOST"),
	}
	if cfg.PrefsNamespace == "" {
		cfg.PrefsNamespace = cfg.ApiHost
	}
	if cfg.Debug {
		cfg.LogLevel = "debug"
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
	MustSetEnv("WEBHOOKDB_LOG_LEVEL", "error")
	MustSetEnv("WEBHOOKDB_API_HOST", "https://api.production.webhookdb.com")
	MustSetEnv("WEBHOOKDB_WEBSITE_HOST", "https://webhookdb.com")
}
