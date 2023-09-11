package compare

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	diff "github.com/olegfedoseev/image-diff"
)

// Basic values
const (
	Key      = "UPDATE"
	KeyValid = "true"
)

// Testing interface is valid for *testing.T
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

// Save `.png` files
func Save(filename string, img image.Image) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	if err = png.Encode(f, img); err != nil {
		return
	}
	if err = f.Close(); err != nil {
		return
	}
	return
}

// TestPng compare `.png` files
// for update test screens run in console:
// UPDATE=true go test
func TestPng(t Testing, filename string, actual image.Image) {
	if os.Getenv(Key) == KeyValid {
		err := Save(filename, actual)
		if err != nil {
			t.Errorf("Cannot save `%s`: %v", filename, err)
			return
		}
	}
	if err := func() error {
		if actual == nil {
			return fmt.Errorf("image is nil")
		}
		// get expect result
		dir, err := os.MkdirTemp("", "actual")
		if err != nil {
			return fmt.Errorf("Cannot create temp folder: %v", err)
		}
		actualFilename := filepath.Join(dir, "a.png")
		err = Save(actualFilename, actual)
		if err != nil {
			return fmt.Errorf("Cannot save `%s`: %v", actualFilename, err)
		}
		diff, percent, err := diff.CompareFiles(filename, actualFilename)
		if err != nil {
			return fmt.Errorf("Cannot compare images: %v", err)
		}
		if percent == 0.0 {
			return nil
		}
		err = fmt.Errorf("Images is different at %.2f percent", percent)
		// save diff image
		errS := Save(filename+".new.png", diff)
		if errS != nil {
			err = fmt.Errorf("Cannot save `%s`: %v. %v", actualFilename, err, errS)
		}
		return err
	}(); err != nil {
		t.Errorf("%s: %v", filename, err)
	}
}
