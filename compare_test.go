package compare_test

import (
	"os"
	"testing"

	"github.com/Konstantin8105/compare"
)

type checker struct {
	iserror bool
}

func (c *checker) Errorf(format string, args ...any) {
	c.iserror = true
}

func TestWrong(t *testing.T) {
	var c checker
	compare.Test(&c, "/EEEe/d", nil)
	if !c.iserror {
		t.Fatal("no wrong_name")
	}
}

func TestStore(t *testing.T) {
	var c checker

	var actual []byte
	for i := 0; i < 2300; i++ {
		actual = append(actual, []byte("good\n")...)
	}

	os.Setenv(compare.Key, compare.KeyValid)
	compare.Test(&c, ".test", actual)
	if c.iserror {
		t.Fatal("add key-value")
	}

	os.Unsetenv(compare.Key)
	compare.Test(&c, ".test", actual)
	if c.iserror {
		t.Fatal("check error")
	}

	wrong := []byte("bad\n")
	compare.Test(&c, ".test", wrong)
	if !c.iserror {
		t.Fatal("check wrong body")
	}
}
