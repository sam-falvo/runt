// vim:ts=8:sw=8:noexpandtab:

package main

import (
	"flag"
	"fmt"
	"github.com/sam-falvo/tx/driver"
	"os"
)

func problem(e error) {
	fmt.Fprintf(os.Stderr, "%s\n", e.Error())
	fmt.Fprintf(os.Stderr, "USAGE: %s batch-dir\n", os.Args[0]);
	os.Exit(1)
	panic("os.Exit() didn't os.Exit()!")
}

func main() {
	var err error

	d := new(driver.Driver)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		problem(fmt.Errorf("You must at least provide a test batch parameter."))
	}
	err = d.UseBatch(args[0])
	if err != nil {
		problem(err)
	}
	err = d.LaunchSuites()
	if err != nil {
		problem(err)
	}
	j, err := d.JsonEvents()
	if err != nil {
		problem(err)
	}
	for _, v := range j {
		fmt.Println(v)
	}
}

