Stories
=======

TO DO
-----

AS A: developer
I WANT: "runt fooBar" to perform no action when directory fooBar contains no executables
SO THAT: we maintain the principle of least surprise.

AS A: developer
I WANT: my command-line tool to block until the batch is complete
SO THAT: I do not have to worry about polling for completion.

AS A: developer
I WANT: a "batch" to refer to all test executables (or links thereto) recursively found in a single named directory
SO THAT: I can launch many tests conveniently.

AS A: developer
I WANT: a test to be invoked as a shell executable
SO THAT: I do not have to depend on Go APIs.

AS A: developer
I WANT: an endpoint that I can hit with my web browser to see how far along a particular batch is
SO THAT: I can keep track of test progress.

AS A: developer
I WANT: the opportunity to see the stdout and stderr streams of any test, successful or not
SO THAT: I can diagnose problems without resorting to applicaion-fragile tooling.

AS A: test runner
I WANT: a test command to return 0 upon success, non-zero otherwise
SO THAT: I can keep track of successful vs. unsuccessful test runs.

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

