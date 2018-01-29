#!/bin/bash -ex

export GOPATH=$(pwd)

go get cue
go install cue

