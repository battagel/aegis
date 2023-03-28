# Antivirus

Antivirus manager service for MinIO.

## Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Make Commands](#make-commands)
- [Usage](#usage)

## Requirements

- GoLang
- MinIO server
- Kafka server
- ClamAV daemon

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

## Usage

This service manages the scanning of new files uploaded to MinIO. It is connected to Kafka to read PUT events, download the file into a cache and then scan. Depending on the output, the files metadata will be updated with the result of the scan.
