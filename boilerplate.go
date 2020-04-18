package main

import (
	"fmt"
)

var LOG_LEVEL = "error"

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err);
	}
}

func debug(msg ...string) {
	if LOG_LEVEL == "debug" {
		fmt.Println(msg)
	}
}
