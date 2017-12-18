package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expect a single argument to underline.  Can be any normal string")
	} else {
		arg := []byte(os.Args[1])
		strike, underline, both := "", "", ""
		for _, b := range arg {
			strike += string([]byte{b, 204, 182})
			underline += string([]byte{b, 205, 159})
			both += string([]byte{b, 204, 182, 205, 159})
		}
		fmt.Printf("Strikethrough: %s\nUnderline: %s\nBoth: %s\n", strike, underline, both)
	}
}
