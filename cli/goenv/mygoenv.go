package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("My go version is %v running on %v", runtime.Version(), runtime.GOOS)
}
