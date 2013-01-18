// vim:ts=8:sw=8:noexpandtab:

package driver

import "os"
import "testing"
import "time"

type myFileInfo struct {
	name string
	mode os.FileMode
	isDir bool
};

func (f *myFileInfo) Name() string { return f.name }
func (f *myFileInfo) Mode() os.FileMode { return f.mode }
func (f *myFileInfo) IsDir() bool { return f.isDir }
func (f *myFileInfo) Sys() interface{} { return nil }
func (f *myFileInfo) ModTime() time.Time { return time.Now() }
func (f *myFileInfo) Size() int64 { return 0 }

var fileExistsButNotADirectory = myFileInfo {
	name: "don't care",
	mode: 0644,
	isDir: false,
}

var validDirectory = myFileInfo {
	name: "don't care",
	mode: 0644|os.ModeDir,
	isDir: true,
}

func aDir(fn string) *myFileInfo {
	return &myFileInfo {
		name: fn,
		mode: 0644|os.ModeDir,
		isDir: true,
	}
}

func anExec(fn string) *myFileInfo {
	return &myFileInfo {
		name: fn,
		mode: 0755,
		isDir: false,
	}
}

func aFile(fn string) *myFileInfo {
	return &myFileInfo {
		name: fn,
		mode: 0644,
		isDir: false,
	}
}

func emptyDirectory(_ string) ([]os.FileInfo, error) {
	return make([]os.FileInfo, 0), nil
}

func basicReadDir(p string) (fis []os.FileInfo, e error) {
	switch p {
	case "blah":
		fis = []os.FileInfo {
			aFile("a"),
			aFile("b"),
			aDir("c"),
			aDir("d"),
			anExec("e"),
			anExec("f"),
		}
	default:
		fis = []os.FileInfo {}
	}
	e = nil
	return
}

func deepReadDir(p string) (fis []os.FileInfo, e error) {
	switch p {
	case "blah":
		fis = []os.FileInfo {
			aFile("a"),
			aFile("b"),
			aDir("c"),
			aDir("d"),
			anExec("e"),
			anExec("f"),
		}
	case "blah/c":
		fis = []os.FileInfo {
			anExec("g"),
			aFile("gg"),
		}
	case "blah/d":
		fis = []os.FileInfo {
			anExec("h"),
			aFile("hh"),
		}
	default:
		fis = []os.FileInfo {}
	}
	e = nil
	return
}

func isElementOf(haystack []string, needle string) (found bool) {
	found = false
	for _, h := range haystack {
		if h == needle {
			found = true
		}
	}
	return
}

type testProc func (d *Driver)

func withSetup(rd readDirFn, fiStat os.FileInfo, eStat error, f testProc) {
	d := new(Driver)
	d.UseStat(func (_ string) (fi os.FileInfo, e error) {
		return fiStat, eStat
	})
	d.UseReadDir(rd)
	f(d)
}

// AS A: developer
// I WANT: to type in a shell command "run fooBar" to run all integration tests collectively part of the fooBar suite
// SO THAT: I can integrate my integration tests with the build environment of my choice.

//	AS A: implementor
//	I WANT: TX to produce a meaningful error if batch parameter proves anything except a directory or a link thereto
//	SO THAT: we can avoid special-case logic for handling individual test files.
//	NOTE: Remember that tests are shell executables.  "run fooBar.exe" is the same as just typing "fooBar.exe".

func TestDriverShouldVerifyBatchIsADirectory(t *testing.T) {
	withSetup(emptyDirectory, &fileExistsButNotADirectory, nil, func(d *Driver) {
		err := d.UseBatch("alksjdhflakjsdf")
		if err != DirectoryExpectedError {
			t.Errorf("Expected indicated path must be a directory.")
		}
	})
}

func TestDriverShouldOnlyUseDirectoriesThatExist(t *testing.T) {
	withSetup(emptyDirectory, nil, os.ErrNotExist, func(d *Driver) {
		err := d.UseBatch("MissingDirectory")
		if err != os.ErrNotExist {
			t.Errorf("Expected indicated path must exist.")
		}
	})
}

func TestDriverShouldNotYieldErrorIfGivenDirectory(t *testing.T) {
	withSetup(emptyDirectory, &validDirectory, nil, func(d *Driver) {
		err := d.UseBatch("akjdhf")
		if err != nil {
			t.Errorf("Unexpected error when given a directory: %s", err.Error())
		}
	})
}

//	AS A: implementor
//	I WANT: runt to iterate through all the directory entries in fooBar, recursively
//	SO THAT: we can decide which are executable files and which are not.

func TestDriverMustIsolateExecutablesInBatchDir(t *testing.T) {
	withSetup(basicReadDir, aDir("blah"), nil, func(d *Driver) {
		_ = d.UseBatch("blah")
		if len(d.Executables()) != 2 {
			t.Errorf("Gave a simple directory with 2 executables, but discovered %d", len(d.Executables()))
		}
	})
}

func TestDriverMustRecurseIntoSubdirectories(t *testing.T) {
	pathsRead := make([]string, 0, 4)
	withSetup(
		func (path string) (fis []os.FileInfo, e error) {
			pathsRead = append(pathsRead, path)
			return basicReadDir(path)
		},
		aDir("blah"),
		nil,
		func (d *Driver) {
			_ = d.UseBatch("blah")
			if len(pathsRead) != 3 {
				t.Errorf("Expected to traverse 3 directories; only %d traversed (%#v).", len(pathsRead), pathsRead)
			}
			if !isElementOf(pathsRead, "blah") {
				t.Errorf("Expected to traverse root of batch directory")
			}
			if !isElementOf(pathsRead, "blah/c") {
				t.Errorf("Expected to traverse blah/c")
			}
			if !isElementOf(pathsRead, "blah/d") {
				t.Errorf("Expected to traverse blah/d")
			}
		},
	)
}

func TestDriverMustQualifyExecutablePathNames(t *testing.T) {
	withSetup(deepReadDir, aDir("blah"), nil, func (d *Driver) {
		_ = d.UseBatch("blah")
		es := d.Executables()
		if !isElementOf(es, "blah/e") {
			t.Errorf("Expected to discover executable blah/e")
		}
		if !isElementOf(es, "blah/f") {
			t.Error("Expected to discover executable blah/f")
		}
		if !isElementOf(es, "blah/c/g") {
			t.Error("Expected to discover executable blah/c/g")
		}
		if !isElementOf(es, "blah/d/h") {
			t.Error("Expected to discover executable blah/d/h")
		}
	})
}


//	AS A: implementor
//	I WANT: TX to run each candidate exactly once.
//	SO THAT: we avoid duplicate invokations of any given test.


func TestDriverMustDispatchOneNameAtATime(t *testing.T) {
	expected := []string{ "blah/e", "blah/f", "blah/c/g", "blah/d/h" }

	withSetup(deepReadDir, aDir("blah"), nil, func (d *Driver) {
		_ = d.UseBatch("blah")
		for i := 0; i < len(expected); i++ {
			n, ok := d.NextExecutable()
			if !ok {
				t.Errorf("Expected dequeue to be OK")
			}
			if !isElementOf(expected, n) {
				t.Errorf("Expected %s to be member of %#v", n, expected)
			}
		}
		_, ok := d.NextExecutable()
		if ok {
			t.Errorf("We've exhausted the queue; this shouldn't be OK")
		}
	})
}

