package server

import (
	"context"
	"crypto/tls"
	slog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/iobestar/go-template/go-service/internal"
	"github.com/iobestar/go-template/go-service/internal/app"
	"github.com/iobestar/go-template/go-service/internal/logging"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const name = "template"

var _ http.Handler = internal.HttpApp(nil)

func Start(cmd *cobra.Command, version, build string) error {
	defer func() {
		log.Info(name + " stopped")
	}()

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Wait()
	}()

	config := NewConfig()

	disposer := internal.NewDisposer()
	defer func() {
		disposer.Dispose()
	}()

	disposeFunc, err := logging.Init(logging.NewConfig())
	if err != nil {
		return err
	}
	disposer.Attach(disposeFunc)

	logFields := log.Fields{
		"pid":     os.Getpid(),
		"cmd":     cmd.Name(),
		"logging": log.GetLevel().String(),
		"version": version,
		"build":   build,
	}
	if len(viper.ConfigFileUsed()) > 0 {
		logFields["config"] = viper.ConfigFileUsed()
	}
	log.WithFields(logFields).Info("starting " + name)

	a, err := app.NewHttp()
	if nil != err {
		return err
	}

	if err := a.Run(ctx, &wg, disposer); nil != err {
		return err
	}

	httpL, err := createListener(config.Network, config.Address)
	if err != nil {
		return err
	}

	errorWriter := logging.NewLevelWriter(log.ErrorLevel)
	defer func() {
		if err := errorWriter.Close(); nil != err {
			log.WithError(err).Warn("error closing error writer")
		}
	}()

	httpS := &http.Server{
		Handler:  a,
		ErrorLog: slog.New(errorWriter, "", 0),
	}

	httpStopCh := make(chan struct{})
	go func() {
		defer close(httpStopCh)
		if len(config.CertPath) > 0 && len(config.KeyPath) > 0 {
			log.WithFields(log.Fields{"address": config.Address, "cert_path": config.CertPath, "key_path": config.KeyPath}).
				Info("starting https server")
			loader, err := internal.NewCertificateLoader(config.CertPath, config.KeyPath)
			if err != nil {
				log.WithError(err).Error("unable to create certificate loader")
				return
			}
			if err := loader.WatchAndReload(ctx, &wg); nil != err {
				log.WithError(err).Error("error watching and reloading certificates")
				return
			}

			httpS.TLSConfig = &tls.Config{GetCertificate: loader.GetCertificate}

			log.WithFields(log.Fields{"address": config.Address}).Info("https server started")
			if err := httpS.Serve(tls.NewListener(httpL, httpS.TLSConfig)); errors.Is(err, http.ErrServerClosed) {
				log.Info("https server stopped")
			} else if nil != err {
				log.WithError(err).Error("http server error")
			}
		} else {
			log.WithFields(log.Fields{"address": config.Address}).Info("starting http server")
			if err := httpS.Serve(httpL); errors.Is(err, http.ErrServerClosed) {
				log.Info("http server stopped")
			} else if nil != err {
				log.WithError(err).Error("http server error")
			}
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		if err := httpS.Shutdown(context.Background()); nil != err {
			log.WithError(err).Error("http server shutdown error")
		}
		<-httpStopCh
	case <-httpStopCh:
	}
	return nil
}

func createListener(network, address string) (net.Listener, error) {
	if runtime.GOOS == "windows" && network == "unix" {
		return nil, errors.New("unix network not allowed on windows")
	}

	if network == "unix" {
		_ = syscall.Unlink(address)
	}
	return net.Listen(network, address)
}
