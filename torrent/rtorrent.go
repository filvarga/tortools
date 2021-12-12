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

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/filvarga/vpptool/run"
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

	output, status := run.Run1("xmlrpc", r.url(), "download_list")

	if !status {
		log.Fatal("error getting xmlrpc result")
		return nil
	}

	re := regexp.MustCompile(`(?m)^\s+Index\s+[0-9]+\s+String:\s+'([0-9A-Z]{40})'\s*$`)
	results := re.FindAllStringSubmatch(string(output), -1)

	for i := 0; i < len(results); i++ {
		torrents = append(torrents, Torrent{r: r, hash: results[i][1]})
	}
	return torrents
}

func (r *Rtorrent) Add(magnet string) bool {
	re := regexp.MustCompile(`(?m)^\s+Integer:\s+([0-9]+)\s+$`)

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
	re := regexp.MustCompile(`(?m)^\s+String:\s+'(.+)'\s+$`)
	results := re.FindStringSubmatch(string(b))
	for i := 0; i < len(results); i++ {
		return results[i], nil
	}
	return nil, fmt.Errorf("Command failed, invalid value returned")
}

func getIntValue(b []byte) (int, error) {
	re := regexp.MustCompile(`(?m)^\s+Integer:\s+([0-9]+)\s+$`)
	results := re.FindStringSubmatch(string(b))
	for i := 0; i < len(results); i++ {
		return strconv.Atoi(results[i]), nil
	}
	return nil, fmt.Errorf("Command failed, invalid value returned")
}

func getInt64Value(b []byte) (int64, error) {
	// int64 !!!
	re := regexp.MustCompile(`(?m)^\s+64-bit\s+integer:\s+([0-9]+)\s+$`)
	results := re.FindStringSubmatch(string(b))
	for i := 0; i < len(results); i++ {
		return strconv.Atoi(results[i]), nil
	}
	return nil, fmt.Errorf("Command failed, invalid value returned")
}

func (t *Torrent) IsActive() bool {
	output, err := run.Run1("xmlrpc", t.r.url(), "d.is_active", t.hash)
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
	output, err := run.Run1("xmlrpc", t.r.url(), "d.complete", t.hash)
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
	output, err := run.Run1("xmlrpc", t.r.url(), "d.start", t.hash)
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
	output, err := run.Run1("xmlrpc", t.r.url(), "d.pause", t.hash)
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
	output, err := run.Run1("xmlrpc", t.r.url(), "d.stop", t.hash)
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
	output, err := run.Run1("xmlrpc", t.r.url(), "d.resume", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getIntValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value == 0

}

func (t *Torrent) GetFilename() string {
	output, err := run.Run1("xmlrpc", t.r.url(), "d.base_filename", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getStrValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) GetPath() string {
	output, err := run.Run1("xmlrpc", t.r.url(), "d.base_path", t.hash)
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
	output, err := run.Run1("xmlrpc", t.r.url(), "d.connection_seed", t.hash)
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
	output, err := run.Run1("xmlrpc", t.r.url(), "d.connection_leech", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getStrValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) GetBytesDone() int64 {
	output, err := run.Run1("xmlrpc", t.r.url(), "d.bytes_done", t.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getInt64Value(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (t *Torrent) GetBytesSize() int64 {
	output, err := run.Run1("xmlrpc", t.r.url(), "d.size_bytes", t.hash)
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
