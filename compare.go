package compare

import (
	"bytes"
	"io/ioutil"
	"os"
)

const (
	Key      = "UPDATE"
	KeyValid = "true"
)

type Testing interface {
	Errorf(format string, args ...any)
}

// for update test screens run in console:
// UPDATE=true go test
func Test(t Testing, filename string, actual []byte) {
	if os.Getenv(Key) == KeyValid {
		if err := ioutil.WriteFile(filename, actual, 0644); err != nil {
			t.Errorf("Cannot write snapshot to file: %v", err)
			return
		}
	}
	// get expect result
	expect, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Cannot read snapshot file: %v", err)
		return
	}
	// compare
	if !bytes.Equal(actual, expect) {
		f2 := filename + ".new"
		if err := ioutil.WriteFile(f2, actual, 0644); err != nil {
			t.Errorf("Cannot write snapshot to file new: %v", err)
			return
		}
		size := 1000
		if size < len(actual) {
			actual = actual[:size]
		}
		if size < len(expect) {
			expect = expect[:size]
		}
		t.Errorf("Snapshots is not same:\nActual:\n%s\nExpect:\n%s\nmeld %s %s",
			actual,
			expect,
			filename, f2,
		)
	}
}
