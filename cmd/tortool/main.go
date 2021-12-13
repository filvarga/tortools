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
	"regexp"
	"strings"

	"github.com/filvarga/tortools/lprov"
	"github.com/filvarga/tortools/rprov"
	"github.com/filvarga/tortools/run"
)

var (
	context = ""
)

type Torrent struct {
	Name       string
	Local      bool
	Downloaded bool
	Percent    float64
	Seeders    int
	Leechers   int
	local      *lprov.Torrent
	remote     *rprov.Torrent
}

type Torrents []Torrent

func (t *Torrent) get(r lprov.Rtorrent) {
	if !t.Local {
		r.Add(t.remote.Magnet)
	}
}

func (t *Torrent) show() {
	fmt.Printf("local: %5t se: %4d le: %4d name: %s\n",
		t.Local, t.Seeders, t.Leechers, t.Name)
}

func (tr Torrents) show() {
	for _, t := range tr {
		t.show()
	}
}

func convertl(t lprov.Torrent) *Torrent {
	percent := 100 / float64(t.GetBytesSize()) *
		float64(t.GetBytesDone())
	return &Torrent{
		Name:       t.GetName(),
		Local:      true,
		Downloaded: t.IsComplete(),
		Percent:    percent,
		Seeders:    t.GetSeeders(),
		Leechers:   t.GetLeechers(),
		local:      &t,
	}
}

func convertr(t rprov.Torrent) *Torrent {
	return &Torrent{
		Name:       t.Title,
		Local:      false,
		Downloaded: false,
		Percent:    0.0,
		Seeders:    t.Seeders,
		Leechers:   t.Leechers,
		remote:     &t,
	}
}

func convertlTorrents(tr lprov.Torrents) Torrents {
	var torrents Torrents
	for _, t := range tr {
		torrents = append(torrents, *convertl(t))
	}
	return torrents
}

func convertrTorrents(tr rprov.Torrents) Torrents {
	var torrents Torrents
	for _, t := range tr {
		torrents = append(torrents, *convertr(t))
	}
	return torrents
}

func contains(s string, substrs []string) bool {
	s = strings.ToLower(s)
	for _, substr := range substrs {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}

func listAllLocal(r lprov.Rtorrent) {
	torrents := convertlTorrents(r.GetTorrents())
	torrents.show()
}

func findAllLocal(r lprov.Rtorrent, s rprov.Search) lprov.Torrents {
	local := r.GetTorrents()
	if s.TV {
		return local.FindAll(regexp.MustCompile(
			fmt.Sprintf(`(?i)%s.*?s%02de%02d`,
				s.Title, s.Season, s.Episode)))
	} else {
		return local.FindAll(regexp.MustCompile(
			fmt.Sprintf(`(?i)%s`, s.Title)))
	}
}

func findAllRemote(s rprov.Search) rprov.Torrents {
	return s.GetTorrents()
}

func findAll(r lprov.Rtorrent, s rprov.Search) Torrents {

	var torrents Torrents

	for _, t := range convertlTorrents(findAllLocal(r, s)) {
		torrents = append(torrents, t)
	}

	for _, t := range convertrTorrents(findAllRemote(s)) {
		torrents = append(torrents, t)
	}

	return torrents
}

func findAllA(r lprov.Rtorrent, s rprov.Search, substrs []string) Torrents {

	var torrents Torrents

	if len(substrs) == 0 {
		return findAll(r, s)
	}

	for _, t := range findAll(r, s) {
		if contains(t.Name, substrs) {
			torrents = append(torrents, t)
		}
	}
	return torrents
}

func findFirst(r lprov.Rtorrent, s rprov.Search) *Torrent {

	local := findAllLocal(r, s)
	if len(local) > 0 {
		return convertl(local[0])
	}

	remote := findAllRemote(s)
	if len(remote) > 0 {
		return convertr(remote[0])
	}
	return nil
}

func findFirstA(r lprov.Rtorrent, s rprov.Search, substrs []string) *Torrent {
	if len(substrs) > 0 {
		for _, t := range convertlTorrents(findAllLocal(r, s)) {
			if contains(t.Name, substrs) {
				return &t
			}
		}

		for _, t := range convertrTorrents(findAllRemote(s)) {
			if contains(t.Name, substrs) {
				return &t
			}
		}
	} else {
		return findFirst(r, s)
	}
	return nil
}

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
	return run.Run3(true, "docker", "rm", "-f", name)
}

func (i *image) build(quiet bool, target string, u user) error {
	return run.Run3(quiet, "docker", "build",
		"--build-arg", fmt.Sprintf("IDU:%d", u.idu),
		"--build-arg", fmt.Sprintf("IDG:%d", u.idg),
		"--network", "host", "--target", target,
		"-t", fmt.Sprintf("%s:%s", i.name, i.ver),
		context)
}

func (i *image) deploy(quiet bool, name string, m mount) error {
	delContainer(name)
	return run.Run3(quiet, "docker", "run", "-it",
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

func print_usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s: <list>|<<find|get> <title> [Season] [Episode]>\n",
		os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

type arraySubstrs []string

func (i *arraySubstrs) String() string {
	return "string..."
}

func (i *arraySubstrs) Set(value string) error {
	*i = append(*i, strings.ToLower(value))
	return nil
}

func main() {

	r := lprov.Rtorrent{
		Host: "localhost",
		Port: 80,
	}

	s := rprov.Search{}

	var substrs arraySubstrs

	flag.Var(&substrs, "contains", "Contains")
	flag.BoolVar(&s.TV, "tv", false, "TV Show")

	flag.Parse()

	s.Title = flag.Arg(1)

	switch flag.Arg(0) {
	default:
		print_usage()
	case "list":
		listAllLocal(r)
	case "find":
		if len(s.Title) == 0 {
			print_usage()
		}
		if s.TV {
			s.Season = lprov.Str2Int(flag.Arg(2), -1)
			s.Episode = lprov.Str2Int(flag.Arg(3), -1)
			if s.Season == -1 || s.Episode == -1 {
				print_usage()
			}
		}
		torrents := findAllA(r, s, substrs)
		torrents.show()
	case "get":
		if len(s.Title) == 0 {
			print_usage()
		}
		if s.TV {
			s.Season = lprov.Str2Int(flag.Arg(2), -1)
			s.Episode = lprov.Str2Int(flag.Arg(3), -1)
			if s.Season == -1 || s.Episode == -1 {
				print_usage()
			}
		}
		torrent := findFirstA(r, s, substrs)
		torrent.show()
		torrent.get(r)
	}
}

/* vim: set ts=2: */
