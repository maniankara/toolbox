#!/bin/bash

if [[ "$1" -eq "" ]];then
  echo "Usage: "
  echo "    ./`basename $0` [on|off]"
  echo "    ./`basename $0` on"
  echo "    ./`basename $0` off"
  exit -1
elif [ $1 -eq "on" ]; then
  echo "Turning vpn on, might need sudo password"
  sudo cp /run/resolvconf/resolv.conf-vpnon /etc/resolv.conf
elif [ $1 -eq "off" ]; then
  echo "Turning vpn off, might need sudo password"
  sudo cp /run/resolvconf/resolv.conf-vpnoff /etc/resolv.conf
fi

