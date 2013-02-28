Runt
====

Runt executes test suites in parallel processes, collects stdout/stderr from each, and produces output compatible with Logstash and Kibana suitable for convenient viewing.

Install
-------

    go get github.com/sam-falvo/runt
    go install github.com/sam-falvo/runt
    export PATH=$PATH:$PWD/bin/runt

Usage
-----

Basic Usage
~~~~~~~~~~~

Assume you have a directory that contains a number of shell-executable files (e.g., binary programs or shell scripts) called 'tests'.  To execute these tests, simply enter:

    runt tests

The tests will execute, their outputs and result codes collected, and when all the executables have completed, a JSON event dump will appear on stdout.

Intended Usage
~~~~~~~~~~~~~~

Typically, you'll have a "log sender" program, such as Beaver, operating which monitors a log file, much as "tail -f" would.  Assuming this log file is called test-results.log, you would invoke runt as follows:

    runt tests >>test-results.log

After the json_event records appear in the log file, the log sender will forward them to the log aggregator of your choice.  I use json_event structures because they're the native format for Logstash, which is the aggregator I use.

