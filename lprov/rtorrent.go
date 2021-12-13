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

type Torrents []Torrent

func (r *Rtorrent) url() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (r *Rtorrent) GetTorrents() Torrents {

	var torrents []Torrent

	output, err := run.Run2("xmlrpc", r.url(), "download_list")
	if err != nil {
		log.Fatal("error getting xmlrpc result")
	}

	re := regexp.MustCompile(`String: '(?P<str>[0-9A-Z]{40})'`)

	matches := re.FindAllSubmatch(output, -1)
	names := re.SubexpNames()

	for _, match := range matches {
		for i, group := range match {
			switch names[i] {
			case "str":
				torrents = append(torrents, Torrent{r: *r, hash: string(group)})
			}
		}
	}
	return torrents
}

func (tr Torrents) Resume() {
	for _, t := range tr {
		t.Resume()
	}
}

func (tr Torrents) Delete() {
	for _, t := range tr {
		t.Delete()
	}
}

func (tr Torrents) Start() {
	for _, t := range tr {
		t.Start()
	}
}

func (tr Torrents) Pause() {
	for _, t := range tr {
		t.Pause()
	}
}

func (tr Torrents) Stop() {
	for _, t := range tr {
		t.Stop()
	}
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

	return 0 == value
}

func str2Int(str string, def int) int {
	if i, err := strconv.Atoi(str); err == nil {
		return i
	}
	return def
}

func getStrValue(b []byte) (string, error) {
	re := regexp.MustCompile(`String: '(?P<str>.*)'`)

	matches := re.FindSubmatch(b)
	names := re.SubexpNames()

	for i, group := range matches {
		switch names[i] {
		case "str":
			return string(group), nil
		}
	}
	return "", fmt.Errorf("Command failed, invalid value returned")
}

func getIntValue(b []byte) (int, error) {
	re := regexp.MustCompile(`Integer: (?P<int>[0-9]+)`)

	matches := re.FindSubmatch(b)
	names := re.SubexpNames()

	for i, group := range matches {
		switch names[i] {
		case "int":
			if s, err := strconv.Atoi(string(group)); err == nil {
				return s, nil
			}
		}
	}
	return 0, fmt.Errorf("Command failed, invalid value returned")
}

func getInt64Value(b []byte) (int, error) {
	re := regexp.MustCompile(`64-bit integer: (?P<int>[0-9]*)`)

	matches := re.FindSubmatch(b)
	names := re.SubexpNames()

	for i, group := range matches {
		switch names[i] {
		case "int":
			if s, err := strconv.Atoi(string(group)); err == nil {
				return s, nil
			}
		}
	}
	return 0, fmt.Errorf("Command failed, invalid value returned")
}

func (t *Torrent) getStrValue(field string) string {
	var (
		output []byte
		err    error
		value  string
	)

	output, err = run.Run2("xmlrpc", t.r.url(), field, t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getStrValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) getIntValue(field string) int {
	var (
		output []byte
		err    error
		value  int
	)

	output, err = run.Run2("xmlrpc", t.r.url(), field, t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getIntValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) getInt64Value(field string) int {
	var (
		output []byte
		err    error
		value  int
	)

	output, err = run.Run2("xmlrpc", t.r.url(), field, t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getInt64Value(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) IsActive() bool {
	return 0 == t.getInt64Value("d.is_active")
}

func (t *Torrent) IsComplete() bool {
	return 0 == t.getInt64Value("d.complete")
}

func (t *Torrent) Resume() bool {
	return 0 == t.getIntValue("d.resume")
}

func (t *Torrent) Delete() bool {
	return 0 == t.getIntValue("d.erase")
}

func (t *Torrent) Start() bool {
	return 0 == t.getIntValue("d.start")
}

func (t *Torrent) Pause() bool {
	return 0 == t.getIntValue("d.pause")
}

func (t *Torrent) Stop() bool {
	return 0 == t.getIntValue("d.stop")
}

func (t *Torrent) GetName() string {
	return t.getStrValue("d.name")
}

func (t *Torrent) GetDirectory() string {
	return t.getStrValue("d.directory")
}

func (t *Torrent) GetSeeders() int {
	return str2Int(t.getStrValue("d.connection_seed"), -1)
}

func (t *Torrent) GetLeechers() int {
	return str2Int(t.getStrValue("d.connection_leech"), -1)
}

func (t *Torrent) GetBytesDone() int {
	return t.getInt64Value("d.bytes_done")
}

func (t *Torrent) GetBytesSize() int {
	return t.getInt64Value("d.size_bytes")
}

/* vim: set ts=2: */
