/*
 * Copyright 2018 by Dave Barach 
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
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
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func scrape() {
	var string_data, login_url, after_tag string
	var correctables, uncorrectables string
	var data []byte
	var resp *http.Response
	var err error
	var start_loc, end_loc int
	var now time.Time
	var year, month, day, hour, minute, second int
	
	/* Log in */
	login_url = "http://192.168.100.1/goform/login"
	form := url.Values {
		"loginUsername": {"admin"},
		"loginPassword": {"password"},
	}
	resp, err = http.PostForm (login_url, form)
	if err != nil {
		log.Fatal (err)
	}
	_, err = ioutil.ReadAll (resp.Body)
	resp.Body.Close()

	/* Try to get useful data */
	resp, err = http.Get ("http://192.168.100.1/RgConnect.asp")
	if err != nil {
		log.Fatal (err)
	}
	data, err = ioutil.ReadAll (resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal (err)
	}

	/* Turn byte array into a string */
	string_data = string(data)

	/* Search for the summary stats */
	start_loc = strings.Index (string_data, "Total Uncorrectables")

	if start_loc == -1 {
		log.Fatal ("Missing <th> Total Uncorrectables")
	}
	
	/* Search for the correctable error count, in the first <td> */
	after_tag = string_data[start_loc:len(string_data)]
	start_loc = strings.Index (after_tag, "<td>")
	if start_loc == -1 {
		log.Fatal ("Missing <td>")
	}
	start_loc = start_loc + 4;
	end_loc = strings.Index (after_tag, "</td>")
	if end_loc == -1 {
		log.Fatal ("Missing </td>")
	}

	correctables = after_tag[start_loc:end_loc]

	/* Chop off the correctable count */
	after_tag = after_tag[end_loc+1:len(after_tag)]
	
	/* Search for the uncorrectable count */
	start_loc = strings.Index (after_tag, "<td>")
	if start_loc == -1 {
		log.Fatal ("Missing <td>")
	}
	start_loc = start_loc + 4;
	end_loc = strings.Index (after_tag, "</td>")
	if end_loc == -1 {
		log.Fatal ("Missing </td>")
	}
	
	/* Extract it */
	uncorrectables = after_tag[start_loc:end_loc]

	now = time.Now()
	year = now.Year()
	month = int(now.Month())
	day = now.Day()
	hour = now.Hour()
	minute = now.Minute()
	second = now.Second()

	fmt.Printf("%02d-%02d-%02d-%02d:%02d:%02d: ",
		year, month, day, hour, minute, second)
	fmt.Printf ("%s correctable errors, %s uncorrectable errors\n", 
		correctables, uncorrectables)
}

func main() {
	var before, after time.Time
	var sleep_time time.Duration
	var work_time time.Duration
	var work_secs, sleep_secs, delay int
	var once_only *bool
	var delay_arg *int

	once_only = flag.Bool ("once", false, "run once and quit")
	delay_arg = flag.Int ("delay", 60, "delay between runs")

	flag.Parse()

	delay = *delay_arg

	for true {
		/* It takes several seconds to scrape the stats... */
		before = time.Now()
		scrape()
		after = time.Now()
		
		if *once_only == true {
			break;
		}

		work_time = after.Sub(before)
		work_secs = int(((work_time + 500000000)/ time.Second))
		sleep_secs = delay - work_secs
		sleep_time = time.Duration(sleep_secs) * time.Second
		time.Sleep (sleep_time)
	}
}

/*
 * Local Variables:
 * eval: (set-variable 'tab-width 4 t)
 * End:
 */
