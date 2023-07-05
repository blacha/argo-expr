package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func execCommand(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	out := new(bytes.Buffer)

	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, out.String(), err
}

func TestRunHello(t *testing.T) {
	t.Run("should output clean result", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "{{=input}}", "--value", "input=hello world")
		assert.Equal(t, "hello world\n", output)
	})

	t.Run("should test basic math", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "{{=asInt(input) + 5}}", "--value", "input=5")
		assert.Equal(t, "10\n", output)
	})

	t.Run("should work with sprig basic math", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "{{= sprig.empty(input) ?  'empty' : input }}", "--value", "input=hello")
		assert.Equal(t, "hello\n", output)
	})

	t.Run("should not need --value", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "{{= sprig.empty(input) ?  'empty' : input }}")
		assert.Equal(t, "empty\n", output)
	})

	t.Run("should output as JSON", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "{{= input }}", "--json", "--value", "input=hello")
		assert.Equal(t, `{"result":"hello","template":"{{= input }}","values":{"input":"hello"}}`+"\n", output)
	})

	t.Run("should hide warnings with --quiet as JSON", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "{{= input }}", "--json", "--value", "input=hello")
		assert.Equal(t, `{"result":"hello","template":"{{= input }}","values":{"input":"hello"}}`+"\n", output)
	})

	t.Run("should read from file", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "--from-file", "test/input.json")
		assert.Equal(t, "hello\n", output)
	})

	t.Run("should read from file allow value override", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "--from-file", "test/input.json", "--value", "inputs.parameters.message=🦄🌈")
		assert.Equal(t, "🦄🌈\n", output)
	})

	t.Run("should read from file allow input override", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "--from-file", "test/input.json", "{{= inputs.parameters.message + ' world'}}", "--quiet")
		assert.Equal(t, "hello world\n", output)
	})

	t.Run("should log a warning when overriding input", func(t *testing.T) {
		_, output, _ := execCommand(create_command(), "--from-file", "test/input.json", "{{= inputs.parameters.message + ' world'}}")
		assert.Equal(t, "Replacing template from:'{{= inputs.parameters.message }}' to:'{{= inputs.parameters.message + ' world'}}'\nhello world\n", output)
	})
}
