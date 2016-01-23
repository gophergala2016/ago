// Copyright (c) 2016 SeongJae Park.
//
// This program is free software; you can redistribute it and/or modify it
// under the terms of the GNU General Public License version 3 as published by
// the Free Software Foundation.

package main

import (
	"fmt"
	"os"
	"log"
)

const (
	CMD = "ago"
	USAGE = "USAGE: " + CMD + " <commands> [argument ...]\n"
	NOARG_ERRMSG = USAGE + "\nFor detail, try " + CMD + " help\n"
	HELP_MSG = "Use the source ;)\n"
)

var (
	agol = log.New(os.Stderr, "[ago] ", 0)
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("No argument.\n")
		fmt.Printf(NOARG_ERRMSG)
		os.Exit(1)
	}
	cmd := args[1]
	args = args[2:]
	switch cmd {
	case "help":
		fmt.Printf(HELP_MSG)
	default:
		agol.Printf("wrong commanad")
		os.Exit(1)
	}
}
