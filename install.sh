# systemd service location
SERVICE=/etc/systemd/system/renews.service

# stop service if running
systemctl stop renews.service || true
# install the renews binary from github releases
mkdir -p /home/root/bin
cd /home/root/bin
wget -O release.zip http://github.com/evidlo/remarkable_news/releases/latest/download/release.zip
unzip -o release.zip

# install systemd service
# mv renews.service ${SERVICE}
wget -O ${SERVICE} "https://evidlo.github.io/remarkable_news/services/${1}.service"

# substitute COOLDOWN and KEYWORDS arguments
if [[ -z $COOLDOWN ]]
then
    COOLDOWN=3600
fi
sed -i "s|COOLDOWN|${COOLDOWN}|" ${SERVICE}
sed -i "s|KEYWORDS|${KEYWORDS}|" ${SERVICE}

# reload systemd and remove extra files
systemctl daemon-reload
systemctl enable --now renews.service
rm renews.x86 release.zip
