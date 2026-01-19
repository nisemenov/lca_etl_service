package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name          string
		env           map[string]string
		wantDBPath    string
		wantPanic     bool
		wantPanicText string
	}{
		{
			name:       "DB_PATH задан через переменную окружения",
			env:        map[string]string{"DB_PATH": "/data/myapp.db"},
			wantDBPath: "/data/myapp.db",
			wantPanic:  false,
		},
		{
			name:          "DB_PATH не задан → должен паниковать",
			env:           map[string]string{},
			wantDBPath:    "",
			wantPanic:     true,
			wantPanicText: "DB_PATH is required",
		},
		{
			name:          "DB_PATH пустая строка → тоже паника",
			env:           map[string]string{"DB_PATH": ""},
			wantDBPath:    "",
			wantPanic:     true,
			wantPanicText: "DB_PATH is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				require.NoError(t, os.Setenv(k, v))
			}

			t.Cleanup(func() {
				os.Clearenv()
			})

			var cfg *Config
			var panicMsg interface{}

			func() {
				defer func() {
					panicMsg = recover()
				}()
				cfg = Load()
			}()

			if tt.wantPanic {
				require.NotNil(t, panicMsg, "ожидалась паника, но её не было")
				assert.Equal(t, tt.wantPanicText, panicMsg,
					"неверное сообщение в панике")
				return
			}

			require.Nil(t, panicMsg, "неожиданная паника")
			require.NotNil(t, cfg, "конфиг не должен быть nil")
			assert.Equal(t, tt.wantDBPath, cfg.DBPath)
		})
	}
}

func Test_getEnv(t *testing.T) {
	t.Run("переменная существует", func(t *testing.T) {
		os.Setenv("TEST_KEY", "custom-value")
		defer os.Unsetenv("TEST_KEY")

		got := getEnv("TEST_KEY", "default")
		assert.Equal(t, "custom-value", got)
	})

	t.Run("переменная не существует → fallback", func(t *testing.T) {
		os.Unsetenv("TEST_KEY")

		got := getEnv("TEST_KEY", "fallback-value")
		assert.Equal(t, "fallback-value", got)
	})

	t.Run("переменная существует, но пустая", func(t *testing.T) {
		os.Setenv("TEST_KEY", "")
		defer os.Unsetenv("TEST_KEY")

		got := getEnv("TEST_KEY", "default")
		assert.Equal(t, "", got, "должна возвращаться именно пустая строка")
	})
}
