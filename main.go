package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"time"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
)

func main() {
	// ----- flag parsing -----

	url := flag.String("url", "", "input URL")
	output := flag.String("output", "", "output image path")
	source := flag.String("source", "", "use builtin source and scaling options")
	format := flag.Bool("strftime", false, "enable strftime formatting in URL")
	verbose := flag.Bool("verbose", false, "enable debug output")
	xpath := flag.String("xpath", "", "xpath to <img> tag in url")
	xpath_title := flag.String("xpath-title", "", "xpath to title in url")
	title_font_path := flag.String("title-font-path", "/usr/share/fonts/ttf/noto/NotoSans-Regular.ttf", "path to TTF title font")
	title_font_size := flag.Float64("title-font-size", 30, "title font size")
	title_font_dpi := flag.Float64("title-font-dpi", 72, "title font dpi")
	test := flag.Bool("test", false, "disable wait-online and cooldown")
	mode := flag.String("mode", "fill", "image scaling mode (fill, center)")
	scale := flag.Float64("scale", 1, "scale image prior to centering")
	// top := flag.Int("top", 0, "crop from top")
	// left := flag.Int("left", 0, "crop from left")
	// right := flag.Int("right", 0, "crop from right")
	// bottom := flag.Int("bottom", 0, "crop from bottom")
	cooldown := flag.Int("cooldown", 3600, "minimum seconds to wait before attempting download again")
	flag.Parse()

	if *verbose {
		LOG_LEVEL = "debug"
	}

	var face font.Face
	if *xpath_title != "" {
		face = loadSystemFont(*title_font_path, *title_font_size, *title_font_dpi)
	}

	handle_image := func(img image.Image, title string) {
		img = adjust(img, *mode, *scale)
		if title != "" {
			addLabelByMiddle(img.(draw.Image), img.Bounds().Max.X/2, 50, face, title)
		}
		imaging.Save(img, *output)
		debug("Image saved to ", *output)
	}

	var img image.Image
	var title string
	var err error

	// download/rescale image, then quit
	if *test {
		// use a built-in image source
		if *source != "" {
			img, err = sources[*source]()
		} else {
			img, title, err = custom(*url, *format, *xpath, *xpath_title)
		}

		if err != nil {
			panic(err)
		}
		handle_image(img, title)
	} else {
		// initialize with zero date
		time_last_success := time.Time{};

		online := make(chan int)
		go wait_online(online)

		// loop forever and wait for network online events
		for {
			// wait for network online message from wpa supplicant
			<- online
			debug("Network online")

			// FIXME - need to wait a few seconds for DNS?
			time.Sleep(5 * time.Second)

			// make sure we don't hammer server every time wifi is turned on
			if time.Now().Sub(time_last_success).Seconds() > float64(*cooldown) {

				if *source != "" {
					img, err = sources[*source]()
				} else {
					img, title, err = custom(*url, *format, *xpath, *xpath_title)
				}

				if err == nil {
					time_last_success = time.Now()
				} else {
					fmt.Println(err)
					continue
				}
			} else {
				debug("Hit cooldown limit")
				continue
			}

			handle_image(img, title)
		}
	}
}
