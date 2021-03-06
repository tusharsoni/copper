package cconfigtest

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/cconfig"
)

// SetupDirWithConfigs creates a temp directory that can store config files.
// It creates base.toml, test.toml, and secrets.toml with the given strings.
// The directory is cleaned up after test run.
func SetupDirWithConfigs(t *testing.T, configs ...string) cconfig.Dir {
	t.Helper()

	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	if len(configs) >= 1 {
		err = ioutil.WriteFile(path.Join(dir, "base.toml"), []byte(configs[0]), os.ModePerm)
		assert.NoError(t, err)
	}

	if len(configs) >= 2 { // nolint:gomnd
		err = ioutil.WriteFile(path.Join(dir, "test.toml"), []byte(configs[1]), os.ModePerm)
		assert.NoError(t, err)
	}

	if len(configs) >= 3 { // nolint:gomnd
		err = ioutil.WriteFile(path.Join(dir, "secrets.toml"), []byte(configs[2]), os.ModePerm)
		assert.NoError(t, err)
	}

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll(dir))
	})

	return cconfig.Dir(dir)
}
