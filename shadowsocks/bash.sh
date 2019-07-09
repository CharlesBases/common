#!/bin/zsh
server=127.0.0.1
server_port=666666
client_port=888888
method=aes-128-cfb # 加密方式[aes-128-cfb, aes-192-cfb, aes-256-cfb, bf-cfb, cast5-cfb, des-cfb, rc4-md5, rc4-md5-6, chacha20, salsa20, rc4]
password=
timeout=

echo install shadowsocks-server
go get github.com/shadowsocks/shadowsocks-go/cmd/shadowsocks-server

echo install shadowsocks-client
go get github.com/shadowsocks/shadowsocks-go/cmd/shadowsocks-local