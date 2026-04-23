package configloader

import (
	"os"
	"testing"
)

// --- argParse ---

func TestArgParse_NoArgs_ReturnsError(t *testing.T) {
	_, err := argParse([]string{})
	if err == nil {
		t.Fatal("expected error for empty args, got nil")
	}
}

func TestArgParse_BaseModuleOnly(t *testing.T) {
	result, err := argParse([]string{"run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BaseModule != "run" {
		t.Errorf("expected BaseModule 'run', got '%s'", result.BaseModule)
	}
	if len(result.Args) != 0 {
		t.Errorf("expected no args, got %d", len(result.Args))
	}
}

func TestArgParse_StringFlagWithValue(t *testing.T) {
	result, err := argParse([]string{"run", "--sourceDir", "/tmp/data"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(result.Args))
	}
	if result.Args[0].Key != "sourceDir" || result.Args[0].Value != "/tmp/data" {
		t.Errorf("unexpected arg: %+v", result.Args[0])
	}
}

func TestArgParse_BoolFlagWithoutValue(t *testing.T) {
	result, err := argParse([]string{"help", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(result.Args))
	}
	if result.Args[0].Key != "help" || result.Args[0].Value != "" {
		t.Errorf("unexpected arg: %+v", result.Args[0])
	}
}

func TestArgParse_MultipleFlags(t *testing.T) {
	result, err := argParse([]string{"run", "--sourceDir", "/src", "--dbPath", "/db.kdbx"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(result.Args))
	}
}

func TestArgParse_UnexpectedValueWithoutFlag_ReturnsError(t *testing.T) {
	_, err := argParse([]string{"run", "orphanValue"})
	if err == nil {
		t.Fatal("expected error for value without preceding flag, got nil")
	}
}

// --- loadConfigFromFile ---

func TestLoadConfigFromFile_ValidYAML(t *testing.T) {
	content := []byte("sourceDir: /tmp/src\ndbPath: /tmp/db.kdbx\nfileExtensions:\n  - .txt\n  - .md\n")
	f, err := os.CreateTemp("", "confy-test-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.Write(content)
	f.Close()

	cfg, err := loadConfigFromFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SourceDir != "/tmp/src" {
		t.Errorf("expected SourceDir '/tmp/src', got '%s'", cfg.SourceDir)
	}
	if cfg.DBPath != "/tmp/db.kdbx" {
		t.Errorf("expected DBPath '/tmp/db.kdbx', got '%s'", cfg.DBPath)
	}
	if len(cfg.FileExtensions) != 2 {
		t.Errorf("expected 2 extensions, got %d", len(cfg.FileExtensions))
	}
}

func TestLoadConfigFromFile_NonExistentFile_ReturnsError(t *testing.T) {
	_, err := loadConfigFromFile("/nonexistent/path/config.yml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfigFromFile_InvalidYAML_ReturnsError(t *testing.T) {
	f, err := os.CreateTemp("", "confy-test-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.WriteString(":: this: is: not: valid: yaml: [")
	f.Close()

	_, err = loadConfigFromFile(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

// --- applyArg ---

func TestApplyArg_StringField(t *testing.T) {
	cfg := &Config{}
	err := applyArg(cfg, CommandLineArg{Key: "sourceDir", Value: "/my/src"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SourceDir != "/my/src" {
		t.Errorf("expected SourceDir '/my/src', got '%s'", cfg.SourceDir)
	}
}

func TestApplyArg_BoolField_SetTrue(t *testing.T) {
	cfg := &Config{}
	err := applyArg(cfg, CommandLineArg{Key: "help", Value: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.HelpOnly {
		t.Error("expected HelpOnly to be true")
	}
}

func TestApplyArg_BoolField_ExplicitFalse(t *testing.T) {
	cfg := &Config{HelpOnly: true}
	err := applyArg(cfg, CommandLineArg{Key: "help", Value: "false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HelpOnly {
		t.Error("expected HelpOnly to be false")
	}
}

func TestApplyArg_SliceField_AppendsValue(t *testing.T) {
	cfg := &Config{}
	applyArg(cfg, CommandLineArg{Key: "fileExtensions", Value: ".txt"})
	applyArg(cfg, CommandLineArg{Key: "fileExtensions", Value: ".md"})
	if len(cfg.FileExtensions) != 2 {
		t.Errorf("expected 2 extensions, got %d", len(cfg.FileExtensions))
	}
	if cfg.FileExtensions[0] != ".txt" || cfg.FileExtensions[1] != ".md" {
		t.Errorf("unexpected extensions: %v", cfg.FileExtensions)
	}
}

func TestApplyArg_UnknownKey_ReturnsError(t *testing.T) {
	cfg := &Config{}
	err := applyArg(cfg, CommandLineArg{Key: "nonExistentField", Value: "x"})
	if err == nil {
		t.Fatal("expected error for unknown key, got nil")
	}
}

// --- file + CLI arg combination ---

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "confy-merge-test-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestMerge_CLIArgOverridesFileValue(t *testing.T) {
	path := writeTempConfig(t, "sourceDir: /from/file\ndbPath: /db.kdbx\n")

	cfg, err := loadConfigFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if err := applyArg(cfg, CommandLineArg{Key: "sourceDir", Value: "/from/cli"}); err != nil {
		t.Fatalf("unexpected error applying arg: %v", err)
	}

	if cfg.SourceDir != "/from/cli" {
		t.Errorf("expected SourceDir '/from/cli', got '%s'", cfg.SourceDir)
	}
	// unrelated file value must remain intact
	if cfg.DBPath != "/db.kdbx" {
		t.Errorf("expected DBPath '/db.kdbx', got '%s'", cfg.DBPath)
	}
}

func TestMerge_CLIArgAddsValueNotInFile(t *testing.T) {
	path := writeTempConfig(t, "sourceDir: /from/file\n")

	cfg, err := loadConfigFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if err := applyArg(cfg, CommandLineArg{Key: "dbPath", Value: "/cli/db.kdbx"}); err != nil {
		t.Fatalf("unexpected error applying arg: %v", err)
	}

	if cfg.DBPath != "/cli/db.kdbx" {
		t.Errorf("expected DBPath '/cli/db.kdbx', got '%s'", cfg.DBPath)
	}
	if cfg.SourceDir != "/from/file" {
		t.Errorf("expected SourceDir '/from/file', got '%s'", cfg.SourceDir)
	}
}

func TestMerge_CLIArgAppendsToSliceFromFile(t *testing.T) {
	path := writeTempConfig(t, "fileExtensions:\n  - .txt\n")

	cfg, err := loadConfigFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if err := applyArg(cfg, CommandLineArg{Key: "fileExtensions", Value: ".md"}); err != nil {
		t.Fatalf("unexpected error applying arg: %v", err)
	}

	if len(cfg.FileExtensions) != 2 {
		t.Fatalf("expected 2 extensions, got %d", len(cfg.FileExtensions))
	}
	if cfg.FileExtensions[0] != ".txt" || cfg.FileExtensions[1] != ".md" {
		t.Errorf("unexpected extensions: %v", cfg.FileExtensions)
	}
}

func TestMerge_FileValuesPreservedWhenNoOverride(t *testing.T) {
	path := writeTempConfig(t, "sourceDir: /src\ndbPath: /db.kdbx\nfileExtensions:\n  - .go\n")

	cfg, err := loadConfigFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	// apply an unrelated arg — file values must be unchanged
	if err := applyArg(cfg, CommandLineArg{Key: "help", Value: ""}); err != nil {
		t.Fatalf("unexpected error applying arg: %v", err)
	}

	if cfg.SourceDir != "/src" || cfg.DBPath != "/db.kdbx" || len(cfg.FileExtensions) != 1 {
		t.Errorf("file values were unexpectedly modified: %+v", cfg)
	}
}
