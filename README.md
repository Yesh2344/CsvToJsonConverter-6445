# CsvToJsonConverter

[![Go Reference](https://pkg.go.dev/badge/github.com/example/csvtojsonconverter.svg)](https://pkg.go.dev/github.com/example/csvtojsonconverter)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

CsvToJsonConverter is a production‑ready command line tool that converts CSV files to JSON while performing schema‑based validation. It demonstrates modern Go practices, including clean architecture, configuration management, structured logging, and comprehensive testing.

## Features

- **Schema validation** – Define column names and types in `config.yaml`.
- **Zero‑dependency CSV parsing** – Uses the standard library.
- **Structured logging** – Powered by Uber’s `zap`.
- **Configuration via Viper** – Supports environment variables and YAML config.
- **CLI usage** – Simple `csv2json` command.
- **HTTP API (optional)** – Expose conversion as a REST endpoint.

## Installation