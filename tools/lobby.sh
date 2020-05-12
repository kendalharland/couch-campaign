#!/bin/bash
go build -o bin/couchcampaign cmd/game/*
go build -o bin/lobby cmd/lobby/*
PATH="$PATH:$(pwd)/bin" ./bin/lobby
