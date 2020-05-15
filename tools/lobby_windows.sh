#!/bin/bash

set -eu

go build -o bin/couchcampaign.exe cmd/game/*
go build -o bin/lobby.exe cmd/lobby/*
PATH="$PATH:$(pwd)/bin" ./bin/lobby.exe
