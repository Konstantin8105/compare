package compare_test

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/Konstantin8105/compare"
)

type checker struct {
	iserror bool
	err     error
}

func (c *checker) Errorf(format string, args ...any) {
	c.iserror = true
	c.err = fmt.Errorf(format, args...)
}

func TestWrong(t *testing.T) {
	var c checker

	os.Setenv(compare.Key, compare.KeyValid)

	for i := 0; i < 2; i++ {
		compare.Test(&c, "/EEEe/d", nil)
		if !c.iserror {
			t.Fatal("no wrong_name")
		}
		compare.TestPng(&c, "/sdsdas3/3/", nil)
		if !c.iserror {
			t.Fatal("no wrong_name")
		}
		compare.TestPng(&c, "/sdsdas3/3/", defaultPng())
		if !c.iserror {
			t.Fatal("no wrong_name")
		}

		os.Unsetenv(compare.Key)
	}
}

func Test(t *testing.T) {
	var c checker

	var actual []byte
	for i := 0; i < 2300; i++ {
		actual = append(actual, []byte("good\n")...)
	}

	os.Setenv(compare.Key, compare.KeyValid)
	compare.Test(&c, ".test", actual)
	if c.iserror {
		t.Fatalf("add key-value: %v", c.err)
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

func defaultPng() image.Image {
	size := 100
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(x*250/100) % 250,
				G: uint8(y*250/100) % 250,
				B: uint8((x+y)/100) % 250,
				A: 225,
			})
		}
	}
	return img
}

func TestPng(t *testing.T) {
	var c checker

	actual := defaultPng()

	os.Setenv(compare.Key, compare.KeyValid)
	compare.TestPng(&c, ".test.png", actual)
	if c.iserror {
		t.Fatalf("add key-value: %v", c.err)
	}

	os.Unsetenv(compare.Key)
	compare.TestPng(&c, ".test.png", actual)
	if c.iserror {
		t.Fatal("check error")
	}

	wrong := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	compare.TestPng(&c, ".test.png", wrong)
	if !c.iserror {
		t.Fatalf("check wrong body: %v", c.err)
	}
}
