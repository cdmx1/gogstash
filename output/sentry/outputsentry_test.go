package outputsentry

import (
	"fmt"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/gogstash/config"
	"github.com/tsaikd/gogstash/config/goglog"
)

func init() {
	goglog.Logger.SetLevel(logrus.DebugLevel)
	config.RegistOutputHandler(ModuleName, InitHandler)
}

func captureWithSentry(level sentry.Level, format string, args ...any) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)
	})
	hub.CaptureMessage(fmt.Sprintf(format, args...))
	hub.Flush(time.Second * 3)
}

func Test_output_sentry_module(t *testing.T) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		TracesSampleRate: 1.0,
	})
	require.NoError(t, err)

	for i := 0; i < 60; i++ {
		i := i
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			switch {
			case i%2 == 0:
				captureWithSentry(sentry.LevelInfo, "Hub Info %d", i)
			case i%2 == 1:
				captureWithSentry(sentry.LevelWarning, "Hub warn %d", i)
			}
		})
	}
}

func TestSentryRecover(t *testing.T) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "", // Use an actual DSN here
		TracesSampleRate: 1.0,
	})
	require.NoError(t, err)

	require.NotPanics(t, func() {
		defer func() {
			if err := recover(); err != nil {
				// Report the panic to Sentry
				sentry.CurrentHub().Recover(err)
				// Flush the buffered events to ensure the panic is sent to Sentry
				sentry.Flush(time.Second * 5)
			}
		}()
		// Code that causes a panic
		panic("panic")
	})
}
