package control

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
)

func (s *StateClientReader) CheckKey(ctx context.Context, key string) (*string, *string, error) {
	stateKey := fmt.Sprintf(TemplateValidKey, key)

	stateValue, err := s.etcdClient.Get(ctx, stateKey)

	if err != nil {
		if errors.Is(err, rpctypes.ErrKeyNotFound) {
			return nil, nil, nil // Doesn't exist
		}

		// else
		return nil, nil, fmt.Errorf("etcd get %s: %w", stateKey, err)
	}

	for _, kv := range stateValue.Kvs {
		if string(kv.Key) == stateKey {
			splits := strings.SplitN(string(kv.Value), ":", 2)
			if len(splits) < 2 {
				return nil, nil, fmt.Errorf("invalid key info: %s", kv.Value)
			}

			// else
			return &splits[0], &splits[1], nil // address, id
		}
	}

	return nil, nil, nil // No match key
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

func NewReader(endpoints []string, username *string, password *string) (*StateClientReader, error) {
	base, err := newBase(endpoints, username, password)
	if err != nil {
		return nil, err
	}

	return &StateClientReader{
		*base,
	}, nil
}
