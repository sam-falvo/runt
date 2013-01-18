#!/bin/bash
export GOPATH=`pwd`
export P=github.com/sam-falvo/tx
export SP=src/$P
export PP=pkg/$P
export I="go install $P/runt"
export M="vim $SP/runt/main.go"

