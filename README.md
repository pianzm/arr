## Implementing Asynchronous Request-Reply using Redis PubSub
## Stack used
  - Go 1.15+
  - Echo Framework v3
  - Redis

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
