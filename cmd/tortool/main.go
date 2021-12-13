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

	"github.com/filvarga/tortools/lprov"
	"github.com/filvarga/tortools/rprov"
	"github.com/filvarga/tortools/run"
)

var (
	context = ""
)

/*
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
*/

type Torrent struct {
	Name       string
	Local      bool
	Downloaded bool
	Percent    float64
	Seeders    int
	Leechers   int
}

type Torrents []Torrent

func (tr Torrents) show() {
	for _, t := range tr {
		fmt.Printf("local: %5t se: %4d le: %4d name: %s\n",
			t.Local, t.Seeders, t.Leechers, t.Name)
	}
}

func buildRegex(s rprov.Search) *regexp.Regexp {
	var str string
	if s.TVShow {
		str = fmt.Sprintf(`^.+%s.+[sS]%02d:[eE]%02d.+$`,
			s.Title, s.Season, s.Episode)
	} else {
		str = fmt.Sprintf(`^.+%s.+$`, s.Title)
	}
	return regexp.MustCompile(str)
}

func getFirst(r lprov.Rtorrent, s rprov.Search) (Torrent, error) {
	var torrent Torrent

	// TODO: call findFirst and just go to the business of downloading

	return torrent, nil
}

func findFirst(r lprov.Rtorrent, s rprov.Search) (Torrent, error) {
	var torrent Torrent

	// TODO: sort

	return torrent, nil
}

func listAll(r lprov.Rtorrent) (Torrents, error) {
	var torrents Torrents
	local := r.GetTorrents()

	for i := 0; i < len(local); i++ {
		percent := 100 / float64(local[i].GetBytesSize()) *
			float64(local[i].GetBytesDone())
		torrents = append(torrents, Torrent{
			Name:       local[i].GetName(),
			Local:      true,
			Downloaded: local[i].IsComplete(),
			Percent:    percent,
			Seeders:    local[i].GetSeeders(),
			Leechers:   local[i].GetLeechers(),
		})
	}

	return torrents, nil
}

func findAll(r lprov.Rtorrent, s rprov.Search) (Torrents, error) {
	var torrents []Torrent

	local := r.GetTorrents()
	re := buildRegex(s)

	for i := 0; i < len(local); i++ {
		name := local[i].GetName()
		if re.MatchString(name) {
			percent := 100 / float64(local[i].GetBytesSize()) *
				float64(local[i].GetBytesDone())
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
	fmt.Fprintf(os.Stderr, "Usage of %s: <list>|<<find|get> <title>>\n",
		os.Args[0])
	flag.PrintDefaults()
}

func main() {

	r := lprov.Rtorrent{
		Host: "localhost",
		Port: 80,
	}

	s := rprov.Search{}

	flag.IntVar(&s.Season, "s", 1, "Season")
	flag.IntVar(&s.Episode, "e", 1, "Episode")
	flag.BoolVar(&s.TVShow, "show", false, "TV Show")

	flag.Parse()

	s.Title = flag.Arg(1)

	switch flag.Arg(0) {
	default:
		print_usage()
	case "list":
		torrents, _ := listAll(r)
		torrents.show()
	case "find":
		if len(s.Title) == 0 {
			print_usage()
		}
		torrents, _ := findAll(r, s)
		torrents.show()
	case "get":
		if len(s.Title) == 0 {
			print_usage()
		}
	}
}

/* vim: set ts=2: */
