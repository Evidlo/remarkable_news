.ONESHELL:
host=10.11.99.1
timezone=America/Chicago
cooldown=3600

renews.arm:
	go get ./...
	env GOOS=linux GOARCH=arm GOARM=5 go build -o renews.arm

renews.x86:
	go get ./...
	go build -o renews.x86

# get latest prebuilt releases
.PHONY: download_prebuilt
download_prebuilt:
	wget https://github.com/Evidlo/remarkable_news/releases/download/1/release.zip
	unzip release.zip

# build release
.PHONY: release
release: renews.arm renews.x86
	zip release.zip renews.arm renews.x86

clean:
	rm -f renews.x86 renews.arm release.zip

define install
	ssh-add
	ssh root@$(host) systemctl stop renews
	scp renews.arm root@$(host):
	# substitute timezone/cooldown arguments
	sed -e "s|TZ|$(timezone)|" \
		-e "s|COOLDOWN|$(cooldown)|" \
		$(1) > renews.service
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

.PHONY: install_xkcd
install_xkcd: renews.arm
	$(call install,services/xkcd.service)

.PHONY: install_wp
install_wp: renews.arm
	$(call install,services/wp.service)

# .PHONY: install_wikipotd
# install_wikipotd: renews.arm
# 	$(call install,services/wikipotd.service)
