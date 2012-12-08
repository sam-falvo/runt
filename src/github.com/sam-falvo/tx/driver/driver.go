// vim:ts=8:sw=8:noexpandtab:

// driver package implements the coordinating logic for the runt command.
// It ultimately is responsible for instantiating a TAS and one or more TX
// instances as per the user's provided configuration.
package driver

import "fmt"
import "os"


// statFn is a prototype matching os.Stat().  Not used for normal code,
// statFn proves useful for configuring a customized stat()-like behavior for
// the purposes of testing edge cases easily.  It, in effect, isolates the driver
// from the Go standard library and host operating system.
type statFn func (string) (os.FileInfo, error)

// This type represents a single test runner instance.
type Driver struct {
	stat	statFn
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

	return nil
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

// stat switches between the driver-specific stat procedure or os.Stat.
func stat(d *Driver, p string) (os.FileInfo, error) {
	if d.stat != nil {
		return d.stat(p);
	}
	return os.Stat(p)
}

