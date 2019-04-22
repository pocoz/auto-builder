#!/bin/sh -xe

echo "post inst"

addgroup --system auto-builder
adduser --system auto-builder --no-create-home --home /nonexistent

if [[ -f "/etc/init.d/auto-builder.sh" ]]; then
  update-rc.d auto-builder defaults
  invoke-rc.d auto-builder start
fi
if [[ -f "/lib/systemd/system/auto-builder.service" ]]; then
  systemctl daemon-reload
  systemctl enable auto-builder.service
  systemctl start auto-builder.service
fi
