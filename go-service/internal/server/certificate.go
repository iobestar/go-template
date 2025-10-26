package server

import (
	"context"
	"crypto/tls"
	_ "embed"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type CertificateLoader interface {
	GetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error)
	GetClientCertificate(_ *tls.CertificateRequestInfo) (*tls.Certificate, error)
	WatchAndReload(ctx context.Context, wg *sync.WaitGroup) error
	LoadAndUpdateCert() error
}

func NewCertificateLoader(certPath, keyPath string) (CertificateLoader, error) {
	if len(certPath) == 0 {
		return nil, errors.New("cert key path is empty")
	}

	if len(keyPath) == 0 {
		return nil, errors.New("cert key path is empty")
	}
	return newFileCertificateLoader(certPath, keyPath)
}

func newFileCertificateLoader(certPath, keyPath string) (CertificateLoader, error) {
	certificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	return &fileCertificateLoader{
		certPath: certPath,
		keyPath:  keyPath,

		certificate: &certificate,

		lock: &sync.RWMutex{},
	}, nil
}

type fileCertificateLoader struct {
	certPath string
	keyPath  string

	certificate *tls.Certificate

	lock *sync.RWMutex
}

func (c *fileCertificateLoader) LoadAndUpdateCert() error {
	certificate, err := tls.LoadX509KeyPair(c.certPath, c.keyPath)
	if err != nil {
		return err
	}

	c.lock.Lock()
	c.certificate = &certificate
	c.lock.Unlock()
	return nil
}

// WatchAndReload function will watch certificate files and reload certificate when files are changed.
func (c *fileCertificateLoader) WatchAndReload(ctx context.Context, wg *sync.WaitGroup) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := watcher.Add(c.certPath); err != nil {
		return err
	}
	if err := watcher.Add(c.keyPath); err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create {
					if err := c.LoadAndUpdateCert(); nil != err {
						log.WithFields(log.Fields{"keyPath": c.keyPath, "certPath": c.certPath}).WithError(err).Error("error reloading certificate")
					} else {
						log.WithFields(log.Fields{"keyPath": c.keyPath, "certPath": c.certPath}).Info("certificate reloaded")
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					if err := watcher.Add(event.Name); err != nil {
						log.WithField("path", event.Name).WithError(err).Error("error re-watching path")
					}
				}
			case <-ctx.Done():
				if err := watcher.Close(); nil != err {
					log.WithError(err).Error("error closing certificate watcher")
				}
			}
		}
	}()
	log.WithFields(log.Fields{"keyPath": c.keyPath, "certPath": c.certPath}).Info("started certificate watcher")
	return nil
}

func (c *fileCertificateLoader) GetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.certificate, nil
}

func (c *fileCertificateLoader) GetClientCertificate(_ *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.certificate, nil
}
