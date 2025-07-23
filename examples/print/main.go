package main

import (
	"fmt"
	"io"
	"os"

	"go.bytecodealliance.org/cm"

	"github.com/a-skua/go-wasi/internal/gen/wasi/cli/run"
)

func init() {
	run.Exports.Run = Run
}

func main() {
	Run()
}

func Run() cm.BoolResult {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <file>...\n", os.Args[0])
		os.Exit(1)
	}

	files := os.Args[1:]
	for i, filename := range files {
		if len(files) > 1 {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("==> %s <==\n", filename)
		}

		err := printFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s: %v\n", filename, err)
		}
	}

	return cm.ResultOK
}

func printFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	return err
}
