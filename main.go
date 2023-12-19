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
	title_str := flag.String("title", "", "title to use instead of xpath-title")
	title_font_path := flag.String("title-font-path", "/usr/share/fonts/ttf/noto/NotoSans-Bold.ttf", "path to TTF title font")
	title_font_size := flag.Float64("title-font-size", 50, "title font size")
	xpath_subtitle := flag.String("xpath-subtitle", "", "xpath to subtitle in url")
	subtitle_str := flag.String("subtitle", "", "subtitle to use instead of xpath-subtitle")
	subtitle_font_path := flag.String("subtitle-font-path", "/usr/share/fonts/ttf/noto/NotoSans-Regular.ttf", "path to TTF subtitle font")
	subtitle_font_size := flag.Float64("subtitle-font-size", 30, "subtitle font size")
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

	var title_face font.Face
	if *xpath_title != "" || *title_str != "" {
		title_face = loadSystemFont(*title_font_path, *title_font_size)
	}
	var subtitle_face font.Face
	if *xpath_subtitle != "" || *subtitle_str != "" {
		subtitle_face = loadSystemFont(*subtitle_font_path, *subtitle_font_size)
	}

	handle_image := func(img image.Image, title, subtitle string) {
		img = adjust(img, *mode, *scale)
		if title == "" {
			title = *title_str
		}
		if title != "" {
			addCenteredLabel(img.(draw.Image), 100, title_face, title)
		}
		if subtitle == "" {
			subtitle = *subtitle_str
		}
		if subtitle != "" {
			addCenteredLabel(img.(draw.Image), 150, subtitle_face, subtitle)
		}
		imaging.Save(img, *output)
		debug("Image saved to ", *output)
	}

	var img image.Image
	var title string
	var subtitle string
	var err error

	// download/rescale image, then quit
	if *test {
		// use a built-in image source
		if *source != "" {
			img, err = sources[*source]()
		} else {
			img, title, subtitle, err = custom(*url, *format, *xpath, *xpath_title, *xpath_subtitle)
		}

		if err != nil {
			panic(err)
		}
		handle_image(img, title, subtitle)
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
					img, title, subtitle, err = custom(*url, *format, *xpath, *xpath_title, *xpath_subtitle)
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

			handle_image(img, title, subtitle)
		}
	}
}
