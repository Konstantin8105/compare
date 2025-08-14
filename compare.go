package compare

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	diff "github.com/olegfedoseev/image-diff"
)

var App = "meld" // linux
func init() {
	if runtime.GOOS == "windows" {
		App = "WinMergeU.exe"
	}
}

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
	// remove ends
	actual = bytes.ReplaceAll(actual, []byte("\r"), nil)
	// comparing
	if os.Getenv(Key) == KeyValid {
		if err := os.WriteFile(filename, actual, 0644); err != nil {
			t.Errorf("Cannot write snapshot to file: %w", err)
			return
		}
	}
	// get expect result
	expect, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Cannot read snapshot file: %w", err)
		return
	}
	// remove ends
	expect = bytes.ReplaceAll(expect, []byte("\r"), nil)
	// compare
	if !bytes.Equal(actual, expect) {
		f2 := filename + ".new"
		if err := os.WriteFile(f2, actual, 0644); err != nil {
			t.Errorf("Cannot write snapshot to file new: %w", err)
			return
		}
		t.Errorf("%v\n%s \"%s\" \"%s\" &",
			Diff(expect, actual),
			App,
			filename, f2)
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
			return fmt.Errorf("cannot create temp folder: %w", err)
		}
		actualFilename := filepath.Join(dir, "a.png")
		err = Save(actualFilename, actual)
		if err != nil {
			return fmt.Errorf("cannot save `%s`: %w", actualFilename, err)
		}
		diff, percent, err := diff.CompareFiles(filename, actualFilename)
		if err != nil {
			return fmt.Errorf("cannot compare images: %w", err)
		}
		if percent == 0.0 {
			return nil
		}
		err = fmt.Errorf("images is different at %.2f percent", percent)
		// save diff image
		errS := Save(filename+".new.png", diff)
		if errS != nil {
			err = fmt.Errorf("cannot save `%s`: %w. %w", actualFilename, err, errS)
		}
		return err
	}(); err != nil {
		t.Errorf("%s: %v", filename, err)
	}
}

// TestDiff will print two strings vertically next to each other so that line
// differences are easier to read.
func TestDiff(t Testing, txt1, txt2 []byte) {
	if err := Diff(txt1, txt2); err != nil {
		t.Errorf("%v", err)
	}
}

// Diff will print two strings vertically next to each other so that line
// differences are easier to read.
func Diff(txt1, txt2 []byte) (err error) {
	a := string(txt1)
	b := string(txt2)
	if a == b {
		return
	}
	//
	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")
	maxLines := int(math.Max(float64(len(aLines)), float64(len(bLines))))
	out := "\n"
	view := false
	var viewAmount int

	for lineNumber := 0; lineNumber < maxLines; lineNumber++ {
		aLine := "<< EMPTY LINE>>"
		bLine := "<< EMPTY LINE>>"

		// Replace NULL characters with a dot. Otherwise the strings will look
		// exactly the same but have different length (and therfore not be
		// equal).
		if lineNumber < len(aLines) {
			aLine = strconv.Quote(aLines[lineNumber])
		}
		if lineNumber < len(bLines) {
			bLine = strconv.Quote(bLines[lineNumber])
		}

		diffFlag := " "
		if aLine != bLine {
			view = true
			diffFlag = "*"
		}
		if !view {
			continue
		}
		viewAmount++
		out += fmt.Sprintf("%s %3d %-40s%s\n", diffFlag, lineNumber+1, aLine, bLine)

		if len(aLines) < lineNumber || len(bLines) < lineNumber || 20 < viewAmount {
			out += "and more other ..."
			break
		}
	}
	err = fmt.Errorf("%s", out)
	return
}
