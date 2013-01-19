// vim:ts=8:sw=8:noexpandtabs

package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

func get_stream(stream io.ReadCloser, finished chan<- [][]byte) {
	list := make([][]byte, 0)
	buf := make([]byte, 256)
	for n, err := stream.Read(buf); (err == nil) && (n > 0); {
		list = append(list, buf)
		buf = make([]byte, 256)
		n, err = stream.Read(buf)
	}

	finished <- list;
}

func main() {
	cmd := exec.Command("asdf/t1.sh")

	stdout,err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr,err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	solc := make(chan [][]byte)
	selc := make(chan [][]byte)

	go get_stream(stdout, solc)
	go get_stream(stderr, selc)

	var sol [][]byte
	var sel [][]byte

	fmt.Println("wait...")
	sol = <- solc
	sel = <- selc
	fmt.Println("done...")

	fmt.Println("STDOUT:")
	for _, v := range sol {
		fmt.Printf("  %s\n", v)
	}

	fmt.Println("STDERR:")
	for _, v := range sel {
		fmt.Printf("# %s\n", v)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

