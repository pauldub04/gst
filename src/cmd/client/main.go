package main

import (
	"client/internal/client"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	filename := pflag.String("filename", "input", "")
	seed := pflag.Uint64("seed", 0, "")
	sizeInMB := pflag.Int32("size", 0, "")
	pflag.Parse()

	if *sizeInMB == 0 {
		fmt.Printf("Usage: %s --size <size_in_mb>\n", os.Args[0])
		return
	}

	client.Run(*filename, *seed, *sizeInMB, false, "local")
	fmt.Printf("Generated data to %s\n", *filename)
}
