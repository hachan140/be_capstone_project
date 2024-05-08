package common_configs

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Shopify/sarama"
	"os"
	"path/filepath"
	"time"
)

func BuildSaramaConfig(config KafkaSaramaConfig) (*sarama.Config, error) {
	samaraConfig := sarama.NewConfig()
	samaraConfig.Version = sarama.V1_1_0_0
	samaraConfig.Producer.Partitioner = sarama.NewHashPartitioner
	samaraConfig.Producer.RequiredAcks = sarama.WaitForAll
	samaraConfig.Producer.Return.Successes = true
	samaraConfig.Producer.Flush.Frequency = 1 * time.Second
	samaraConfig.Producer.Retry.Max = config.MaxRetry

	if config.EnableTLS {
		samaraConfig.Net.TLS.Enable = config.EnableTLS
		tlsConfig, err := newTLSConfig(
			config.InsecureSkipVerify,
			config.AuthSDKCertFile,
			config.AuthSDKKeyFile,
			config.AuthSDKCACertFile,
		)
		if err != nil {
			return nil, err
		}

		samaraConfig.Net.TLS.Config = tlsConfig
	}

	return samaraConfig, nil
}

func newTLSConfig(insecureSkipVerify bool, clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := &tls.Config{}
	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	tlsConfig.InsecureSkipVerify = insecureSkipVerify
	if !insecureSkipVerify {
		caCert, err := os.ReadFile(filepath.Clean(caCertFile))
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}
