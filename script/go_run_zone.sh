#!/bin/sh

go run ../main.go \
    --cert-path="../../data/certs" \
    --ntp-bind="0.0.0.0:1123" \
    --bind-eth="eth0" \
    --tac-servername="s1.ta.ntsc.ac.cn" \
    --tac-addr="tcp://10.10.10.250:1358" \
    --cql-endpoints="10.10.10.218" \
    --cql-ks="" \
    --cql-pwd="1234qwer"
