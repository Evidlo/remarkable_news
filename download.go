package main

import (
	"time"
	"net/http"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/jsonquery"
	// "fmt"
	"net/url"
	"errors"
	"github.com/lestrrat-go/strftime"
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
		debug("Failed to fetch url")
		return response, err
	}
	// defer response.Body.Close()

	// if http error code, wait for next reconnect
	if response.StatusCode != 200 {
		// body, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(body)
		debug("URL 404")
		return response, Err404
	}

	return response, nil
}


func xpath_html(url, xpath string) (string, error) {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		debug("Failed to parse HTML")
		return "", err
	}

	list, err := htmlquery.QueryAll(doc, "//meta/text()")
	check(err, "Invalid XPath")

	if len(list) == 0 {
		debug("No XPath matches found")
		return "", err
	}

	return htmlquery.InnerText(list[0]), nil
}


func get_xpath(url, xpath, data_format string) (string, error) {
	// load the given URL and query the document with the given XPath expression
	// returns string result

	if data_format == "json" {
		doc, err := jsonquery.LoadURL(url)
		if err != nil {
			debug("Failed to parse JSON")
			return "", err
		}

		list, err := jsonquery.QueryAll(doc, xpath)
		check(err, "Invalid XPath")

		if len(list) == 0 {
			debug("No XPath matches found")
			// FIXME, err is nil
			return "", err
		}

		return list[0].InnerText(), nil
	} else if data_format == "html" {
		doc, err := htmlquery.LoadURL(url)
		if err != nil {
			debug("Failed to parse HTML")
			return "", err
		}

		list, err := htmlquery.QueryAll(doc, xpath)
		check(err, "Invalid XPath")

		if len(list) == 0 {
			debug("No XPath matches found")
			// FIXME, err is nil
			return "", err
		}

		return htmlquery.InnerText(list[0]), nil
	}

	panic(`Invalid data_format`)
}

func format_url(url, timezone string) string {
	// format url containing strftime-style datecodes

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

	return url
}
