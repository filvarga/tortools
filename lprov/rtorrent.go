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

package lprov

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/filvarga/tortools/run"
)

type Rtorrent struct {
	Host string
	Port int
}

type Torrent struct {
	r    Rtorrent
	hash string
}

func (r *Rtorrent) url() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (r *Rtorrent) GetTorrents() []Torrent {

	var torrents []Torrent

	output, err := run.Run2("xmlrpc", r.url(), "download_list")
	if err != nil {
		log.Fatal("error getting xmlrpc result")
	}

	// TODO: redo to subgroups
	re := regexp.MustCompile(`(?m)^\s+Index\s+[0-9]+\s+String:\s+'([0-9A-Z]{40})'\s*$`)
	results := re.FindAllStringSubmatch(string(output), -1)

	for i := 0; i < len(results); i++ {
		torrents = append(torrents, Torrent{r: *r, hash: results[i][1]})
	}
	return torrents
}

func (r *Rtorrent) Add(magnet string) bool {
	var value int

	output, err := run.Run2("xmlrpc", r.url(), "load.normal", "", magnet)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getIntValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value == 0
}

func getStrValue(b []byte) (string, error) {
	re := regexp.MustCompile(`String: '(?P<str>.*)'`)

	matches := re.FindStringSubmatch(string(b))
	names := re.SubexpNames()

	for i, match := range matches {
		switch names[i] {
		case "str":
			return match, nil
		}
	}
	return "", fmt.Errorf("Command failed, invalid value returned")
}

func getIntValue(b []byte) (int, error) {
	re := regexp.MustCompile(`Integer: (?P<int>[0-9]+)`)

	matches := re.FindStringSubmatch(string(b))
	names := re.SubexpNames()

	for i, match := range matches {
		switch names[i] {
		case "int":
			if s, err := strconv.Atoi(match); err == nil {
				return s, nil
			}
		}
	}
	return 0, fmt.Errorf("Command failed, invalid value returned")
}

func getInt64Value(b []byte) (int, error) {
	re := regexp.MustCompile(`64-bit integer: (?P<int>[0-9]*)`)

	matches := re.FindStringSubmatch(string(b))
	names := re.SubexpNames()

	for i, match := range matches {
		switch names[i] {
		case "int":
			if s, err := strconv.Atoi(match); err == nil {
				return s, nil
			}
		}
	}
	return 0, fmt.Errorf("Command failed, invalid value returned")
}

func (t *Torrent) IsActive() bool {
	var value int

	output, err := run.Run2("xmlrpc", t.r.url(), "d.is_active", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getInt64Value(output)
	if err != nil {
		log.Fatal(err)
	}

	return value == 0
}

func (t *Torrent) IsComplete() bool {
	var value int

	output, err := run.Run2("xmlrpc", t.r.url(), "d.complete", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getInt64Value(output)
	if err != nil {
		log.Fatal(err)
	}

	return value == 0
}

func (t *Torrent) Start() bool {
	var value int

	output, err := run.Run2("xmlrpc", t.r.url(), "d.start", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getIntValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value == 0
}

func (t *Torrent) Pause() bool {
	var value int

	output, err := run.Run2("xmlrpc", t.r.url(), "d.pause", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getIntValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value == 0
}

func (t *Torrent) Stop() bool {
	var value int

	output, err := run.Run2("xmlrpc", t.r.url(), "d.stop", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getIntValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value == 0
}

func (t *Torrent) Resume() bool {
	var value int

	output, err := run.Run2("xmlrpc", t.r.url(), "d.resume", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getIntValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value == 0

}

func (t *Torrent) GetName() string {
	var value string

	output, err := run.Run2("xmlrpc", t.r.url(), "d.name", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getStrValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) GetDirectory() string {
	var value string

	output, err := run.Run2("xmlrpc", t.r.url(), "d.directory", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getStrValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) GetSeeders() string {
	var value string

	output, err := run.Run2("xmlrpc", t.r.url(), "d.connection_seed", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getStrValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) GetLeechers() string {
	var value string

	output, err := run.Run2("xmlrpc", t.r.url(), "d.connection_leech", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getStrValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) GetBytesDone() int {
	var value int

	output, err := run.Run2("xmlrpc", t.r.url(), "d.bytes_done", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getInt64Value(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) GetBytesSize() int {
	var value int

	output, err := run.Run2("xmlrpc", t.r.url(), "d.size_bytes", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getInt64Value(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

/* vim: set ts=2: */
