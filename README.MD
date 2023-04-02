# Aegis

This service manages the scanning of new files uploaded to MinIO. It is connected to Kafka to read PUT events, download the file into a cache and then scan. Depending on the output, the files tags will be updated with the result of the scan.

## Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Make Commands](#make-commands)
- [Usage](#usage)

## Requirements

- Go
- Mockery (testing)
- MinIO Object Store Server
- Kafka Server
- ClamAV Daemon

## Installation

Clone the source

```
git clone git@github.com:battagel/aegis.git
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
