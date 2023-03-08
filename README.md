# remarkable_news

reMarkable service to automatically download daily newspaper/comic as your suspend screen.  No cloud needed.

![demo](pic.png)


## Install

    wget -O - http://evidlo.github.io/remarkable_news/install.sh | sh /dev/stdin nyt

This will install and start the update service on the reMarkable.  Every time you connect to WiFi, it will try to grab the latest front page from The New York Times.  See below for more image sources.

By default, downloads are rate limited to once per hour (3600 s).  This can be overriden by modifying `/etc/systemd/system/renews.service`

Requires [remarkable-hacks](https://github.com/ddvk/remarkable-hacks) to be installed for software versions >=2.5.0.27

## Install (Windows)

Install [WSL](https://docs.microsoft.com/en-us/learn/modules/get-started-with-windows-subsystem-for-linux/2-enable-and-install), then follow the Linux/OSX instructions.  This has not been tested.

## Supported News/Comics Sources

- XKCD 
    - `wget -O - http://evidlo.github.io/remarkable_news/install.sh | sh /dev/stdin xkcd`
- Washington Post (only updates weekdays) 
    - `wget -O - http://evidlo.github.io/remarkable_news/install.sh | sh /dev/stdin wp`
- New York Times (slightly low resolution) 
    - `wget -O - http://evidlo.github.io/remarkable_news/install.sh | sh /dev/stdin nyt`
- New York Times (high quality provided by [acrogenesis/nyt-today](https://github.com/acrogenesis/nyt-today)) 
    - `wget -O - http://evidlo.github.io/remarkable_news/install.sh | sh /dev/stdin xkcd`
- Picsum (random images) 
    - `wget -O - http://evidlo.github.io/remarkable_news/install.sh | sh /dev/stdin picsum`
- LoremFlickr (random images) 
    - `wget -O - http://evidlo.github.io/remarkable_news/install.sh | KEYWORDS=nature,cats sh /dev/stdin loremflicker`
    - replace 'nature,cats' with your own keywords
- Unsplash (random images)
    - `wget -O - http://evidlo.github.io/remarkable_news/install.sh | KEYWORDS=nature sh /dev/stdin unsplash`
    - replace 'nature' with your own keywords.  Only one keyword supported
- Calvin and Hobbes 
    - `wget -O - http://evidlo.github.io/remarkable_news/install.sh | sh /dev/stdin cah`
- The Guardian 
    - `wget -O - http://evidlo.github.io/remarkable_news/install.sh | sh /dev/stdin uk_tg`
<!-- - Wikipedia Picture of the Day - `make install_wikipotd` -->


## Debugging

On the reMarkable

    journalctl --unit renews -f

Then disconnect and reconnect WiFi to trigger a download.  `remarkable_news` will only download at a maximum of once per hour to avoid burdening the server.

## Contributing

See [contributing.md](contributing.md)
