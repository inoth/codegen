package util

import (
	"fmt"
	"os"
)

func Must(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
