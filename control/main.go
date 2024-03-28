package control

import (
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type StateClientBase struct {
	etcdClient *clientv3.Client
}

func newBase(endpoints []string, username *string, password *string) (*StateClientBase, error) {
	// Initialize configurations
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}

	if username != nil {
		cfg.Username = *username
	}

	if password != nil {
		cfg.Password = *password
	}

	// Start client
	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("dial etcd: %w", err)
	}

	return &StateClientBase{
		etcdClient: cli,
	}, nil
}

func (base *StateClientBase) Stop() {
	_ = base.etcdClient.Close() // Ignore errors
}
