package control

import (
	"context"
	"errors"
	"fmt"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
)

func (s *StateClientReader) CheckKey(ctx context.Context, key string) (*string, error) {
	stateKey := fmt.Sprintf(TemplateValidKey, key)

	stateValue, err := s.etcdClient.Get(ctx, stateKey)

	if err != nil {
		if errors.Is(err, rpctypes.ErrKeyNotFound) {
			return nil, nil // Doesn't exist
		}

		// else
		return nil, fmt.Errorf("etcd get %s: %w", stateKey, err)
	}

	for _, kv := range stateValue.Kvs {
		if string(kv.Key) == stateKey {
			value := string(kv.Value)
			return &value, nil
		}
	}

	return nil, nil // No match key
}

func (s *StateClientReader) CheckAccountPaused(ctx context.Context, account string) (bool, error) {
	stateKey := fmt.Sprintf(TemplatePausedAccount, account)

	stateValue, err := s.etcdClient.Get(ctx, stateKey)
	if err != nil {
		return false, fmt.Errorf("etcd get %s: %w", stateKey, err)
	}

	return stateValue.Count > 0, nil
}

type StateClientReader struct {
	StateClientBase
}

func NewReader(endpoints []string) (*StateClientReader, error) {
	base, err := newBase(endpoints)
	if err != nil {
		return nil, err
	}

	return &StateClientReader{
		*base,
	}, nil
}
