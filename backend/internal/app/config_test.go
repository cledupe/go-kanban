package app

import "testing"

func TestLoadConfigFromEnvRejectsMissingOrInvalidSettings(t *testing.T) {
	t.Run("missing host", func(t *testing.T) {
		t.Setenv("HOST", "")
		t.Setenv("PORT", "8080")

		if _, err := LoadConfigFromEnv(); err == nil {
			t.Fatal("expected error for missing HOST")
		}
	})

	t.Run("missing port", func(t *testing.T) {
		t.Setenv("HOST", "127.0.0.1")
		t.Setenv("PORT", "")

		if _, err := LoadConfigFromEnv(); err == nil {
			t.Fatal("expected error for missing PORT")
		}
	})

	t.Run("invalid port", func(t *testing.T) {
		t.Setenv("HOST", "127.0.0.1")
		t.Setenv("PORT", "abc")

		if _, err := LoadConfigFromEnv(); err == nil {
			t.Fatal("expected error for non-numeric PORT")
		}
	})
}
