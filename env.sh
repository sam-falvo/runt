#!/bin/bash
export GOPATH=`pwd`
export P=github.com/sam-falvo/tx
export SP=src/$P
export PP=pkg/$P
export I="go install $P/runt"
export M="vim $SP/runt/main.go"
export D="vim $SP/driver/driver.go"
export DT="vim $SP/driver/driver_test.go"

