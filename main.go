package main

import (
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"os"
	"time"
)

func main() {
	// ----- flag parsing -----

	url := flag.String("url", "", "input URL")
	output := flag.String("output", "", "output image path")
	source := flag.String("source", "", "use builtin source and scaling options")
	format := flag.Bool("strftime", false, "enable strftime formatting in URL")
	verbose := flag.Bool("verbose", false, "enable debug output")
	xpath := flag.String("xpath", "", "xpath to <img> tag in url")
	test := flag.Bool("test", false, "disable wait-online and cooldown")
	mode := flag.String("mode", "fill", "image scaling mode (fill, center)")
	scale := flag.Float64("scale", 1, "scale image prior to centering")
	cooldown := *flag.Int64("cooldown", 3600, "minimum seconds to wait before attempting download again")
	flag.Parse()

	if *verbose {
		LOG_LEVEL = "debug"
	}

	var img image.Image
	var err error

	// download/rescale image, then quit
	if *test {
		// use a built-in image source
		if *source != "" {
			img, err = sources[*source]()
		} else {
			img, err = custom(*url, *format, *xpath)
		}

		if err != nil {
			panic(err)
		}
		// img = adjust(img, *top, *left, *right, *bottom)
		img = adjust(img, *mode, *scale)
		imaging.Save(img, *output)
		debug("Image saved to ", *output)
		return
	}

	var time_last_success time.Time

	if stat, err := os.Stat(*output); err == nil {
		time_last_success = stat.ModTime()
	} else {
		time_last_success = time.Time{}
	}

	// loop forever and wait for network online events
	for {
		if time.Now().Sub(time_last_success).Seconds() < float64(cooldown) {
			debug("Hit cooldown limit")
			time.Sleep(time.Duration(cooldown) * time.Second)
			continue
		}

		if *source != "" {
			img, err = sources[*source]()
		} else {
			img, err = custom(*url, *format, *xpath)
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		img = adjust(img, *mode, *scale)
		imaging.Save(img, *output)
		debug("Image saved to ", *output)
		time_last_success = time.Now()

		wait_for := time.Duration(cooldown) * time.Second
		debug(fmt.Sprintf("Sleeping for %v", wait_for))
		time.Sleep(wait_for)
	}
}
