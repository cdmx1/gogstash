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

func Test_output_sentry_module(t *testing.T) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		TracesSampleRate: 1.0,
	})
	require.NoError(t, err)

	hubInfo := sentry.CurrentHub().Clone()
	hubInfo.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelInfo)
	})
	hubWarn := sentry.CurrentHub().Clone()
	hubWarn.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelWarning)
	})
	fmt.Println("Test output sentry module")
	for i := 0; i < 60; i++ {
		i := i
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			switch {
			case i%2 == 0:
				hubInfo.CaptureMessage(fmt.Sprintf("Hub Info %d", i))
			case i%2 == 1:
				hubWarn.CaptureMessage(fmt.Sprintf("Hub Warn %d", i))
			}
		})
	}

	t.Cleanup(func() {
		hubInfo.Flush(5 * time.Second)
		hubWarn.Flush(5 * time.Second)
	})
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
