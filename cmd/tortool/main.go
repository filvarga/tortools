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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/filvarga/tortools/download"
	"github.com/filvarga/tortools/run"
	"github.com/filvarga/tortools/search"
)

var (
	context = ""
)

type user struct {
	idu int
	idg int
}

type mount struct {
	downloads string
	session   string
}

type image struct {
	name string
	ver  string
}

func delContainer(name string) error {
	return run.Exec3(true, "docker", "rm", "-f", name)
}

func (i *image) build(quiet bool, target string, u user) error {
	return run.Exec3(quiet, "docker", "build",
		"--build-arg", fmt.Sprintf("IDU:%d", u.idu),
		"--build-arg", fmt.Sprintf("IDG:%d", u.idg),
		"--network", "host", "--target", target,
		"-t", fmt.Sprintf("%s:%s", i.name, i.ver),
		context)
}

func (i *image) deploy(quiet bool, name string, m mount) error {
	_ = delContainer(name)
	return run.Exec3(quiet, "docker", "run", "-it",
		"-v", fmt.Sprintf("%s:/app/downloads", m.downloads),
		"-v", fmt.Sprintf("%s:/app/session", m.session),
		"-d", "--network", "host", "--name", name,
		fmt.Sprintf("%s:%s", i.name, i.ver))
}

func build() {
	var err error

	u := user{1000, 1000}
	app := image{"app", "latest"}
	web := image{"web", "latest"}

	err = app.build(false, "app", u)
	if err != nil {
		log.Fatal(err)
	}

	err = web.build(false, "web", u)
	if err != nil {
		log.Fatal(err)
	}
}

func deploy() {
	var err error

	m := mount{"/tmp/downloads", "/tmp/session"}
	app := image{"app", "latest"}
	web := image{"web", "latest"}

	err = app.deploy(false, "app", m)
	if err != nil {
		log.Fatal(err)
	}

	err = web.deploy(false, "web", m)
	if err != nil {
		log.Fatal(err)
	}
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

type arrayTags []string

func (i *arrayTags) String() string {
	return "string..."
}

func (i *arrayTags) Set(value string) error {
	*i = append(*i, strings.ToLower(value))
	return nil
}

func buildSearch() search.Search {
	s := search.Search{
		Season:  1,
		Episode: 1,
	}
	s.Title = flag.Arg(2)
	if len(s.Title) == 0 {
		printUsage()
	}
	if len(flag.Arg(3)) > 0 {
		s.Season = download.Str2Int(flag.Arg(3), -1)
		if len(flag.Arg(4)) > 0 {
			s.Episode = download.Str2Int(flag.Arg(4), -1)
		}
		if s.Season == -1 || s.Episode == -1 {
			printUsage()
		}
		s.Type = search.TV
	}
	return s
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
		case "list":
			ListAllDownloads(r)
		case "del":
		}
	case "search":
		switch flag.Arg(1) {
		default:
			printUsage()
		case "find":
		case "get":
		}
	case "find":
		switch flag.Arg(1) {
		default:
			printUsage()
		case "first":
			m := FindFirstB(r, buildSearch(), tags)
			fmt.Println(m.String())
		case "all":
			ListB(r, buildSearch(), tags)
		}
	case "get":
		switch flag.Arg(1) {
		default:
			printUsage()
		case "first":
			m := FindFirstB(r, buildSearch(), tags)
			fmt.Println(m.String())
			m.Get(r)
		case "all":
		}
	case "del":
		switch flag.Arg(1) {
		default:
			printUsage()
		case "first":
		case "all":
		}
	case "run":
	}
}

/* vim: set ts=2: */
