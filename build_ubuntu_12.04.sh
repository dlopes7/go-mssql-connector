#!/usr/bin/env bash

docker run -v $(pwd):/src --rm dlopes7/go-ubuntu:latest go build -v -o releases/go-mssql-connector