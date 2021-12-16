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
	"fmt"
	"regexp"
	"strings"

	"github.com/filvarga/tortools/download"
	"github.com/filvarga/tortools/search"
)

type Media struct {
	Name     string
	Local    bool
	Type     int
	Season   int
	Episode  int
	torrent  *search.Torrent
	download *download.Download
}

type Medias []Media

func (m *Media) String() string {
	if m.Local {
		return fmt.Sprintf("local:  %s", m.Name)
	} else {
		return fmt.Sprintf("remote: %s", m.Name)
	}
}

func (m *Media) Show() {
	fmt.Println(m.String())
}

func (m *Media) Get(r download.Rtorrent) bool {
	if !m.Local {
		d := r.AddDownload(m.torrent.Magnet)
		if d == nil {
			return false
		}
		// local has precedence over remote
		m.Name = d.GetName()
		m.Local = true
		m.download = d
	}
	return true
}

func (m *Media) Del() bool {
	if m.Local {
		return m.download.Delete()
	}
	return false
}

func (ms Medias) Show() {
	for _, m := range ms {
		m.Show()
	}
}

func (ms Medias) Get(r download.Rtorrent) {
	for _, m := range ms {
		m.Get(r)
	}
}

func (ms Medias) Del() {
	for _, m := range ms {
		m.Del()
	}
}

func convertTorrent(t search.Torrent) *Media {
	return &Media{
		Name:    t.Title,
		Local:   false,
		torrent: &t,
	}
}

func convertDownload(d download.Download) *Media {
	return &Media{
		Name:     d.GetName(),
		Local:    true,
		download: &d,
	}
}

func convertTorrents(ts search.Torrents) Medias {
	var ms Medias
	for _, t := range ts {
		ms = append(ms, *convertTorrent(t))
	}
	return ms
}

func convertDownloads(ds download.Downloads) Medias {
	var ms Medias
	for _, d := range ds {
		ms = append(ms, *convertDownload(d))
	}
	return ms
}

func findTorrents(s search.Search) search.Torrents {
	return s.GetTorrents()
}

func findDownloads(r download.Rtorrent, s search.Search) download.Downloads {
	ds := r.GetDownloads()
	if s.Type == search.TV {
		return ds.FindAll(regexp.MustCompile(
			fmt.Sprintf(`(?i)%s.*?s%02de%02d`,
				s.Title, s.Season, s.Episode)))
	} else {
		return ds.FindAll(regexp.MustCompile(
			fmt.Sprintf(`(?i)%s`, s.Title)))
	}
}

func ListAllDownloads(r download.Rtorrent) Medias {
	return convertDownloads(r.GetDownloads())
}

func FindTorrentsA(s search.Search) Medias {
	return convertTorrents(findTorrents(s))
}

func FindTorrentsB(s search.Search, tags []string) Medias {
	var ms Medias
	for _, m := range FindTorrentsA(s) {
		if contains(m.Name, tags) {
			ms = append(ms, m)
		}
	}
	return ms
}

func FindDownloadsA(r download.Rtorrent, s search.Search) Medias {
	return convertDownloads(findDownloads(r, s))
}

func FindDownloadsB(r download.Rtorrent, s search.Search, tags []string) Medias {
	var ms Medias
	for _, m := range FindDownloadsA(r, s) {
		if contains(m.Name, tags) {
			ms = append(ms, m)
		}
	}
	return ms
}

func FindFirstDownloadA(r download.Rtorrent, s search.Search) *Media {
	downloads := findDownloads(r, s)
	if len(downloads) > 0 {
		return convertDownload(downloads[0])
	}
	return nil
}

func FindFirstDownloadB(r download.Rtorrent, s search.Search, tags []string) *Media {
	if len(tags) > 0 {
		for _, m := range convertDownloads(findDownloads(r, s)) {
			if contains(m.Name, tags) {
				return &m
			}
		}
	} else {
		return FindFirstDownloadA(r, s)
	}
	return nil
}

func FindAllA(r download.Rtorrent, s search.Search) Medias {
	var ms Medias
	for _, t := range convertDownloads(findDownloads(r, s)) {
		ms = append(ms, t)
	}
	for _, t := range convertTorrents(findTorrents(s)) {
		ms = append(ms, t)
	}
	return ms
}

func FindAllB(r download.Rtorrent, s search.Search, tags []string) Medias {
	var ms Medias
	for _, m := range FindAllA(r, s) {
		if contains(m.Name, tags) {
			ms = append(ms, m)
		}
	}
	return ms
}

func FindFirstA(r download.Rtorrent, s search.Search) *Media {
	downloads := findDownloads(r, s)
	if len(downloads) > 0 {
		return convertDownload(downloads[0])
	}
	torrents := findTorrents(s)
	if len(torrents) > 0 {
		return convertTorrent(torrents[0])
	}
	return nil
}

func FindFirstB(r download.Rtorrent, s search.Search, tags []string) *Media {
	if len(tags) > 0 {
		for _, m := range convertDownloads(findDownloads(r, s)) {
			if contains(m.Name, tags) {
				return &m
			}
		}
		for _, m := range convertTorrents(findTorrents(s)) {
			if contains(m.Name, tags) {
				return &m
			}
		}
	} else {
		return FindFirstA(r, s)
	}
	return nil
}

func contains(s string, tags []string) bool {
	s = strings.ToLower(s)
	for _, substr := range tags {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}

/* vim: set ts=2: */
