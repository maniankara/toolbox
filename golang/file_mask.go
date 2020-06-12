package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Wrong cli args !")
		fmt.Println("Usage: ./file_mask <string to mask> <file whose lines to be masked>")
		fmt.Println("E.g.: ./file_mask uanoop ./personal.txt")
		os.Exit(-1)
	}

	// open given file for read/write
	src, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println("Error: Unable to open file for writing: ", os.Args[2], err)
		os.Exit(-1)
	}

	// tmp file for backup
	err = ioutil.WriteFile(os.Args[2]+".tmp", src, 0755)
	if err != nil {
		fmt.Println("Error: Unable to open file for writing: ", os.Args[2]+".tmp", err)
		os.Exit(-1)
	}

	// replace contents of file
	lines := strings.Split(string(src), "\n")

	for i, line := range lines {
		if strings.Contains(line, os.Args[1]) {
			lines[i] = "X_REDACTED_X"
		}
	}

	replaced := strings.Join(lines, "\n")
	err = ioutil.WriteFile(os.Args[2], []byte(replaced), 0755)
	if err != nil {
		fmt.Println("Error: Unable to open file for writing: ", os.Args[2]+".tmp", err)
		os.Exit(-1)
	}
}
