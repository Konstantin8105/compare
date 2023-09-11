# compare
Compare testing

```
package compare // import "github.com/Konstantin8105/compare"


CONSTANTS

const (
	Key      = "UPDATE"
	KeyValid = "true"
)
    Basic values


FUNCTIONS

func Save(filename string, img image.Image) (err error)
    Save `.png` files

func Test(t Testing, filename string, actual []byte)
    for update test screens run in console: UPDATE=true go test

func TestPng(t Testing, filename string, actual image.Image)
    TestPng compare `.png` files for update test screens run in console:
    UPDATE=true go test


TYPES

type Testing interface {
	Errorf(format string, args ...any)
}
    Testing interface is valid for *testing.T
```
