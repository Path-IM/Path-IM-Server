package xetcd

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/zeromicro/go-zero/core/discov"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"time"
)

func NewClient(config discov.EtcdConf) *clientv3.Client {
	var tlsConf *tls.Config
	if config.HasTLS() {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.CertKeyFile)
		if err != nil {
			panic(err)
		}
		caData, err := ioutil.ReadFile(config.CACertFile)
		if err != nil {
			panic(err)
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)
		tlsConf = &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
			RootCAs:            pool,
			Certificates:       []tls.Certificate{cert},
		}
	}
	cfg := clientv3.Config{
		Endpoints:            config.Hosts,
		AutoSyncInterval:     time.Minute,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    5 * time.Second,
		DialKeepAliveTimeout: 5 * time.Second,
		RejectOldCluster:     true,
		Username:             config.User,
		Password:             config.Pass,
		TLS:                  tlsConf,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		time.Sleep(time.Second * 2)
		return NewClient(config)
		//panic(err)
	}
	return client
}
