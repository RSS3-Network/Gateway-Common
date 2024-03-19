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
	writer, err := control.NewWriter(etcdEndpoints)

	if err != nil {
		t.Fatal(fmt.Errorf("create writer: %w", err))
	}

	defer writer.Stop()

	// Create reader client
	reader, err := control.NewReader(etcdEndpoints)

	if err != nil {
		t.Fatal(fmt.Errorf("create reader: %w", err))
	}

	defer reader.Stop()

	// Prepare a key
	demoKey := "84b01bc1-4dad-4694-99ce-514c37b88f9a"
	// ...and an accounts
	demoAccount := "0xD3E8ce4841ed658Ec8dcb99B7a74beFC377253EA"
	// ...and context
	ctx := context.Background()

	// Test flow
	var (
		account *string
		exist   bool
	)
	//// State 0: nothing
	if account, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account != nil {
		t.Error("key should not exist")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if exist {
		t.Error("account should not exist")
	}

	//// State 1: create key, no paused account
	if err = writer.CreateKey(ctx, demoKey, demoAccount); err != nil {
		t.Error(fmt.Errorf("create key: %w", err))
	}

	if account, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account == nil {
		t.Error("key should exist")
	} else if *account != demoAccount {
		t.Error("wrong account")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if exist {
		t.Error("account should not exist")
	}

	// State 2: create key, pause account

	if err = writer.PauseAccount(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("pause account: %w", err))
	}

	if account, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account == nil {
		t.Error("key should exist")
	} else if *account != demoAccount {
		t.Error("wrong account")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if !exist {
		t.Error("account should exist")
	}

	// State 3: no key, pause account
	if err = writer.DeleteKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("delete key: %w", err))
	}

	if account, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account != nil {
		t.Error("key should not exist")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if !exist {
		t.Error("account should exist")
	}

	// State 4: nothing again
	if err = writer.ResumeAccount(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("resume account: %w", err))
	}

	if account, err = reader.CheckKey(ctx, demoKey); err != nil {
		t.Error(fmt.Errorf("check key: %w", err))
	} else if account != nil {
		t.Error("key should not exist")
	}

	if exist, err = reader.CheckAccountPaused(ctx, demoAccount); err != nil {
		t.Error(fmt.Errorf("check account: %w", err))
	} else if exist {
		t.Error("account should not exist")
	}
}
