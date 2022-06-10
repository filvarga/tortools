/*
 * Copyright 2021 Filip Varga
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package torclient

import (
	"flag"
	"github.com/filvarga/tortools/download"
	"log"
	"os"
	"strings"
)

type arrayTags []string

func (i *arrayTags) String() string {
	return "string..."
}

func (i *arrayTags) Set(value string) error {
	*i = append(*i, strings.ToLower(value))
	return nil
}

func printUsage() {
	usage := `Usage:
%[1]s	download purge
%[1]s	download list
%[1]s [-tag <tag> ...] download del <title> [season] [episode]
%[1]s	[-tag <tag> ...] search find|get <title> [season] [episode]
%[1]s [-tag <tag> ...] find all|first <title> [season] [episode]
%[1]s	[-tag <tag> ...] get all|first <title> [season] [episode]
%[1]s	[-tag <tag> ...] del all|first <title> [season] [episode]
%[1]s run

Flags:
`
	log.Printf(usage, os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	var (
		tags arrayTags
		r    download.Rtorrent
	)

	flag.StringVar(&r.Host, "host", "localhost", "Provider host")
	flag.IntVar(&r.Port, "port", 80, "Provider port")
	flag.Var(&tags, "tag", "Contains tag")

	flag.Parse()

	switch flag.Arg(0) {
	default:
		printUsage()
	case "download":
		switch flag.Arg(1) {
		default:
			printUsage()
		case "purge":
			// purge all downloads
		case "list":
			// list all downloads
		case "del":
			// del all downloads matching search pattern
		}
	case "search":
		switch flag.Arg(1) {
		default:
			printUsage()
		case "find":
			// find all searches matching search pattern
		case "get":
			// get all searches matching search pattern
		}
	case "find":
		switch flag.Arg(1) {
		default:
			printUsage()
		case "first":
			// find first managed media matching search pattern
		case "all":
			// find all managed media matching search pattern
		}
	case "get":
		switch flag.Arg(1) {
		default:
			printUsage()
		case "first":
			// get first managed media matching search pattern
		case "all":
			// get all managed media matching search pattern
		}
	case "del":
		switch flag.Arg(1) {
		default:
			printUsage()
		case "first":
			// del first managed media matching search pattern
		case "all":
		}
	case "run":
	}
}

/* vim: set ts=2: */
