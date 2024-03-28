package control_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rss3-network/gateway-common/control"
)

func TestControl(t *testing.T) {
	t.Parallel()

	// Prepare configs
	etcdEndpoints := []string{"localhost:2379"}

	// Create writer client
	writer, err := control.NewWriter(etcdEndpoints, nil, nil)

	if err != nil {
		t.Fatal(fmt.Errorf("create writer: %w", err))
	}

	defer writer.Stop()

	// Create reader client
	reader, err := control.NewReader(etcdEndpoints, nil, nil)

	if err != nil {
		t.Fatal(fmt.Errorf("create reader: %w", err))
	}

	defer reader.Stop()

	// Prepare a key
	demoKey := "84b01bc1-4dad-4694-99ce-514c37b88f9a"
	demoKeyID := "83298679882397449"
	// ...and an accounts
	demoAccount := "0xD3E8ce4841ed658Ec8dcb99B7a74beFC377253EA"
	// ...and context
	ctx := context.Background()

	// Test flow

	state0(ctx, t, reader, writer, demoKey, demoKeyID, demoAccount)

	state1(ctx, t, reader, writer, demoKey, demoKeyID, demoAccount)

	state2(ctx, t, reader, writer, demoKey, demoKeyID, demoAccount)

	state3(ctx, t, reader, writer, demoKey, demoKeyID, demoAccount)

	state4(ctx, t, reader, writer, demoKey, demoKeyID, demoAccount)
}

func state0(ctx context.Context, t *testing.T, reader *control.StateClientReader, _ *control.StateClientWriter, demoKey string, _ string, demoAccount string) {
	//// State 0: nothing
	var (
		account *string
		keyID   *string
		exist   bool
		err     error
	)

	if account, keyID, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account != nil || keyID != nil {
		t.Error("key should not exist")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if exist {
		t.Error("account should not exist")
	}
}

func state1(ctx context.Context, t *testing.T, reader *control.StateClientReader, writer *control.StateClientWriter, demoKey string, demoKeyID string, demoAccount string) {
	//// State 1: create key, no paused account
	var (
		account *string
		keyID   *string
		exist   bool
		err     error
	)

	if err = writer.CreateKey(ctx, demoAccount, demoKeyID, demoKey); err != nil {
		t.Error(fmt.Errorf("create key: %w", err))
	}

	if account, keyID, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account == nil || keyID == nil {
		t.Error("key should exist")
	} else if *account != demoAccount || *keyID != demoKeyID {
		t.Error("wrong key info")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if exist {
		t.Error("account should not exist")
	}
}

func state2(ctx context.Context, t *testing.T, reader *control.StateClientReader, writer *control.StateClientWriter, demoKey string, demoKeyID string, demoAccount string) {
	// State 2: create key, pause account
	var (
		account *string
		keyID   *string
		exist   bool
		err     error
	)

	if err = writer.PauseAccount(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("pause account: %w", err))
	}

	if account, keyID, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account == nil || keyID == nil {
		t.Error("key should exist")
	} else if *account != demoAccount || *keyID != demoKeyID {
		t.Error("wrong key info")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if !exist {
		t.Error("account should exist")
	}
}

func state3(ctx context.Context, t *testing.T, reader *control.StateClientReader, writer *control.StateClientWriter, demoKey string, _ string, demoAccount string) {
	// State 3: no key, pause account
	var (
		account *string
		keyID   *string
		exist   bool
		err     error
	)

	if err = writer.DeleteKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("delete key: %w", err))
	}

	if account, keyID, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account != nil || keyID != nil {
		t.Error("key should not exist")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if !exist {
		t.Error("account should exist")
	}
}

func state4(ctx context.Context, t *testing.T, reader *control.StateClientReader, writer *control.StateClientWriter, demoKey string, _ string, demoAccount string) {
	// State 4: nothing again
	var (
		account *string
		keyID   *string
		exist   bool
		err     error
	)

	if err = writer.ResumeAccount(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("resume account: %w", err))
	}

	if account, keyID, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account != nil || keyID != nil {
		t.Error("key should not exist")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if exist {
		t.Error("account should not exist")
	}
}
