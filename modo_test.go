package jsonschema

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadingGeneratedSchemas(t *testing.T) {
	t.Parallel()

	filesystem := os.DirFS("testdata")
	_ = fs.WalkDir(filesystem, "modo", func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err)
		if d.IsDir() {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			assert.NotEqual(t, "", d.Name())
			f, err := filesystem.Open(path)
			require.NoError(t, err)
			data, err := io.ReadAll(f)
			require.NoError(t, err)

			var s Schema
			err = json.Unmarshal(data, &s)
			require.NoError(t, err)

			_, err = s.ValidateBytes(context.Background(), []byte(`{}`))
			require.NoError(t, err)
		})
		return nil
	})
}
