package outputsentry

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/tsaikd/gogstash/config"
	"github.com/tsaikd/gogstash/config/logevent"
)

const ModuleName = "sentry"

// OutputConfig holds the configuration json fields and internal objects
type OutputConfig struct {
	config.OutputConfig
	DSN string `json:"dsn"`
}

// DefaultOutputConfig returns an OutputConfig struct with default values
func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		OutputConfig: config.OutputConfig{
			CommonConfig: config.CommonConfig{
				Type: ModuleName,
			},
		},
		DSN: "",
	}
}

func InitHandler(
	ctx context.Context,
	raw config.ConfigRaw,
	control config.Control,
) (config.TypeOutputConfig, error) {
	conf := DefaultOutputConfig()

	if err := config.ReflectConfig(raw, &conf); err != nil {
		return nil, err
	}

	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              conf.DSN,
		TracesSampleRate: 1.0,
		Transport:        sentrySyncTransport,
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	return &conf, nil
}

func (o *OutputConfig) Output(ctx context.Context, event logevent.LogEvent) (err error) {
	return
}
