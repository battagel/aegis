# Aegis for MinIO

This service manages the scanning of new files uploaded to MinIO. It is connected to Kafka to read PUT events, download the file into a cache and then scan. Depending on the output, the files tags will be updated with the result of the scan.

## Contents

- [Setup](#setup)
- [Config](#config)
- [Make Commands](#make-commands)
- [Usage](#usage)

## Setup

### Local Environment

External dependencies that need to be pre-installed:

- Go
- MinIO
- Prometheus
- Kafka
- Clamav
- Postgresql
- Grafana (optional)

Once installed, start the external dependencies. This can be done using the
commands below. Prometheus requires a config file to be passed in shown below.
Any changes to the following commands may require changes to the config.env.

MinIO:

```console
minio server <path_to_data>
```

Kafka:

```console
kafka-topics.sh --create --topic minio-put-event --bootstrap-server localhost:9092
```

Prometheus

```console
prometheus --config.file=./local-configs/prometheus.yml
```

Postgresql

```console
pg_ctl -D <postgresql_location> start
```

Finally, start the Aegis service.

```console
make build && make run
```

Thats it! Navigate to [http://localhost:9000](http://localhost:9000) to view
interact with MinIO. Checkout the [usage](#usage) section for how to configure the MinIO client.

### Kubernetes Environment

External dependencies that need to be pre-installed:

- k3d
- Helm
- Docker
- Kubernetes

Make sure dependencies are installed and configured correctly.

From within project root you can use the make command to launch the kubernetes environment created by k3d.

```console
make docker-build && make create-cluster
```

Using kubectl you can verify that all pods and services have started correctly.

```console
kubectl get pods
kubectl get svc
```

Navigate to [http://localhost:9001](http://localhost:9001) to view interact with
MinIO. Checkout the [usage](#usage) section for how to configure the MinIO client.

Default credentials:
user: minioadmin
password: minioadmin

## Testing

Dependencies:

- Go
- mockery

```console
make mock && make test
```

## Config

All configuration is stored within the config.env file but can be overridden by environment variables.

Default configuration

```env
LOGGER_LEVEL=info
LOGGER_ENCODING=console

MINIO_ENDPOINT=127.0.0.1:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false

KAFKA_BROKERS="127.0.0.1:9092"
KAFKA_TOPIC=minio-put-event
KAFKA_GROUP_ID=g1
KAFKA_MAX_BYTES=10

CLAMAV_REMOVE_AFTER_SCAN=true
CLAMAV_DATETIME_FORMAT="01-02-2006 15:04:05"
CLAMAV_PATH=cache/

PROMETHEUS_ENDPOINT=127.0.0.1:2112
PROMETHEUS_PATH=/metrics

POSTGRES_USERNAME=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_ENDPOINT=127.0.0.1:5432
POSTGRES_DATABASE=aegis_antivirus
POSTGRES_TABLE=aegis_audit_logs
```

## Make Commands

```console
help             Print this message
build            Create the binary
run              Run the binary
vendor           Download the vendored dependencies
test             Run the tests
mock             Generate the mocks for testing
docker-build     Build the docker image
create-cluster   Create the k3d cluster
delete-cluster   Delete the k3d cluster
rebuild-cluster  Delete and recreate the cluster
```

## Usage

### Overview

Aegis is designed to bring security to the open-source object store MinIO. Aegis listens for PUT notifications from MinIO which passed on to Kafkas event queue. These notifications are then read in by Aegis where files are retrieved into a temporary cache. Aegis then distributes the files to various antivirus engines for scanning. The results are collated and then the file is tagged with the results. A Postgresql audit log is also collecting to track change over time. Metrics are also collected and exposed to Prometheus which can be consumed in Grafana for visualisation.

### MinIO

When the Aegis local or cluster setup has been created, MinIO need to be configured to send notifications to Kafka. In dashboard > Events > Add Event Destination. Select Kafka as the queue type and enter the brokers address and topic name. The brokers address change depending on your setup but these are normally 127.0.0.1:9092 for local or aegis-kafka.default.svc.cluster.local:9092 for kubernetes. The topic name can be anything you like but it must match the topic name in the config.env file.

Next, you should configure every bucket that requires file scanning. Buckets > <bucket> > Settings > Events > Subscribe to Event. From the drop down, select the Kafka event queue ARN and also the PUT event type. This will send a notification to Kafka every time a file is uploaded to the bucket.

This should be all the setup you need to utilise Aegis in your object storage stack.

### Grafana

Prometheus metrics are exposed to [http://localhost:9090](http://localhost:9090). These can be consumed in Grafana to visualise the data.

### ClamAV

When starting the deployment within kubernetes the clamav will start and request a signature database update. Sometimes this update can trigger too many times causing a request cooldown. This will cause errors with object scanning but can be solved by simply waiting for this cooldown to expire.
