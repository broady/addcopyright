// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command addcopyright adds copyright headers to several files at once.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	headerFile     = flag.String("header", "", "Path to file containing full header text")
	sentinelString = flag.String("sentinel", "/"+"/ Copyright", "String to look for")
	copyrightOwner = flag.String("owner", "", "Copyright owner")
)

const defaultHeader = `/` + `/ Copyright %d %s. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

`

func main() {
	flag.Parse()

	if *copyrightOwner == "" && *headerFile == "" {
		fmt.Fprintln(os.Stderr, "must provide -owner or -header flag. usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	head := []byte(fmt.Sprintf(defaultHeader, time.Now().Year(), *copyrightOwner))

	if *headerFile != "" {
		h, err := ioutil.ReadFile(*headerFile)
		if err != nil {
			log.Fatalf("could not read header file: %v", err)
		}
		head = h
	}

	for _, file := range flag.Args() {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("could not read %s: %v", file, err)
		}
		if !strings.Contains(string(f), *sentinelString) {
			log.Printf("prepending header to %s", file)
			err := ioutil.WriteFile(file, append(head, f...), 0)
			if err != nil {
				log.Fatalf("could not write %s: %v. existing contents were: %#q", file, err)
			}
		}
	}
	log.Printf("all done")
}
