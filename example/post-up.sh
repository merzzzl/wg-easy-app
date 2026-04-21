#!/usr/bin/env bash
set -e

ip route add 149.154.166.110/32 dev wg1
ip route add default dev wg1 table 100
ip rule add iif wg0 table 100 priority 100

iptables -I FORWARD 1 -i wg0 -o wg1 -j ACCEPT
iptables -I FORWARD 1 -i wg1 -o wg0 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
iptables -t nat -I POSTROUTING 1 -s 10.8.10.0/24 -o wg1 -j MASQUERADE