# remarkable_news

reMarkable service to automatically download daily newspaper/comic as your suspend screen.  No cloud needed.

![demo](pic.png)

## Install (Linux/OSX)

Assuming you have Go installed

    git clone http://github.com/evidlo/remarkable_news && cd remarkable_news
    make install_nyt
    
This will install and start the newspaper fetch service on the reMarkable.  Every time you connect to WiFi, it will try to grab the latest front page from The New York Times.
    
Alternatively you can use the prebuilt release if you don't want to install Go

    git clone http://github.com/evidlo/remarkable_news && cd remarkable_news
    make download_prebuilt
    make install_nyt
    
## Install (Windows)

Install [WSL](https://docs.microsoft.com/en-us/learn/modules/get-started-with-windows-subsystem-for-linux/2-enable-and-install), then follow the Linux/OSX instructions.
    
## Supported News/Comics Sources

- XKCD - `make install_xkcd`
- Washington Post (only updates weekdays) - `make install_wp`
- New York Times (slightly low resolution) - `make install_nyt`
- Wikipedia Picture of the Day - `make install_wikipotd`
    
## Contributing

#### Additional Sources

I'm looking for help adding more comics/news sources.  Currently remarkable_news supports `.jpg`, `.png`, `.tiff`, and `.bmp` sources.  New source configurations can be added in the `services/` folder and a new target should be added to the [Makefile](Makefile) in the *Sources* section.

The easiest news/comics sources to add are those that have a static link to the latest image.  However, this is often not the case, so remarkable_news can handle these situations in two ways:

- date dependent URLs - See [this file](/services/nyt.service) for an example.
- <img> tag parsing from html (via [xpath expressions](https://www.webperformance.com/load-testing-tools/blog/articles/real-browser-manual/building-a-testcase/how-locate-element-the-page/xpath-locator-examples/)) - See [this file](/services/xkcd.service) for an example.

#### New features

Also, there are some additional features I would like to get added

- more options for scaling, margins
- parse image titles from html (also via xpaths)
- parse image description from html (via xpaths), would be great for Wikipedia picture of the day to have a caption

#### Testing on host machine

Run `renews.x86` with the `-test` option.  This disables download cooldown and waiting for WiFi connect.

## Debugging

    journalctl --unit renews -f
    
Then disconnect and reconnect WiFi to trigger a download.  remarkable_news will only download at a maximum of once per hour to avoid burdening the server.
