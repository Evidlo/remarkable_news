.ONESHELL:
# .SILENT:
SHELL:=/bin/bash

host=10.11.99.1
timezone=America/Chicago
cooldown=3600

renews.arm:
	go get ./...
	env GOOS=linux GOARCH=arm GOARM=7 go build -o renews.arm

renews.x86:
	go get ./...
	go build -o renews.x86

# get latest prebuilt releases
.PHONY: download_prebuilt
download_prebuilt:
	wget http://github.com/evidlo/remarkable_news/releases/latest/download/release.zip
	unzip release.zip

# build release
.PHONY: release
release: renews.arm renews.x86
	zip release.zip renews.arm renews.x86

clean:
	rm -f renews.x86 renews.arm release.zip

define install
	eval "$(ssh-agent -s)"
	eval $(shell ssh-agent -s)
	ssh-add
	# stop running service, ignore failure to stop
	ssh root@$(host) systemctl stop renews || true
	scp renews.arm root@$(host):
	# substitute timezone/cooldown/KEYWORDS arguments
	sed -e "s|TZ|$(timezone)|" \
		-e "s|COOLDOWN|$(cooldown)|" \
		-e "s|KEYWORDS|$(KEYWORDS)|" \
		$(1) > renews.service
	# back up suspend screen.  don't overwrite existing file
	ssh root@$(host) "cd /usr/share/remarkable/; ls suspended_back.png 2> /dev/null || cp suspended.png suspended_back.png"
	# copy service to remarkable and enable
	scp renews.service root@$(host):/etc/systemd/system/renews.service
	ssh root@$(host) <<- ENDSSH
		systemctl daemon-reload
		systemctl enable renews
		systemctl restart renews
	ENDSSH
endef

# ----- Sources -----

.PHONY: install_nyt
install_nyt: renews.arm
	$(call install,services/nyt.service)

.PHONY: install_nyt_hq
install_nyt_hq: renews.arm
	$(call install,services/nyt-hq.service)

.PHONY: install_xkcd
install_xkcd: renews.arm
	$(call install,services/xkcd.service)

.PHONY: install_wp
install_wp: renews.arm
	$(call install,services/wp.service)

.PHONY: install_picsum
install_picsum: renews.arm
	$(call install,services/picsum.service)

.PHONY: install_loremflickr
install_loremflickr: renews.arm
	$(call install,services/loremflickr.service)

.PHONY: install_cah
install_cah: renews.arm
	$(call install,services/cah.service)

# .PHONY: install_wikipotd
# install_wikipotd: renews.arm
# 	$(call install,services/wikipotd.service)
