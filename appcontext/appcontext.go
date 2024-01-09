package appcontext

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lithictech/go-aperitif/logctx"
	"github.com/sirupsen/logrus"
	"github.com/webhookdb/webhookdb-cli/client"
	"github.com/webhookdb/webhookdb-cli/config"
	"github.com/webhookdb/webhookdb-cli/prefs"
	"github.com/webhookdb/webhookdb-cli/whfs"
	"os"
	"runtime"
)

type AppContext struct {
	FS          whfs.FS
	Config      config.Config
	Resty       *resty.Client
	GlobalPrefs *prefs.GlobalPrefs
	Prefs       prefs.Prefs
	Auth        client.Auth
	logger      *logrus.Entry
}

func (ac AppContext) Logger() *logrus.Entry {
	return ac.logger
}

func New(command string, cfg config.Config) (ac AppContext, err error) {
	ac.FS = whfs.New()
	ac.Config = cfg
	if ac.GlobalPrefs, err = prefs.Load(ac.FS); err != nil {
		return
	}
	ac.Prefs = ac.GlobalPrefs.GetNS(cfg.PrefsNamespace)
	ac.Auth = client.Auth{Token: ac.Prefs.AuthToken}
	ac.Resty = newResty(cfg)
	if ac.logger, err = logctx.NewLogger(logctx.NewLoggerInput{
		Level:     cfg.LogLevel,
		Format:    cfg.LogFormat,
		File:      cfg.LogFile,
		BuildSha:  config.BuildSha,
		BuildTime: config.BuildTime,
		Fields:    logrus.Fields{"command": command},
	}); err != nil {
		return
	}
	// TODO: For now, always use stderr instead of stdout because we are running this as a CLI,
	// not an application, and want the caller to be able to collect logs easily.
	// But we should make this better configurable, maybe with a change to logctx.
	ac.logger.Logger.SetOutput(os.Stderr)
	return
}

func (ac AppContext) SavePrefs() error {
	ac.GlobalPrefs.SetNS(ac.Config.PrefsNamespace, ac.Prefs)
	return prefs.Save(ac.FS, ac.GlobalPrefs)
}

func NewTestContext() AppContext {
	cfg := config.LoadConfig()

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	pr := &prefs.GlobalPrefs{}
	ac := AppContext{
		FS:          whfs.New(),
		logger:      logger.WithFields(nil),
		Config:      cfg,
		Resty:       newResty(cfg),
		GlobalPrefs: pr,
		Prefs:       pr.GetNS(""),
	}
	return ac
}

const ctxKey = "appcontext"

func InContext(parent context.Context, ac AppContext) context.Context {
	return context.WithValue(parent, ctxKey, ac)
}

func FromContext(c context.Context) AppContext {
	return c.Value(ctxKey).(AppContext)
}

func newResty(cfg config.Config) *resty.Client {

	r := resty.New().
		SetHostURL(cfg.ApiHost).
		SetHeader("User-Agent", userAgent).
		SetHeader("Whdb-User-Agent", userAgent).
		SetHeader("Whdb-Short-Session", shortSession)
	platformUA := platformUserAgent()
	if platformUA != "" {
		r = r.SetHeader("Whdb-Platform-User-Agent", platformUA)
	}
	r.Debug = cfg.Debug
	return r
}

var userAgent string

func init() {
	userAgent = fmt.Sprintf("WebhookDB/v1 webhookdb-cli/%s (%s; %s) Built/%s https://webhookdb.com",
		config.BuildShaShort, runtime.GOOS, runtime.GOARCH, config.BuildTime)
}
