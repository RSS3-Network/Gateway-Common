package control

import (
	"context"
	"fmt"
)

func (s *StateClientWriter) CreateKey(ctx context.Context, account string, keyID string, key string) error {
	stateKey := fmt.Sprintf(TemplateValidKey, key)

	if _, err := s.etcdClient.Put(ctx, stateKey, fmt.Sprintf("%s:%s", account, keyID)); err != nil {
		return fmt.Errorf("etcd put %s: %w", stateKey, err)
	}

	return nil
}

func (s *StateClientWriter) DeleteKey(ctx context.Context, key string) error {
	stateKey := fmt.Sprintf(TemplateValidKey, key)

	if _, err := s.etcdClient.Delete(ctx, stateKey); err != nil {
		return fmt.Errorf("etcd put %s: %w", stateKey, err)
	}

	return nil
}

func (s *StateClientWriter) PauseAccount(ctx context.Context, account string) error {
	stateKey := fmt.Sprintf(TemplatePausedAccount, account)

	if _, err := s.etcdClient.Put(ctx, stateKey, ""); err != nil {
		return fmt.Errorf("etcd put %s: %w", stateKey, err)
	}

	return nil
}

func (s *StateClientWriter) ResumeAccount(ctx context.Context, account string) error {
	stateKey := fmt.Sprintf(TemplatePausedAccount, account)

	if _, err := s.etcdClient.Delete(ctx, stateKey); err != nil {
		return fmt.Errorf("etcd put %s: %w", stateKey, err)
	}

	return nil
}

type StateClientWriter struct {
	StateClientBase
}

func NewWriter(endpoints []string) (*StateClientWriter, error) {
	base, err := newBase(endpoints)
	if err != nil {
		return nil, err
	}

	return &StateClientWriter{
		*base,
	}, nil
}
