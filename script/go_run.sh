#!/bin/sh

go run ../main.go \
    --etcd-endpoints="http://10.10.10.219:2379,http://10.10.10.220:2379,http://10.10.10.221:2379" \
    --data-center="testing-dc" \
    --cert-path="../../tsa/data/certs" \
    --ntp-bind="0.0.0.0:1123" \
    --bind-eth="wlp3s0" \
    --tac-servername="s1.ta.ntsc.ac.cn" \
    --tac-addr="tcp://10.10.10.250:1358" \
    --es-endpoints="http://10.10.10.179:9200,http://10.10.10.180:9200,http://10.10.10.181:9200" \
    --fix-time=false
