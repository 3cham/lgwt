package context

import "testing"

// This calls the Program and we could cancelled it with ctrl-C and still invoke the cleaning event
func TestProgram(t *testing.T) {
	t.Skip("skipping because interrupt signal is not sent by intellij")
	Program()
}
