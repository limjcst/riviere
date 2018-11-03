# Rivière

[![Build Status](https://travis-ci.org/limjcst/riviere.svg?branch=master)](https://travis-ci.org/limjcst/riviere)
[![Coverage Status](https://coveralls.io/repos/github/limjcst/riviere/badge.svg?branch=master)](https://coveralls.io/github/limjcst/riviere?branch=master)

## Overview

Rivière is a tool aiming at forwarding connections between ports and remote ports.
And the ability of supporting hot editing is the most important difference between this tool and other routing or proxying softwares.

## Usage

Run `make` to generate executale file.
`GOPATH` is required by `go-swagger` to create `swagger.json`, which describes the usages of APIs.

As for the database, `sqlite3` and `postgres` are avaliable.
To extend the supported database list, import the driver in `config/config.go`, and compile again.
