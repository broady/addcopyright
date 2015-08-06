// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	apache         = flag.Bool("apache2", false, "Use Apache 2.0 header")
)

const apache2Header = `/` + `/ Copyright %d %s. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

`

func usage(err string) {
	fmt.Fprintf(os.Stderr, "%s. usage:\n", err)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Parse()

	var head []byte

	if !*apache && *headerFile == "" {
		usage("must provide a license type: -apache2 or custom (-header)")
	}

	if *headerFile == "-" {
		h, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("could not from stdin: %v", err)
		}
		head = h
	} else if *headerFile != "" {
		h, err := ioutil.ReadFile(*headerFile)
		if err != nil {
			log.Fatalf("could not read header file: %v", err)
		}
		head = h
	}

	if *apache {
		if *copyrightOwner == "" {
			usage("must provide -owner with -apache2")
		}

		head = []byte(fmt.Sprintf(apache2Header, time.Now().Year(), *copyrightOwner))
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
