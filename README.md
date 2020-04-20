# remarkable_news

Automatically download daily newspaper/comic as your suspend screen.  No cloud needed.

![demo](pic.png)

## Install

Assuming you have Go installed

    make install_nyt
    
This will install and start the newspaper fetch service on the reMarkable.  Every time you connect to WiFi, it will try to grab the latest front page from The New York Times.
    
Alternatively you can use the prebuilt release if you don't have Go

    make download_prebuilt
    make install_nyt
    
## Supported News/Comics Sources

- XKCD - `make install_xkcd`
- Washington Post - `make install_wp`
- New York Times - `make install_nyt`
    
## Contributing

I'm looking for help adding more comics/news sources.  Currently remarkable_news supports `.jpg`, `.png`, `.tiff`, and `.bmp` sources.  New source URLs can be added to the [Makefile](Makefile) in the *Sources* section.

URLs can be date dependent.  For example, this URL template

    http://example.com/%%Y/%%m/%%d/news.jpg
    
would be converted to this URL by remarkable_news.

    http://example.com/2020/04/18/news.jpg
    
The full list of date formatting options are listed [here](https://github.com/lestrrat-go/strftime#supported-conversion-specifications).  Two percent signs should be used instead of just one, as in the example above.

## Debugging

    journalctl --unit renews -f
    
Then disconnect and reconnect WiFi to trigger a download.  remarkable_news will only download at a maximum of once per hour to avoid burdening the server.
