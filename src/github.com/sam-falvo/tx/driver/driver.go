// vim:ts=8:sw=8:noexpandtab:

// driver package implements the coordinating logic for the runt command.
// It ultimately is responsible for instantiating a TAS and one or more TX
// instances as per the user's provided configuration.
package driver

import "fmt"
import "io/ioutil"
import "os"


// statFn is a prototype matching os.Stat().  Not used for normal code,
// statFn proves useful for configuring a customized stat()-like behavior for
// the purposes of testing edge cases easily.  It, in effect, isolates the driver
// from the Go standard library and host operating system.
type statFn func (string) (os.FileInfo, error)

// readDirFn is the prototye for the matching io.ioutil.ReadDir() function.
type readDirFn func (string) ([]os.FileInfo, error)

// This type represents a single test runner instance.
type Driver struct {
	stat	statFn
	readdir	readDirFn
	executables []string
}

// DirectoryExpectedError represents an error condition where a directory name was
// specified by a caller, but it actually refers to a non-directory object, such as
// an object or socket.
var DirectoryExpectedError error = fmt.Errorf("Directory expected")

// UseBatch specifies the batch of tests to work with.  The parameter names a
// directory in the local filesystem, within which zero or more test executables
// reside.
//
// Returns DirectoryExpectedError if the path provided exists, but isn't a
// directory.  Otherwise, any lower-level errors bubble up verbatim.
func (my *Driver) UseBatch(path string) error {
	fi, err := stat(my, path)
	if err != nil { return err }

	if !fi.IsDir() {
		return DirectoryExpectedError
	}

	my.executables = make([]string, 0)

	return discoverExecutables(my, path)
}

func discoverExecutables(d *Driver, dir string) error {
	fis, err := readdir(d, dir)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		qn := fmt.Sprintf("%s/%s", dir, fi.Name())
		if isExecutable(fi) {
			d.executables = append(d.executables, qn)
		}
		if fi.IsDir() {
			err = discoverExecutables(d, qn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isExecutable(fi os.FileInfo) bool {
	return (fi.Mode() & 0110) != 0
}

// Executables yields the driver's understanding of its batch of test executables
// to run.  Each filename appears relative to the batch directory; for example,
// if three executables A, B, and C appear inside the T batch directory, this
// function will return T/A, T/B, and T/C.
func (my *Driver) Executables() []string {
	return my.executables
}

// UseStat tells the driver which implementation of the stat system call to use.
// By default, it resorts to using os.Stat().  However, this call cannot be
// intercepted for unit-testing purposes; thus, unit tests rely on UseStat() to
// replace the Go library's implementation with their own.
//
// Any call to UseStat in production code is an error.
func (my *Driver) UseStat(s statFn) {
	my.stat = s
}

// UseReadDir tells the driver which implementation of the readdir function
// to use.  By default, it resorts to using io.ioutil.ReadDir().  However,
// this call cannot be intercepted for unit-testing purposes; thus, unit tests
// rely on UseReadDir() to replace the Go library's implementation with their
// own.
//
// Any call to UseReadDir in production code is an error.
func (my *Driver) UseReadDir(rd readDirFn) {
	my.readdir = rd
}

// stat switches between the driver-specific stat procedure or os.Stat.
func stat(d *Driver, p string) (os.FileInfo, error) {
	if d.stat != nil {
		return d.stat(p);
	}
	return os.Stat(p)
}

// readdir switches between the driver-specific ReadDir procedure or io.ioutil.ReadDir
func readdir(d *Driver, p string) ([]os.FileInfo, error) {
	if d.readdir != nil {
		return d.readdir(p)
	}
	return ioutil.ReadDir(p)
}

