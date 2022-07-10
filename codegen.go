package main

import (
	"fmt"
	"os"

	"codegenfornet/cmd"
	_ "codegenfornet/src/plugins/struct_crawl/all"
)

func main() {
	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				fmt.Printf("%v\n", err)
			} else {
				panic(exception)
			}
			os.Exit(1)
		}
	}()
	cmd.Execute()
}
