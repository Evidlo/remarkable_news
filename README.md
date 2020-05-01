# remarkable_news

reMarkable service to automatically download daily newspaper/comic as your suspend screen.  No cloud needed.

![demo](pic.png)

## Install (Linux/OSX)

Assuming you have Go installed and with the reMarkable connected via USB

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
- Picsum (random images) - `make install_picsum`
- Calvin and Hobbes - `make install_cah`
<!-- - Wikipedia Picture of the Day - `make install_wikipotd` -->
    

## Debugging

On the reMarkable

    journalctl --unit renews -f
    
Then disconnect and reconnect WiFi to trigger a download.  remarkable_news will only download at a maximum of once per hour to avoid burdening the server.
