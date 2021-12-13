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

package rprov

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Search struct {
	TV      bool
	Title   string
	Season  int
	Episode int
}

type Torrent struct {
	Title    string
	Link     string
	Magnet   string
	Seeders  int
	Leechers int
}

type Torrents []Torrent

func get(title string) []Torrent {
	var torrents []Torrent

	title = strings.ReplaceAll(strings.ToLower(title), " ", "-")
	url := fmt.Sprintf("https://www.magnetdl.com/%c/%s/se/desc", title[0], title)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	doc.Find("tbody tr").Each(func(_ int, tr *goquery.Selection) {
		var torrent Torrent

		tr.Find("td").Each(func(ix int, td *goquery.Selection) {

			switch td.AttrOr("class", "") {
			case "m":
				td.Find("a").Each(func(_ int, a *goquery.Selection) {
					if len(torrent.Magnet) == 0 {
						torrent.Magnet = a.AttrOr("href", "")
					}
				})
			case "n":
				td.Find("a").Each(func(_ int, a *goquery.Selection) {
					if len(torrent.Link) == 0 {
						torrent.Link = a.AttrOr("href", "")
					}
					if len(torrent.Title) == 0 {
						torrent.Title = a.AttrOr("title", "")
					}
				})
			case "s":
				if i, err := strconv.Atoi(td.Text()); err == nil {
					torrent.Seeders = i
				}
			case "l":
				if i, err := strconv.Atoi(td.Text()); err == nil {
					torrent.Leechers = i
				}
			}
		})

		if len(torrent.Title) > 0 && len(torrent.Magnet) > 0 &&
			torrent.Seeders > 0 {
			torrents = append(torrents, torrent)
		}
	})

	return torrents
}

func (s *Search) GetTorrents() []Torrent {
	if s.TV {
		return get(fmt.Sprintf("%s s%02de%02d", s.Title,
			s.Season, s.Episode))
	} else {
		return get(s.Title)
	}
}

/* vim: set ts=2: */
