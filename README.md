# Runt

Runt executes test suites in parallel processes, collects stdout/stderr from each, and produces output compatible with [Logstash](http://logstash.net/) and [Kibana](http://kibana.org/) suitable for convenient viewing.

Although I wrote Runt using Go, it is not in any way restricted to invoking Go test suites.  In fact, my original use-case for this tool is to invoke Whiskey-based Node.js integration tests.

## Install

    go get github.com/sam-falvo/runt
    go install github.com/sam-falvo/runt
    export PATH=$PATH:$PWD/bin/runt

## Usage
### Basic Usage

Let's pretend you're writing a large, Enterprise-like application that requires long-running integration or systems tests.  Ideally, these tests should run in parallel, so as to consume as little time as possible.  Let's put these test suites into a directory named "suites":

    mkdir suites
    cat <<EOF >suites/suite-1.sh
    #!/bin/bash
    echo "Hello from suite 1"
    sleep 15
    EOF
    cat <<EOF >suites/suite-2.sh
    #!/bin/bash
    echo "Suite 2 says it's a fast suite."
    sleep 2
    EOF
    cat <<EOF >suites/suite-3.sh
    #!/bin/bash
    echo "Suite 3 is suicidal.  You can't trust it."
    sleep 4
    kill -11 \$\$
    EOF
    cat <<EOF >suites/suite-4.sh
    #!/bin/bash
    echo "uh oh -- I generate a warning, but otherwise finish to completion" >&2
    sleep 6
    EOF
    chmod a+x suites/*.sh

To execute these test-suites, simply enter:

    runt suites

The test suites will execute, their outputs and result codes collected, and when all the executables have completed, a JSON event dump will appear on stdout.  By default, Runt will restrict itself to running no more than four processes at once, to prevent accidentally fork-bombing itself.

You can confirm for yourself that runt is, in fact, executing these tests in parallel, either by looking at a process listing, or via the time command:

    time runt suites

You should see the execution time not much longer than 15 seconds, which is the slowest suite we've created above.

The resulting output is in Logstash's JSON-event format, and should look something like this:

    {"@timestamp":"2013-02-27T22:48:00.058099-08:00","@tags":[],"@type":"ShellCommand","@source":"Runt Demo","@fields":{"Executable":"suites/suite-2.sh","Stdout":"Suite 2 says it's a fast suite.\n","Stderr":""},"@message":"Command completed successfully."}
    {"@timestamp":"2013-02-27T22:48:00.058124-08:00","@tags":[],"@type":"ShellCommand","@source":"Runt Demo","@fields":{"Executable":"suites/suite-3.sh","Stdout":"Suite 3 is suicidal.  You can't trust it.\n","Stderr":""},"@message":"Error: signal 11"}
    {"@timestamp":"2013-02-27T22:48:00.058127-08:00","@tags":[],"@type":"ShellCommand","@source":"Runt Demo","@fields":{"Executable":"suites/suite-4.sh","Stdout":"","Stderr":"uh oh -- I generate a warning, but otherwise finish to completion\n"},"@message":"Command completed successfully."}
    {"@timestamp":"2013-02-27T22:48:00.05813-08:00","@tags":[],"@type":"ShellCommand","@source":"Runt Demo","@fields":{"Executable":"suites/suite-1.sh","Stdout":"Hello from suite 1\n","Stderr":""},"@message":"Command completed successfully."}

Observe how stdout, stderr, and the shell command's final status appear in the log entries.

### Intended Usage

Typically, you'll have a "log sender" program, such as Beaver, operating which monitors a log file, much as "tail -f" would.  Assuming this log file is called test-results.log, you would invoke runt as follows:

    runt tests >>test-results.log

After the json_event records appear in the log file, the log sender will forward them to the log aggregator of your choice.  I use json_event structures because they're the native format for Logstash, which is the aggregator I use.

