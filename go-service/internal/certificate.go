package internal

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

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

const (
	reloadInterval = 6 * time.Hour
	retryBackoff   = 5 * time.Minute
	maxRetries     = 3
)

// WatchAndReload periodically reloads certificate files every 6 hours.
// On reload failure, it retries with a fixed 5-minute backoff up to 3 times before
// falling back to the regular reload interval.
func (c *fileCertificateLoader) WatchAndReload(ctx context.Context, wg *sync.WaitGroup) error {
	wg.Add(1)
	go func() {
		defer wg.Done()

		reloadTimer := time.NewTimer(reloadInterval)
		defer reloadTimer.Stop()

		logFields := log.Fields{"keyPath": c.keyPath, "certPath": c.certPath}
		retryCount := 0

		for {
			select {
			case <-reloadTimer.C:
				if err := c.LoadAndUpdateCert(); err != nil {
					retryCount++
					if retryCount <= maxRetries {
						log.WithFields(logFields).WithError(err).Errorf("error reloading certificate, retry %d/%d", retryCount, maxRetries)
						reloadTimer.Reset(retryBackoff)
					} else {
						log.WithFields(logFields).WithError(err).Error("error reloading certificate, max retries exceeded")
						retryCount = 0
						reloadTimer.Reset(reloadInterval)
					}
				} else {
					log.WithFields(logFields).Info("certificate reloaded")
					retryCount = 0
					reloadTimer.Reset(reloadInterval)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	log.WithFields(log.Fields{"keyPath": c.keyPath, "certPath": c.certPath}).Info("started certificate reload scheduler")
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
