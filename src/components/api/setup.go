package api

import (
	"github.com/bsphere/le_go"
	"github.com/kelseyhightower/envconfig"
	"github.com/ml-tv/tv-api/src/services/tmdb"
	"github.com/ml-tv/tv-api/src/core/logger"
	"github.com/ml-tv/tv-api/src/core/notifiers/mailer"
	"github.com/ml-tv/tv-api/src/core/storage/db"
)

// Args represents the app args
type Args struct {
	Port            string `default:"5000"`
	PostgresURI     string `required:"true" envconfig:"postgres_uri"`
	TMDbAPIKey      string `envconfig:"tmdb_api_key"`
	LogEntriesToken string `envconfig:"logentries_token"`
	EmailAPIKey     string `envconfig:"email_api_key"`
	EmailFrom       string `envconfig:"email_default_from"`
	EmailTo         string `envconfig:"email_default_to"`
	Debug           bool   `default:"false"`
}

// Setup parses the env, sets the app globals and returns the params
func Setup() *Args {
	var params Args
	if err := envconfig.Process("", &params); err != nil {
		panic(err)
	}

	if err := db.Setup(params.PostgresURI); err != nil {
		panic(err)
	}

	// LogEntries
	if params.LogEntriesToken != "" {
		le, err := le_go.Connect(params.LogEntriesToken)
		if err != nil {
			panic(err)
		}
		logger.LogEntries = le
	}

	// Sendgrid
	if params.EmailAPIKey != "" {
		mailer.Emailer = mailer.NewMailer(params.EmailAPIKey, params.EmailFrom, params.EmailTo)
	}

	// TMDb
	tmdb.APIKey = params.TMDbAPIKey

	return &params
}

// SetupIfNeeded parses the env, and sets the app globals
func SetupIfNeeded() {
	if db.Writer != nil {
		return
	}

	Setup()
}
