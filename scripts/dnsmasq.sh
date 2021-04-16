#!/bin/bash

sudo firewall-cmd --add-port={67/tcp,69/udp,69/tcp,4011/udp}
sudo dnsmasq -d -C config/dnsmasq.conf
