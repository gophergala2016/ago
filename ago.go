// Copyright (c) 2016 SeongJae Park.
//
// This program is free software; you can redistribute it and/or modify it
// under the terms of the GNU General Public License version 3 as published by
// the Free Software Foundation.

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
)

const (
	CMD          = "ago"
	USAGE        = "USAGE: " + CMD + " <commands> [argument ...]\n"
	NOARG_ERRMSG = USAGE + "\nFor detail, try " + CMD + " help\n"
	HELP_MSG     = "Use the source ;)\n"
	ANDRD        = "android"
	ANDRD_TMPDIR = "/data/local/tmp"
	DOCDIR       = "docs"
	DOCINFO      = "info"
	WORDINFO     = "words"
)

var (
	errl        = log.New(os.Stderr, "[err] ", 0)
	dbgl        = log.New(os.Stderr, "[dbg] ", 0)
	metadat_dir = "/tmp/.ago"
	docs_dir    = ""
)

// main is the entry point of `ago`.
// ago usage is similar to familiar tools:
// 	ago <command> [argument ...]
//
// commands are:
// ls-docs, add-docs, rm-docs: list, add, remove documentation[s].
// test: start a test. Number of questions can be specified as option.
//
// The description above is lie because this program is nothing for now. It is
// just a plan.
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
	case "ls-docs":
		fmt.Printf("list docs\n")
	case "add-docs":
		fmt.Printf("add doc\n")
	case "rm-docs":
		fmt.Printf("remove doc\n")
	case "test":
		fmt.Printf("let the game begin\n")
	case "help":
		fmt.Printf(HELP_MSG)
	default:
		errl.Printf("wrong commanad")
		os.Exit(1)
	}
}

// init initializes few things for `ago`.
// Internally, ago uses a metadata directory for state saving. Path of the
// directory is `$HOME/.ago`. If `$HOME` is not exists, `/tmp` is used as
// default. For Android support in future, it should be `/data/local/tmp` at
// future.
//
// Hierarchy of the directory is as:
// .ago/docs/doc1
//          /info
//     /words
//
// Documents added by user resides under .ago/docs/ with its own directory. The
// document own directories are named as doc[id] which id is an integer.
// Metadata about those documents are recorded under .ago/docs/info file. The
// metadata contains original document name and current location under .ago
// directory. Because current ago support only text file, this structure is
// unnecessary. Actually, the struct is for future scaling. In future, ago will
// be an document organizer like Mendeley[1] and will support not only text
// file, but also pdf, odt, url, etc.
//
// File `words` under the .ago/ directory contains all data for words in the
// documents. It contains each word and its frequency in the docs(in total and
// per each doc), score of user, and meaning of the word. Calculated importance
// can be in there maybe but not yet decided to add it.
//
// [1] https://www.mendeley.com/
func init() {
	if runtime.GOOS == ANDRD {
		metadat_dir = ANDRD_TMPDIR
	}

	if os.Getenv("HOME") != "" {
		metadat_dir = os.Getenv("HOME")
	}
	metadat_dir = path.Join(metadat_dir, ".ago")
	dbgl.Printf("metadata dir is at %s\n", metadat_dir)

	docs_dir = path.Join(metadat_dir, DOCDIR)
	// docs dir is already exists.
	if _, err := os.Stat(docs_dir); err == nil {
		return
	}

	dbgl.Printf("docs dir is not exists. Create it.\n")
	err := os.MkdirAll(docs_dir, 0700)
	if err != nil {
		errl.Printf("docs dir %s creation failed: %s\n", docs_dir, err)
		os.Exit(1)
	}

	f, err := os.Create(path.Join(docs_dir, DOCINFO))
	if err != nil {
		errl.Printf("docs info file creation failed: %s\n", err)
		os.Exit(1)
	}
	f.Close()

	f, err = os.Create(path.Join(metadat_dir, WORDINFO))
	if err != nil {
		errl.Printf("words info file creation failed: %s\n", err)
		os.Exit(1)
	}
	f.Close()
}
