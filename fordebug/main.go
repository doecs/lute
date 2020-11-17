package main

import (
	"fmt"
	"lute"
)

//
func main() {
	luteEngine := lute.New()
	dom := luteEngine.Md2VditorIRDOM("a\n\n\n\n\nb")
	fmt.Printf("%v", dom)
	if dom == "<p data-block=\"0\">a\n</p><p data-block=\"0\">b\n</p>" {
		fmt.Printf("equal")
	} else {
		fmt.Printf("no equal")
	}
	// fmt.Printf("test")
}
