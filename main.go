package main

import (
	"fmt"
	"net/http"
	"flag"
	"time"
	"io/ioutil"
	"image"
	"github.com/lestrrat-go/strftime"
	"github.com/pkg/errors"
	"github.com/disintegration/imaging"
)

var LOG_LEVEL = "error"

func check(err error, msg string) {
	if err != nil {
		panic(errors.Wrap(err, msg));
	}
}

func debug(msg ...string) {
	if LOG_LEVEL == "debug" {
		fmt.Println(msg)
	}
}


func main() {
	// ----- flag parsing -----

	url := flag.String("url", "", "input URL")
	output := flag.String("output", "", "output image path")
	format := flag.Bool("strftime", false, "enable strftime formatting in URL")
	verbose := flag.Bool("verbose", false, "enable debug output")
	timezone := flag.String("timezone", "", "override timezone (tzinfo format)")
	top := flag.Int("top", 0, "crop from top")
	left := flag.Int("left", 0, "crop from left")
	right := flag.Int("right", 0, "crop from right")
	bottom := flag.Int("bottom", 0, "crop from bottom")
	cooldown := flag.Int("cooldown", 3600, "minimum seconds to wait before attempting download again")
	flag.Parse()

	if *verbose {
		LOG_LEVEL = "debug"
	}

	time_last_success := time.Time{}


	online := make(chan int)
	go wait_online(online)

	for {
		// ----- network wait online -----

		// wait for network online message from wpa supplicant
		<- online
		debug("Network online")

		// FIXME - need to wait a few seconds for DNS?
		time.Sleep(5 * time.Second)

		// allow strftime formatting for date-dependent urls
		if *format {
			if *timezone != "" {
				tz, err := time.LoadLocation(*timezone)
				check(err, "")
				*url, err = strftime.Format(*url, time.Now().In(tz))
				check(err, "")
			} else {
				var err error
				*url, err = strftime.Format(*url, time.Now())
				check(err, "")
			}
			debug("strftime formatted URL:", *url)
		}

		// make sure we don't hammer server every time wifi is turned on
		if time.Now().Sub(time_last_success).Seconds() > float64(*cooldown) {
			debug("Beginning download")

			// ----- image download -----

			// if http failure, wait for next reconnect
			response, err := http.Get(*url)
			if err != nil {
				fmt.Println(err)
				continue

			}
			defer response.Body.Close()
			// if http error code, wait for next reconnect
			if response.StatusCode != 200 {
				body, _ := ioutil.ReadAll(response.Body)
				fmt.Println(body)
				continue
			}
			time_last_success = time.Now()
			img, err := imaging.Decode(response.Body)
			check(err, "")

			// ----- image resizing/cropping -----

			rect := img.Bounds()
			img = imaging.Crop(
				img,
				image.Rect(
					rect.Min.X + *left,
					rect.Min.Y + *top,
					rect.Max.X - *right,
					rect.Max.Y - *bottom,
				),
			)
			// fit image
			img = imaging.Fill(img, 1404, 1872, imaging.Top, imaging.NearestNeighbor)
			// img = imaging.Fill(img, 1404, 1872, imaging.Top, imaging.Linear)
			// img = imaging.Fill(img, 1404, 1872, imaging.Top, imaging.Box)
			// img = imaging.Fill(img, 1404, 1872, imaging.Top, imaging.Lanczos)

			imaging.Save(img, *output)
			debug("Image saved to ", *output)
		} else {
			debug("Hit cooldown limit")
		}
	}
}
