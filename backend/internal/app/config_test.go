package app

import "testing"

func TestLoadConfigFromEnvRejectsMissingOrInvalidSettings(t *testing.T) {
	t.Run("missing host", func(t *testing.T) {
		t.Setenv("HOST", "")
		t.Setenv("PORT", "8080")
		t.Setenv("DB_PATH", "/tmp/test.db")

		if _, err := LoadConfigFromEnv(); err == nil {
			t.Fatal("expected error for missing HOST")
		}
	})

	t.Run("missing port", func(t *testing.T) {
		t.Setenv("HOST", "127.0.0.1")
		t.Setenv("PORT", "")
		t.Setenv("DB_PATH", "/tmp/test.db")

		if _, err := LoadConfigFromEnv(); err == nil {
			t.Fatal("expected error for missing PORT")
		}
	})

	t.Run("invalid port", func(t *testing.T) {
		t.Setenv("HOST", "127.0.0.1")
		t.Setenv("PORT", "abc")
		t.Setenv("DB_PATH", "/tmp/test.db")

		if _, err := LoadConfigFromEnv(); err == nil {
			t.Fatal("expected error for non-numeric PORT")
		}
	})

	t.Run("missing db_path", func(t *testing.T) {
		t.Setenv("HOST", "127.0.0.1")
		t.Setenv("PORT", "8080")
		t.Setenv("DB_PATH", "")

		if _, err := LoadConfigFromEnv(); err == nil {
			t.Fatal("expected error for missing DB_PATH")
		}
	})
}
