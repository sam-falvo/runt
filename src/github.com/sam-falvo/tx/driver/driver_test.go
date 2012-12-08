// vim:ts=8:sw=8:noexpandtab:

package driver

import "os"
import "testing"
import "time"

type fileExistsButNotADirectory struct {};
type validDirectory struct { fileExistsButNotADirectory };

func (fs *fileExistsButNotADirectory) Name() string { return "don't care" }
func (fs *fileExistsButNotADirectory) Size() int64 { return 0 }
func (fs *fileExistsButNotADirectory) Mode() os.FileMode { return 0644 }
func (fs *fileExistsButNotADirectory) ModTime() time.Time { return time.Now() }
func (fs *fileExistsButNotADirectory) IsDir() bool { return false }
func (fs *fileExistsButNotADirectory) Sys() interface{} { return nil }

func (fs *validDirectory) Name() string { return "don't care" }
func (fs *validDirectory) Mode() os.FileMode { return 0644 | os.ModeDir }
func (fs *validDirectory) IsDir() bool { return true }

// AS A: developer
// I WANT: to type in a shell command "run fooBar" to run all integration tests collectively part of the fooBar suite
// SO THAT: I can integrate my integration tests with the build environment of my choice.

//	AS A: implementor
//	I WANT: TX to produce a meaningful error if batch parameter proves anything except a directory or a link thereto
//	SO THAT: we can avoid special-case logic for handling individual test files.
//	NOTE: Remember that tests are shell executables.  "run fooBar.exe" is the same as just typing "fooBar.exe".

func TestDriverShouldVerifyBatchIsADirectory(t *testing.T) {
	d := new(Driver)
	d.UseStat(func (_ string) (fi os.FileInfo, e error) {
		return &fileExistsButNotADirectory{}, nil
	})

	err := d.UseBatch("alksjdhflakjsdf")
	if err != DirectoryExpectedError {
		t.Errorf("Expected indicated path must be a directory.")
	}
}

func TestDriverShouldOnlyUseDirectoriesThatExist(t *testing.T) {
	d := new(Driver)
	d.UseStat(func (_ string) (fi os.FileInfo, e error) {
		return nil, os.ErrNotExist
	})

	err := d.UseBatch("MissingDirectory")
	if err != os.ErrNotExist {
		t.Errorf("Expected indicated path must exist.")
	}
}

func TestDriverShouldNotYieldErrorIfGivenDirectory(t *testing.T) {
	d := new(Driver)
	d.UseStat(func (_ string) (fi os.FileInfo, e error) {
		return &validDirectory{}, nil
	})

	err := d.UseBatch("akjdhf")
	if err != nil {
		t.Errorf("Unexpected error when given a directory: %s", err.Error())
	}
}

