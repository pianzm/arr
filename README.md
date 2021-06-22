## Prove of concept
## Stack used
  - Go 1.15+
  - Echo Framework v3
  - Kafka
  - Zookeeper
  - Postgres 9.6
  - MongoDB

## Available endpoint
  - GET /v1/members/entity
  - GET /v1/members/entity/status
  - GET /v1/members/entity/download
## Build Application Image
  ```shell
    $ copy ~/.ssh/id_rsa .
    $ make app-image
  ```

## Docker Compose
  ```shell
    $ docker compose up -d
  ```

## Test
  ```shell
    $ make test
  ```