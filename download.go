package main

import (
	"time"
	"net/http"
	// "fmt"
	"net/url"
	"image"
	"errors"
	"github.com/lestrrat-go/strftime"
	"github.com/disintegration/imaging"
	"github.com/antchfx/htmlquery"
)

var Err404 = errors.New("Err404")

func to_absurl(base, rel string) (string, error) {
	base_url, err := url.Parse(base)
	if err != nil {
		debug("Error parsing URL:", base)
		return "", err
	}
	rel_url, err := url.Parse(rel)
	if err != nil {
		debug("Error parsing URL:", rel)
		return "", err
	}

	absurl := base_url.ResolveReference(rel_url)

	return absurl.String(), nil
}

func get_url(url string) (*http.Response, error){
	// if http failure, wait for next reconnect
	response, err := http.Get(url)
	if err != nil {
		return response, err
	}
	// defer response.Body.Close()

	// if http error code, wait for next reconnect
	if response.StatusCode != 200 {
		// body, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(body)
		return response, Err404
	}

	return response, nil
}

func download(url string, format bool, timezone, xpath string) (image.Image, error){

	debug("Beginning download")

	// ----- URL strftime formatting -----

	if format {
		// override default timezone
		if timezone != "" {
			tz, err := time.LoadLocation(timezone)
			check(err, "")
			url, err = strftime.Format(url, time.Now().In(tz))
			check(err, "")
		} else {
			var err error
			url, err = strftime.Format(url, time.Now())
			check(err, "")
		}
		debug("strftime formatted URL:", url)
	}

	// ----- fetch html/image -----

	response, err := get_url(url)
	if err != nil {
		debug("Failed to fetch url")
		return nil, err
	}

	// ----- image XPath handling -----

	var imgresponse *http.Response
	// if xpath is provided, assume url is HTML
	if xpath != "" {
		debug("Got -xpath.  Trying to extract img url from provided url")

		doc, err := htmlquery.Parse(response.Body)
		if err != nil {
			debug("Failed to parse html")
			panic(err)
		}

		e := htmlquery.FindOne(doc, xpath)

		// imgurl := e.Attr[0].Val
		imgurl, err := to_absurl(url, e.Attr[0].Val)

		debug("Image url", imgurl)

		// if http failure, wait for next reconnect
		imgresponse, err = get_url(imgurl)
		if err != nil {
			debug("Failed to fetch image")
			return nil, err
		}

	} else {
		imgresponse = response
	}

	// ----- image resizing/cropping -----

	// fmt.Println(imgresponse.StatusCode)
	// body, _ := ioutil.ReadAll(response.Body)
	// ioutil.WriteFile("/tmp/dump", body, 0755)
	img, err := imaging.Decode(imgresponse.Body)
	if err != nil {
		debug("Failed to decode image")
		return nil, err
	}

	return img, nil

}

func adjust(img image.Image, top, left, right, bottom int) image.Image {
	rect := img.Bounds()
	img = imaging.Crop(
		img,
		image.Rect(
			rect.Min.X + left,
			rect.Min.Y + top,
			rect.Max.X - right,
			rect.Max.Y - bottom,
		),
	)
	// fit image
	img = imaging.Fill(img, 1404, 1872, imaging.Top, imaging.NearestNeighbor)
	// img = imaging.Fill(img, 1404, 1872, imaging.Top, imaging.Linear)
	// img = imaging.Fill(img, 1404, 1872, imaging.Top, imaging.Box)
	// img = imaging.Fill(img, 1404, 1872, imaging.Top, imaging.Lanczos)

	return img
}
