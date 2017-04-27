#!/bin/bash
GOPATH=/home/salil/gopath:/home/salil/workspaces/cuckoo-go/cuckoofilter INTBITS=64 CGO_LDFLAGS=-lcrypto go install -x
