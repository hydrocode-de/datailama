package cli_interface

import (
	"bytes"
	"testing"

	"github.com/hydrocode-de/datailama/internal/version"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestVersionOutput(t *testing.T) {
	// Test short version output
	t.Run("short version", func(t *testing.T) {
		var buf bytes.Buffer
		app := &cli.App{
			Commands: GetCommands(),
			Writer:   &buf,
		}

		err := app.Run([]string{"app", "v"})
		assert.NoError(t, err)
		assert.Equal(t, version.Version+"\n", buf.String())
	})

	// Test long version output
	t.Run("long version", func(t *testing.T) {
		var buf bytes.Buffer
		app := &cli.App{
			Commands: GetCommands(),
			Writer:   &buf,
		}

		err := app.Run([]string{"app", "version"})
		assert.NoError(t, err)

		expected := "DataILama CLI\n" +
			"Version: " + version.Version + "\n" +
			"Build Time: " + version.BuildTime + "\n" +
			"Git Commit: " + version.GitCommit + "\n"
		assert.Equal(t, expected, buf.String())
	})
}
