package vault

import (
	"testing"
)

const testPassword = "super-secret-test-password"

func TestNewKeepassVaultCreatesFile(t *testing.T) {

	filePath := t.TempDir() + "/test.kdbx"

	if err := NewKeepassVault(filePath, testPassword); err != nil {
		t.Fatalf("NewKeepassVault() returned error: %v", err)
	}

	v, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() after creation returned error: %v", err)
	}
	defer func() {
		if err := v.Close(); err != nil {
			t.Fatalf("Close() returned error: %v", err)
		}
	}()
}

func TestOpenKeepassVaultWithWrongPasswordFails(t *testing.T) {

	filePath := t.TempDir() + "/test.kdbx"

	if err := NewKeepassVault(filePath, testPassword); err != nil {
		t.Fatalf("NewKeepassVault() returned error: %v", err)
	}

	_, err := OpenKeepassVault(filePath, "wrong-password")
	if err == nil {
		t.Fatal("OpenKeepassVault() with wrong password expected error, got nil")
	}
}

func TestSetAndGetEntry(t *testing.T) {

	filePath := t.TempDir() + "/test.kdbx"

	if err := NewKeepassVault(filePath, testPassword); err != nil {
		t.Fatalf("NewKeepassVault() returned error: %v", err)
	}

	v, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() returned error: %v", err)
	}
	defer func() {
		if err := v.Close(); err != nil {
			t.Fatalf("Close() returned error: %v", err)
		}
	}()

	const key = "api-key"
	const value = "123456789"

	if err := v.SetEntry(key, value); err != nil {
		t.Fatalf("SetEntry() returned error: %v", err)
	}

	got, err := v.GetEntry(key)
	if err != nil {
		t.Fatalf("GetEntry() returned error: %v", err)
	}

	if got != value {
		t.Fatalf("GetEntry() = %q, want %q", got, value)
	}
}

func TestGetEntryForMissingKeyReturnsError(t *testing.T) {

	filePath := t.TempDir() + "/test.kdbx"

	if err := NewKeepassVault(filePath, testPassword); err != nil {
		t.Fatalf("NewKeepassVault() returned error: %v", err)
	}

	v, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() returned error: %v", err)
	}
	defer func() {
		if err := v.Close(); err != nil {
			t.Fatalf("Close() returned error: %v", err)
		}
	}()

	_, err = v.GetEntry("does-not-exist")
	if err == nil {
		t.Fatal("GetEntry() for missing key expected error, got nil")
	}
}

func TestSetEntryUpdatesExistingValue(t *testing.T) {

	filePath := t.TempDir() + "/test.kdbx"

	if err := NewKeepassVault(filePath, testPassword); err != nil {
		t.Fatalf("NewKeepassVault() returned error: %v", err)
	}

	v, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() returned error: %v", err)
	}
	defer func() {
		if err := v.Close(); err != nil {
			t.Fatalf("Close() returned error: %v", err)
		}
	}()

	const key = "token"

	if err := v.SetEntry(key, "old-value"); err != nil {
		t.Fatalf("SetEntry(old) returned error: %v", err)
	}

	if err := v.SetEntry(key, "new-value"); err != nil {
		t.Fatalf("SetEntry(new) returned error: %v", err)
	}

	got, err := v.GetEntry(key)
	if err != nil {
		t.Fatalf("GetEntry() returned error: %v", err)
	}

	if got != "new-value" {
		t.Fatalf("GetEntry() = %q, want %q", got, "new-value")
	}
}

func TestEntryLifecycleAcrossMultipleOpens(t *testing.T) {
	filePath := t.TempDir() + "/test.kdbx"

	if err := NewKeepassVault(filePath, testPassword); err != nil {
		t.Fatalf("NewKeepassVault() returned error: %v", err)
	}

	const key = "service-token"
	const initialValue = "initial-secret"
	const updatedValue = "updated-secret"

	// Erstes Öffnen: Wert schreiben
	v1, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() first open returned error: %v", err)
	}

	if err := v1.SetEntry(key, initialValue); err != nil {
		t.Fatalf("SetEntry() initial write returned error: %v", err)
	}

	if err := v1.Close(); err != nil {
		t.Fatalf("Close() after initial write returned error: %v", err)
	}

	// Zweites Öffnen: Wert lesen und ändern
	v2, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() second open returned error: %v", err)
	}

	got, err := v2.GetEntry(key)
	if err != nil {
		t.Fatalf("GetEntry() after reopen returned error: %v", err)
	}

	if got != initialValue {
		t.Fatalf("GetEntry() after first reopen = %q, want %q", got, initialValue)
	}

	if err := v2.SetEntry(key, updatedValue); err != nil {
		t.Fatalf("SetEntry() update returned error: %v", err)
	}

	got, err = v2.GetEntry(key)
	if err != nil {
		t.Fatalf("GetEntry() after update returned error: %v", err)
	}

	if got != updatedValue {
		t.Fatalf("GetEntry() after update = %q, want %q", got, updatedValue)
	}

	if err := v2.Close(); err != nil {
		t.Fatalf("Close() after update returned error: %v", err)
	}

	// Drittes Öffnen: geänderten Wert erneut lesen
	v3, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() third open returned error: %v", err)
	}
	defer func() {
		if err := v3.Close(); err != nil {
			t.Fatalf("Close() on third open returned error: %v", err)
		}
	}()

	got, err = v3.GetEntry(key)
	if err != nil {
		t.Fatalf("GetEntry() after second reopen returned error: %v", err)
	}

	if got != updatedValue {
		t.Fatalf("GetEntry() after second reopen = %q, want %q", got, updatedValue)
	}
}

func TestClosePersistsData(t *testing.T) {

	filePath := t.TempDir() + "/test.kdbx"

	if err := NewKeepassVault(filePath, testPassword); err != nil {
		t.Fatalf("NewKeepassVault() returned error: %v", err)
	}

	v, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() returned error: %v", err)
	}

	if err := v.SetEntry("username", "alice"); err != nil {
		t.Fatalf("SetEntry() returned error: %v", err)
	}

	if err := v.Close(); err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}

	v2, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() after reopen returned error: %v", err)
	}
	defer func() {
		if err := v2.Close(); err != nil {
			t.Fatalf("Close() on reopened vault returned error: %v", err)
		}
	}()

	got, err := v2.GetEntry("username")
	if err != nil {
		t.Fatalf("GetEntry() after reopen returned error: %v", err)
	}

	if got != "alice" {
		t.Fatalf("GetEntry() after reopen = %q, want %q", got, "alice")
	}
}

func TestSetDescriptionForExistingEntry(t *testing.T) {

	filePath := t.TempDir() + "/test.kdbx"

	if err := NewKeepassVault(filePath, testPassword); err != nil {
		t.Fatalf("NewKeepassVault() returned error: %v", err)
	}

	v, err := OpenKeepassVault(filePath, testPassword)
	if err != nil {
		t.Fatalf("OpenKeepassVault() returned error: %v", err)
	}
	defer func() {
		if err := v.Close(); err != nil {
			t.Fatalf("Close() returned error: %v", err)
		}
	}()

	if err := v.SetEntry("service-a", "secret"); err != nil {
		t.Fatalf("SetEntry() returned error: %v", err)
	}

	if err := v.SetDescription("service-a", "my test description"); err != nil {
		t.Fatalf("SetDescription() returned error: %v", err)
	}
}
