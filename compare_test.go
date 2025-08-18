package compare_test

import (
	"bytes"
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
	t.Run("Wrong", func(t *testing.T) {
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
	})
	t.Run("no file", func(t *testing.T) {
		var c checker
		compare.Test(&c, "no file or wrong", []byte(""))
		compare.Test(t, ".NoFile", []byte(c.err.Error()))
	})
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

func Example() {
	{
		old := compare.App
		compare.App = "universal"
		defer func() {
			compare.App = old
		}()
	}

	var c checker

	wrong := []byte("bad\n")
	compare.Test(&c, ".test", wrong)
	if !c.iserror {
		panic("check wrong body")
	}
	fmt.Fprintf(os.Stdout, "%v", c.err)

	// Output:
	// *   1 "good"                                  "bad"
	// *   2 "good"                                  ""
	// *   3 "good"                                  << EMPTY LINE>>
	// *   4 "good"                                  << EMPTY LINE>>
	// and more other ...
	// universal ".test" ".test.new" &
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

func TestDiff(t *testing.T) {
	tcs := [][2]string{{
		"Output:1\nsame\nSecond",
		"Output:2\nsame\nFirst",
	}, {
		"Channel1",
		"Channel2\nOne\nTwo\n",
	}, {
		"o1\no2\no3\n04\n05\n06\n08\n09\n10",
		"o1\no2\no3\n04\n05\n06\n08\n19\n10",
	}}
	{
		var buf bytes.Buffer
		for i := 0; i < 1500; i++ {
			fmt.Fprintf(&buf, "%06d\n", i)
		}
		base := buf.String()
		diff := buf.Bytes()
		diff[301] = '1'
		tcs = append(tcs, [2]string{base, string(diff)})
	}
	for index, tc := range tcs {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			var c checker
			compare.TestDiff(
				&c,
				[]byte(tc[0]),
				[]byte(tc[1]),
			)
			if !c.iserror {
				return
			}
			compare.Test(t, fmt.Sprintf(".ShowDiff%d", index), []byte(c.err.Error()))
		})
	}
}
