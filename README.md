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
    
## Contributing

I'm looking for help adding more comics/news sources.  Currently remarkable_news supports `.jpg`, `.png`, `.tiff`, and `.bmp` sources.  New source configurations can be added in the `services/` folder and a new target should be added to the [Makefile](Makefile) in the *Sources* section.

URLs can be date dependent.  See [this file](/services/nyt.service) for an example.
    
The full list of date formatting options are listed [here](https://github.com/lestrrat-go/strftime#supported-conversion-specifications).  Two percent signs should be used instead of just one, as in the example above.

remarkable_news also supports parsing images out of webpages using [XPath expression](https://www.webperformance.com/load-testing-tools/blog/articles/real-browser-manual/building-a-testcase/how-locate-element-the-page/xpath-locator-examples/).  See [this file](/services/xkcd.service) for an example.

## Debugging

    journalctl --unit renews -f
    
Then disconnect and reconnect WiFi to trigger a download.  remarkable_news will only download at a maximum of once per hour to avoid burdening the server.

