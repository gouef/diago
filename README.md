<img align=right width="168" src="docs/gouef_logo.png">

# Diago
Diago: the addictive tool. Friendly design, logging

[![GoDoc](https://pkg.go.dev/badge/github.com/gouef/diago.svg)](https://pkg.go.dev/github.com/gouef/diago)
[![GitHub stars](https://img.shields.io/github/stars/gouef/diago?style=social)](https://github.com/gouef/diago/stargazers)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouef/diago)](https://goreportcard.com/report/github.com/gouef/diago)

## How it looks ? 
![Diago](docs/diago.png)


## Introduction
Diago is diagnostic panel wrote in go for websites. It should help you improve, debug your application.

## Latency Extension
⚠️ If you want use Latency Extension than You must register it as first!

## Route Extension
Route extension was moved to [gouef/router](github.com/gouef/router)

## Custom Extension
You can write your own extension, Diago is ready for it!

You must use interface [DiagoExtension](/diagoExtension.go).

For inspiration I recommend look to exists extensions.