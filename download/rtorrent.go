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

package download

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/filvarga/tortools/run"
)

type Rtorrent struct {
	Host string
	Port int
}

type Download struct {
	r    Rtorrent
	hash string
}

type Downloads []Download

func (r *Rtorrent) getInt64Value(field string) int {
	var (
		output []byte
		err    error
		value  int
	)

	output, err = run.Run2("xmlrpc", r.url(), field)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getInt64Value(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (r *Rtorrent) url() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (r *Rtorrent) GetSystemTime() int {
	return r.getInt64Value("system.time")
}

func (r *Rtorrent) GetDownloads() Downloads {

	var downloads []Download

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
				downloads = append(downloads, Download{r: *r, hash: string(group)})
			}
		}
	}
	return downloads
}

func (r *Rtorrent) TryCreateDownload(magnet string) bool {
	var value int

	output, err := run.Run2("xmlrpc", r.url(), "load.start", "", magnet)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getIntValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return 0 == value
}

func diffDownloads(dn Downloads, do Downloads) Downloads {
	var (
		found bool
		out   Downloads
	)

	for _, n := range dn {
		found = false
		for _, o := range do {
			if o.hash == n.hash {
				found = true
				break
			}
		}
		if !found {
			out = append(out, n)
		}
	}
	return out
}

func (r *Rtorrent) AddDownload(magnet string) *Download {

	// this isn't bulet proof, if for instance some other
	// tool or rtorrent (from watch dir) itself adds
	// torrent

	retries := 2

	old := r.GetDownloads()
	t1 := r.GetSystemTime()

	if !r.TryCreateDownload(magnet) {
		return nil
	}

	t2 := r.GetSystemTime()

	for i := 0; i < retries; i++ {
		diff := diffDownloads(r.GetDownloads(), old)
		if len(diff) > 1 {
			log.Fatal("error multiple new downloads found")
		}
		if len(diff) == 1 {
			t := diff[0].GetLoadDate()
			if t > t1 && t < t2 {
				return &diff[0]
			}
		}
		time.Sleep(time.Second)
	}

	// timed out ...
	return nil
}

func (ds Downloads) Resume() {
	for _, d := range ds {
		d.Resume()
	}
}

func (ds Downloads) Delete() {
	for _, d := range ds {
		d.Delete()
	}
}

func (ds Downloads) Start() {
	for _, d := range ds {
		d.Start()
	}
}

func (ds Downloads) Pause() {
	for _, d := range ds {
		d.Pause()
	}
}

func (ds Downloads) Stop() {
	for _, d := range ds {
		d.Stop()
	}
}

func (ds Downloads) FindAll(re *regexp.Regexp) Downloads {
	var downloads Downloads
	for _, d := range ds {
		if re.MatchString(d.GetName()) {
			downloads = append(downloads, d)
		}
	}
	return downloads
}

func (ds Downloads) FindFirst(re *regexp.Regexp) *Download {
	for _, d := range ds {
		if re.MatchString(d.GetName()) {
			return &d
		}
	}
	return nil
}

func Str2Int(str string, def int) int {
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

func (d *Download) getStrValue(field string) string {
	var (
		output []byte
		err    error
		value  string
	)

	output, err = run.Run2("xmlrpc", d.r.url(), field, d.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getStrValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (d *Download) getIntValue(field string) int {
	var (
		output []byte
		err    error
		value  int
	)

	output, err = run.Run2("xmlrpc", d.r.url(), field, d.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getIntValue(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (d *Download) getInt64Value(field string) int {
	var (
		output []byte
		err    error
		value  int
	)

	output, err = run.Run2("xmlrpc", d.r.url(), field, d.hash)
	if err != nil {
		log.Fatal(err)
	}

	value, err = getInt64Value(output)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

func (d *Download) IsActive() bool {
	return 0 == d.getInt64Value("d.is_active")
}

func (d *Download) IsStarted() bool {
	return 1 == d.getInt64Value("d.state")
}

func (d *Download) IsComplete() bool {
	return 0 == d.getInt64Value("d.complete")
}

func (d *Download) Resume() bool {
	return 0 == d.getIntValue("d.resume")
}

func (d *Download) Delete() bool {
	return 0 == d.getIntValue("d.erase")
}

func (d *Download) Start() bool {
	return 0 == d.getIntValue("d.start")
}

func (d *Download) Pause() bool {
	return 0 == d.getIntValue("d.pause")
}

func (d *Download) Stop() bool {
	return 0 == d.getIntValue("d.stop")
}

func (d *Download) GetName() string {
	return d.getStrValue("d.name")
}

func (d *Download) GetDirectory() string {
	return d.getStrValue("d.directory")
}

func (d *Download) GetSeeders() int {
	return Str2Int(d.getStrValue("d.connection_seed"), -1)
}

func (d *Download) GetLeechers() int {
	return Str2Int(d.getStrValue("d.connection_leech"), -1)
}

func (d *Download) GetBytesDone() int {
	return d.getInt64Value("d.bytes_done")
}

func (d *Download) GetBytesSize() int {
	return d.getInt64Value("d.size_bytes")
}

func (d *Download) GetLoadDate() int {
	return d.getInt64Value("d.load_date")
}

func (d *Download) GetPercentDone() float64 {
	return 100 / float64(d.GetBytesSize()) *
		float64(d.GetBytesDone())
}

/* vim: set ts=2: */
