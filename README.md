PACKAGE
=======

package main

    import "github.com/sam-falvo/tx/tas"

Package main implements a proof-of-concept for a test automation server
(TAS).

The idea is simple: a central server maintains a to-do list for test
executors (TXs). TXs and user-facing clients communicate with the TAS
via a REST interface. This TAS provides no authentication; a TX or
user-client may invoke any REST end-point it chooses.

This TAS is a spike, both to learn Go and to learn more about the
relationship between TAS, TX, and user client.

Eventually, this code will be thrown out and a real TAS/TX/client
architecture documented, constructed, and tested.

