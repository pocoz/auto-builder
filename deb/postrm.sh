#!/bin/sh -xe

echo "post rm"

if [[ -f "/etc/init.d/auto-builder.sh" ]]; then
  invoke-rc.d auto-builder stop
  update-rc.d auto-builder disable
fi
if [[ -f "/lib/systemd/system/auto-builder.service" ]]; then
  systemctl stop auto-builder.service
  systemctl disable auto-builder.service
fi
