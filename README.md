# DEPRECATED
This repository is part of the initial [multiserver](https://github.com/HimbeerserverDE/multiserver) project
and should not be used. Its successor is the [mt-auth-convert](https://github.com/HimbeerserverDE/mt-multiserver-proxy/blob/main/doc/auth_backends.md#mt-auth-convert) tool
of the [mt-multiserver-proxy](https://github.com/HimbeerserverDE/mt-multiserver-proxy) project.

# multiserver_converter
MT to multiserver auth database converter

## Installation
`go get github.com/HimbeerserverDE/multiserver_converter`

## Usage
`$GOPATH/bin/multiserver_converter <sqlite3 <in> <out> | psql <in_db> <in_user> <in_password> <in_host> <in_port> <out_db> <out_user> <out_password> <out_host> <out_port>>`
