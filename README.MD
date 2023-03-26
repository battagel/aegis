# Antivirus

Antivirus service for MinIO

## Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Make Commands](#make-commands)

## Requirements

The only requirement for this program is [Golang 1.16+](https://go.dev/dl/)

## Installation

Clone the source

```
git clone git@github.com:battagel/antivirus.git
```

Build and start the service

```
make build && make run
```

## Make Commands

```
help     Print this message
build    Create the binary
run      Run the binary
vendor   Download the vendored dependencies
```
