package main

import (
	"fmt"
	"lemin/internal"
	"time"
)

func main() {
	start := time.Now()

	file := internal.GetFile()
	internal.ValidateFileFormat(file)
	internal.ValidateConnectivity()
	internal.FindAllPaths()
	internal.FindBestPaths()
	internal.Simulate()

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}
