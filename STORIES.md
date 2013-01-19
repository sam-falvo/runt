Stories
=======

TO DO
-----

AS A: test runner
I WANT: a test command to return 0 upon success, non-zero otherwise
SO THAT: I can keep track of successful vs. unsuccessful test runs.

QUEUED
------

AS A: developer
I WANT: the opportunity to see the stdout and stderr streams of any test, successful or not
SO THAT: I can diagnose problems without resorting to applicaion-fragile tooling.


DOING
-----

DONE
----

AS A: developer
I WANT: to type in a shell command "runt fooBar" to run all integration tests collectively part of the fooBar suite
SO THAT: I can integrate my integration tests with the build environment of my choice.

	AS A: implementor
	I WANT: TX to produce a meaningful error if batch parameter proves anything except a directory or a link thereto
	SO THAT: we can avoid special-case logic for handling individual test files.

	AS A: implementor
	I WANT: runt to iterate through all the directory entries in fooBar, recursively
	SO THAT: we can decide which are executable files and which are not.

	AS A: implementor
	I WANT: TX to run each candidate exactly once.
	SO THAT: we avoid duplicate invokations of any given test.

AS A: developer
I WANT: "runt fooBar" to perform no action when directory fooBar contains no executables
SO THAT: we maintain the principle of least surprise.

AS A: developer
I WANT: a "batch" to refer to all test executables (or links thereto) recursively found in a single named directory
SO THAT: I can launch many tests conveniently.

AS A: developer
I WANT: a test to be invoked as a shell executable
SO THAT: I do not have to depend on Go APIs.

WBS:
	Create a type, say ChildResult, that encapsulates a forked process and its stdout and stderr streams.
	- executable string
	- subprocessError error
	- stdout [][]byte
	- stderr [][]byte
	Create a function that takes an executable as input, and synchronously returns a ChildResult.
	- When dispatching an executable, send a value to the semaphore.
	  - If the channel is full, the sending process will block!
	  - Wait for the child process to complete.
	  - Read from the semaphore channel.  This will unblock any waiting goroutines.
	  - Send the child process results back to the runt driver.
	Wrap each child process in a goroutine.
	  - Use os.exec.Command to launch the command by name.  Pass no arguments.
	  - Use os.StdoutPipe() to create a custom io.ReadCloser for the child's stdout stream.
	  - Use os.StderrPipe() to create a custom io.ReadCloser for the child's stderr stream.
	  - Use cmd.Start() to fork and exec the command process.
	  - Use a goroutine to bulk-read data in from stdout.
	  - Use a goroutine to bulk-read data in from stderr.
	  - Use cmd.Wait() to wait for the command to finish and return its status to the parent.
	Create a driver method to launch the subprocesses and block until all have finished, returning an array of ChildResult instances.
	- Create a semaphore channel with 4 elements.
	- For each executable discovered, launch a goroutine to dispatch the executable.  Make sure it uses the semaphore channel.

AS A: developer
I WANT: my command-line tool to block until the batch is complete
SO THAT: I do not have to worry about polling for completion.


TRASH
-----

AS A: developer
I WANT: an endpoint that I can hit with my web browser to see how far along a particular batch is
SO THAT: I can keep track of test progress.


