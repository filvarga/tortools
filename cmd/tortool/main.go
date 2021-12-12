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

	"github.com/filvarga/tortools/run"
)

var (
	context = ""
	git_url = ""
)

func filter(src []Torrent, substr []string) []Torrent {
	var dst []Torrent

	if len(substr) == 0 {
		return src
	}

	for i := 0; i < len(src); i++ {

		for j := 0; j < len(substr); j++ {
			if strings.Contains(src[i].Title, substr[j]) {
				dst = append(dst, src[i])
			}
		}
	}

	return dst
}

type Torrent struct {
	Name       string
	Local      bool
	Downloaded bool
	Percent    float
	Seeders    string
	Leechers   string
}

func buildRegex(s magnetdl.Search) regexp.Regex {
	if m.TVShow {
		s := fmt.Sprintf(`^.+%s.+[sS]%02d:[eE]%02d.+$`,
			m.Title, m.Season, m.Episode)
	} else {
		s := fmt.Sprintf(`^.+%s.+$`, m.Title)
	}
	return regexp.MustCompile(s)
}

func getFirst(r rtorrent.Rtorrent, s magnetdl.Search) (Torrent, error) {
	var torrent Torrent

	// TODO: call findFirst and just go to the business of downloading

	return torrent, nil
}

func findFirst(r rtorrent.Rtorrent, s magnetdl.Search) (Torrent, error) {
	var torrent Torrent

	// TODO: sort

	return torrent, nil
}

func findAll(r rtorrent.Rtorrent, s magnetdl.Search) ([]Torrent, error) {
	var torrents []Torrent

	local := r.GetTorrents()
	re := buildRegex(s)

	for i := 0; i < len(local); i++ {
		name := local[i].GetFileName()
		if re.MatchString(name) {
			percent := 100 / local[i].GetBytesSize() *
				local[i].GetBytesDone()
			torrents = append(torrents, Torrent{
				Name:       name,
				Local:      true,
				Downloaded: local[i].IsComplete(),
				Percent:    percent,
				Seeders:    local[i].GetSeeders(),
				Leechers:   local[i].GetLeechers(),
			})
		}
	}

	remote := s.GetTorrents()

	for i := 0; i < len(remote); i++ {
		torrents = append(torrents, Torrent{
			Name:       remote[i].Title,
			Local:      false,
			Downloaded: false,
			Percent:    0.0,
			Seeders:    remote[i].Seeders,
			Leechers:   remote[i].Leechers,
		})
	}

	return torrents, nil
}

// TODO: return new type of generic Torrent ??
// like for example with stuff like Downloaded && Local and so on
func (r rtorrent) FindAllTorrents(m magnetdl.Media) []rtorrent.Torrent {

	// .Stop()
	// .Start()
	// .IsActive()
	// .GetFileName()
	// .IsComplete()
	// .GetSeeders()
	// .GetLeechers()
	// .GetBytesDone()
	// .GetBytesSize()

	// .Title
	// .Link
	// .Magnet
	// .Seeders
	// .Leechers

	torrents := r.GetTorrents()

	for i := 0; i <= len(torrents); i++ {
		name := torrents[i].GetFilename()
		if re.MatchString(name) {
			// append to some ...
		}
	}
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
		"-v", fmt.Sprintf("%s:/app/downloads", downloads),
		"-v", fmt.Sprintf("%s:/app/session", session),
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
	fmt.Fprintf(os.Stderr, "Usage of %s: <show|other> <list|download> <title>\n",
		os.Args[0])
	flag.PrintDefaults()
}

func main() {

	substr := []string{}

	s := flag.Int("s", 1, "Season")
	e := flag.Int("e", 1, "Episode")
	b := flag.Bool("show", false, "TV Show")

	p480 := flag.Bool("480p", false, "")
	p720 := flag.Bool("720p", false, "")
	p1080 := flag.Bool("1080p", false, "")

	flag.Parse()

	if *p480 {
		substr = append(substr, "480p")
	}
	if *p720 {
		substr = append(substr, "720p")
	}
	if *p1080 {
		substr = append(substr, "1080p")
	}

	switch flag.Arg(0) {
	default:
		print_usage()
	case "other":
		m.Type = Other
	case "show":
		m.Type = Show
	}

	media := magnetdl.Media{
		Title:   flag.Arg(2),
		Seasson: s,
		Episode: e,
		Show:    b}

	if len(m.Title) == 0 {
		print_usage()
		return
	}

	switch flag.Arg(1) {
	default:
		print_usage()
	case "list":
		torrents := FindAllTorrents(m, contains)
		for i := 0; i < len(torrents); i++ {
			fmt.Printf("se: %s\tle: %s\ttitle: %s\n",
				torrents[i].Seeders, torrents[i].Leechers,
				torrents[i].Title)
		}
	case "download":
		torrents := FindAllTorrents(m, contains)
		if len(torrents) > 0 {
			// TODO: highest seed rate and > 0
			//r.add(torrents[0])
		}
	}
}

/* vim: set ts=2: */
