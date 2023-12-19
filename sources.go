package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	// FIXME - resizing already built into imaging, but this is much faster
	"github.com/nfnt/resize"
)

var sources = map[string] func() (image.Image, error) {
	"natgeo": natgeo,
}

func natgeo() (image.Image, error){
	url := "https://www.nationalgeographic.com/photography/photo-of-the-day/_jcr_content/.gallery.json"

	imgurl, err := get_xpath(url, "/items/*[1]/image/uri", "json")
	check(err, "")

	caption, err := get_xpath(url, "/items/*[1]/image/caption", "json")
	caption = strings.TrimSuffix(strings.TrimPrefix(caption, "<p>"), "</p>\n")
	fmt.Println(caption)
	check(err, "")

	// if http failure, wait for next reconnect
	response, err := get_url(imgurl)
	if err != nil {
		debug("Failed to fetch image")
		return nil, err
	}

	img, err := imaging.Decode(response.Body)
	if err != nil {
		debug("Failed to decode image")
		return nil, err
	}

	return img, nil
}

// function for grabbing custom sources
func custom(url string, format bool, xpath, xpath_title, xpath_subtitle string) (image.Image, string, string, error) {

	debug("Beginning download")

	// ----- URL strftime formatting -----

	if format {
		url = format_url(url)
	}

	// ----- image XPath handling -----

	var title string
	var subtitle string
	var response *http.Response
	var err error

	// if xpath is provided, assume url is HTML
	if xpath != "" {
		if xpath_title != "" {
			debug("Got -xpath-title. Trying to extract title from provided url")
			title, err = get_xpath(url, xpath_title, "html")
			check(err, "")
		}
		if xpath_subtitle != "" {
			debug("Got -xpath-subtitle. Trying to extract subtitle from provided url")
			subtitle, err = get_xpath(url, xpath_subtitle, "html")
			check(err, "")
		}

		debug("Got -xpath. Trying to extract img url from provided url")

		result, err := get_xpath(url, xpath, "html")
		check(err, "")

		// imgurl := e.Attr[0].Val
		imgurl, err := to_absurl(url, result)

		debug("Image url", imgurl)

		// if http failure, wait for next reconnect
		response, err = get_url(imgurl)
		if err != nil {
			debug("Failed to fetch image")
			return nil, "", "", err
		}

	} else {
		response, err = get_url(url)
		check(err, "")
	}

	// ----- image loading -----

	img, err := imaging.Decode(response.Body)
	if err != nil {
		debug("Failed to decode image")
		return nil, "", "", err
	}

	return img, title, subtitle, nil

}

// scale, inset image to reMarkable display size
func adjust(img image.Image, mode string, scale float64) image.Image {

	debug("Adjusting image")

	re_width := 1404
	re_height := 1872

	if mode == "fill" {
		// scale image to remarkable width
		// imaging resize is slow for some reason, use other library
		// img = imaging.Resize(img, re_width, 0, imaging.Linear)
		img = resize.Resize(uint(re_width), 0, img, resize.Bilinear)
		// cut off parts of image that overflow
		img = imaging.Crop(img, image.Rect(0, 0, re_width, re_height))

	} else if mode == "center" {
	} else {
		debug("Invalid mode")
	}
	if scale != 1 {
		img_width := float64(img.Bounds().Max.X)
		img = resize.Resize(uint(scale * img_width), 0, img, resize.Bilinear)
	}

	// put image in center of screen
	background := imaging.New(
		re_width,
		re_height,
		color.RGBA{255, 255, 255, 255},
	)
	img = imaging.PasteCenter(background, img)

	return img
}

func loadSystemFont(path string, size float64) font.Face {
	fontdata, err := os.ReadFile(path)
	check(err, "Failed to open font file")
	font, err := truetype.Parse(fontdata)
	check(err, "Failed to parse font")

	return truetype.NewFace(font, &truetype.Options{
		Size: size,
	})
}

func addCenteredLabel(img draw.Image, y int, face font.Face, label string) {
	x := img.Bounds().Max.X / 2
	b, _ := font.BoundString(face, label)
	label_width := (b.Max.X - b.Min.X).Ceil()

	// Hack: Wrap labels that are too long once
	if label_width > img.Bounds().Max.X {
		debug("Label too long, splitting into two lines")
		left := strings.TrimSpace(label[:len(label)/2])
		right := strings.TrimSpace(label[len(label)/2:])
		right_parts := strings.SplitN(right, " ", 2)
		first_line := left + right_parts[0]

		b, _ = font.BoundString(face, first_line)
		label_width = (b.Max.X - b.Min.X).Ceil()
		label_height := (b.Max.Y - b.Min.Y).Ceil()
		addLabel(img, x-label_width/2, y, face, first_line)

		if len(right_parts) == 2 {
			second_line := right_parts[1]
			b, _ = font.BoundString(face, second_line)
			label_width = (b.Max.X - b.Min.X).Ceil()
			addLabel(img, x-label_width/2, y+label_height+10, face, second_line)
		}
	} else {
		addLabel(img, x-label_width/2, y, face, label)
	}

}

func addLabel(img draw.Image, x, y int, face font.Face, label string) {
	debug("Adding label", label)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	d.DrawString(label)
}
