package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-vfs/v2"

	"github.com/twpayne/chezmoi/internal/chezmoitest"
)

func TestManagedCmd(t *testing.T) {
	for _, tc := range []struct {
		name           string
		root           interface{}
		args           []string
		expectedOutput string
	}{
		{
			name: "simple",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi/dot_file": "# contents of .file\n",
			},
			expectedOutput: chezmoitest.JoinLines(
				".file",
			),
		},
		{
			name: "template",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi/dot_template.tmpl": "{{ fail \"Template should not be executed\" }}\n",
			},
			expectedOutput: chezmoitest.JoinLines(
				".template",
			),
		},
		{
			name: "create_template",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi/create_dot_file.tmpl": "{{ fail \"Template should not be executed\" }}\n",
			},
			expectedOutput: chezmoitest.JoinLines(
				".file",
			),
		},
		{
			name: "modify_template",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi/modify_dot_file.tmpl": "{{ fail \"Template should not be executed\" }}\n",
			},
			expectedOutput: chezmoitest.JoinLines(
				".file",
			),
		},
		{
			name: "script_template",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi/run_script.tmpl": "{{ fail \"Template should not be executed\" }}\n",
			},
			args: []string{
				"--include", "scripts",
			},
			expectedOutput: chezmoitest.JoinLines(
				"script",
			),
		},
		{
			name: "symlink_template",
			root: map[string]interface{}{
				"/home/user/.local/share/chezmoi/symlink_dot_symlink.tmpl": "{{ fail \"Template should not be executed\" }}\n",
			},
			expectedOutput: chezmoitest.JoinLines(
				".symlink",
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			chezmoitest.WithTestFS(t, tc.root, func(fs vfs.FS) {
				stdout := &bytes.Buffer{}
				require.NoError(t, newTestConfig(t, fs, withStdout(stdout)).execute(append([]string{"managed"}, tc.args...)))
				assert.Equal(t, tc.expectedOutput, stdout.String())
			})
		})
	}
}
