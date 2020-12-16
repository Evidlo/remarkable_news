package main

import (
	"image"
	"image/color"
	"net/http"
	"fmt"
	"strings"

	"github.com/disintegration/imaging"

	// FIXME - resizing already built into imaging, but this is much faster
	"github.com/nfnt/resize"

	"io/ioutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	// "golang.org/x/image/font/inconsolata"
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
func custom(url string, format bool, xpath string) (image.Image, error) {

	debug("Beginning download")

	// ----- URL strftime formatting -----

	if format {
		url = format_url(url)
	}

	// ----- image XPath handling -----

	var response *http.Response
	var err error
	// if xpath is provided, assume url is HTML
	if xpath != "" {
		debug("Got -xpath.  Trying to extract img url from provided url")

		result, err := get_xpath(url, xpath, "html")
		check(err, "")

		// imgurl := e.Attr[0].Val
		imgurl, err := to_absurl(url, result)

		debug("Image url", imgurl)

		// if http failure, wait for next reconnect
		response, err = get_url(imgurl)
		if err != nil {
			debug("Failed to fetch image")
			return nil, err
		}

	} else {
		response, err = get_url(url)
		check(err, "")
	}

	// ----- image loading -----

	img, err := imaging.Decode(response.Body)
	if err != nil {
		debug("Failed to decode image")
		return nil, err
	}

	return img, nil

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

func addText(img image.Image, y int, label string) image.Image {

	otfData, err := ioutil.ReadFile("./xkcd-Regular.otf")
	check(err, "Couldn't load OTF font")

	otf, err := opentype.Parse(otfData)
	check(err, "Couldn't parse font data")

	xkcdFont, err := opentype.NewFace(otf, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		// Hinting: font.HintingFull,
	})
	check(err, "Couldn't create font face")

	textColor := color.RGBA{50, 50, 50, 255}
	upperCaseLabel := strings.ToUpper(label)

	d := &font.Drawer{
		Dst:  img.(*image.NRGBA),
		Src:  image.NewUniform(textColor),
		Face: xkcdFont,
	}

	width := d.MeasureString(upperCaseLabel)
	x := (fixed.I(1404) - width) / 2.0

	d.Dot = fixed.Point26_6{
		X: x,
		Y: fixed.I(y),
	}
	d.DrawString(upperCaseLabel)

	return img
}
