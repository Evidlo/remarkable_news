## Contributing

#### Additional Sources

I'm looking for help adding more comics/news sources.  Currently remarkable_news supports `.jpg`, `.png`, `.tiff`, and `.bmp` sources.  New source configurations can be added in the `services/`.  See [this file](/services/picsum.service) for a simple example with a fixed link to a `.jpg`.

The easiest news/comics sources to add are those that have a static link to the latest image.  However, this is often not the case, so remarkable_news can handle these situations in two ways:

- date dependent URLs - See [this file](/services/nyt.service) for an example.
- <img> tag parsing from html (via [xpath expressions](https://www.webperformance.com/load-testing-tools/blog/articles/real-browser-manual/building-a-testcase/how-locate-element-the-page/xpath-locator-examples/)) - See [this file](/services/xkcd.service) for an example.

#### Testing on host machine

[Download](https://github.com/evidlo/remarkable_news/releases/latest/download/release.zip) and run `renews.x86` with the `-test` option.  This disables download cooldown and waiting for WiFi connect.

Here is an example command which I used for testing while creating the Calvin and Hobbes source:

    ./renews.x86 -output test.png -verbose -url https://www.gocomics.com/random/calvinandhobbes -xpath '//picture[@class="item-comic-image"]/img/@src' -mode fill -scale 0.9 -test
    
This outputs to `test.png`.

#### Usage

    [evan@blackbox remarkable_news] ./renews.x86 -h
    Usage of ./renews.x86:
      -cooldown int
            minimum seconds to wait before attempting download again (default 3600)
      -mode string
            image scaling mode (fill, center) (default "fill")
      -output string
            output image path
      -scale float
            scale image prior to centering (default 1)
      -source string
            use builtin source and scaling options
      -strftime
            enable strftime formatting in URL
      -test
            disable wait-online and cooldown
      -url string
            input URL
      -verbose
            enable debug output
      -xpath string
            xpath to <img> tag in url

#### New features

Also, there are some additional features I would like to get added

- more options for scaling, margins
- parse image titles from html (also via xpaths)
- parse image description from html (via xpaths), would be great for Wikipedia picture of the day to have a caption
