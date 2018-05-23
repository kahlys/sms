# Golang SMS-API Wrapper

[![godoc](https://godoc.org/github.com/kahlys/sms?status.svg)](https://godoc.org/github.com/kahlys/sms)
[![build status](https://api.travis-ci.org/kahlys/sms.svg?branch=master)](https://travis-ci.org/kahlys/sms)
[![go report](https://goreportcard.com/badge/github.com/kahlys/sms)](https://goreportcard.com/report/github.com/kahlys/sms)
[![Coverage Status](https://coveralls.io/repos/github/kahlys/sms/badge.svg?branch=master)](https://coveralls.io/github/kahlys/sms?branch=master)

This package provides a generic interface around some api-sms wrappers in order

## Instalation

With a correctly configured [Go toolchain](https://golang.org/doc/install):

```sh
$ go get github.com/kahlys/sms
$ go get github.com/kahlys/driver/<drivername>
```

## Drivers

- **Mobimel**: [Documentation FR](http://www.mobimel.com/envoi-automatise-par-requetes-http) - [Example](#mobimel)
- **OVH**: [Documentation FR](https://docs.ovh.com/fr/sms/envoyer_des_sms_depuis_une_url_-_http2sms/) - [Example](#ovh)

## Examples

### Mobimel

```go
import (
    "github/kahlys/sms"
    _ "github/kahlys/sms/driver/mobimel"
)

func main() {
    param := map[string]string{
        "login":"bruce", 
        "password": "91939", 
        "sender": "wayne",
    }
    sender, _ := sms.Init("mobimel", param)
    sender.Send("Meet me at the roof !", "+33666666666")
}
```

### OVH

```go
import (
    "github.com/kahlys/sms"
    _ "github.com/kahlys/sms/driver/ovh"
)

func main() {
    param := map[string]string{
        "account":"sms-xx4242-7",
        "login":"bruce", 
        "password": "91939", 
        "sender": "wayne",
    }
    sender, _ := sms.Init("ovh", param)
    sender.Send("Meet me at the roof !", "+33666666666")
}
```
