# Kong-To-APISIX


[![Go Report Card](https://goreportcard.com/badge/github.com/api7/kong-to-apisix)](https://goreportcard.com/report/github.com/api7/kong-to-apisix)
[![Build Status](https://github.com/api7/kong-to-apisix/actions/workflows/e2e.yml/badge.svg)](https://github.com/api7/kong-to-apisix/actions)
[![Codecov](https://codecov.io/gh/api7/kong-to-apisix/branch/master/graph/badge.svg)](https://codecov.io/gh/api7/kong-to-apisix)

Kong-To-APISIX is a migration tool helping you migrate configuration data of your API gateway from Kong to Apache APISIX. It aims to help people to dip their toes in APISIX and also reduce the operations cost.

Only tested with APISIX 2.7 and Kong 2.4 for now.

## How to use
1. Setup APISIX and Kong, if you don't have them

   Recommend to use `docker compose` to deploy APISIX or Kong:
   - APISIX docker compose guide: https://github.com/apache/apisix-docker#quickstart-via-docker-compose
   - Kong docker compose guide: https://github.com/Kong/docker-kong/tree/master/compose

2. Dump Kong Configuration with Deck. See https://docs.konghq.com/deck/1.7.x/guides/backup-restore/ for details.

3. Run Kong-To-APISIX, and it would generate `apisix.yaml` as declarative configuration file for APISIX.

   ```shell
   make build
   export KONG_YAML_PATH="/PATH/TO/YOUR/Kong.yaml"
   ./bin/kong2apisix
   ```

4. Configure APISIX using `apisix.yaml` by move it to `/PATH/TO/APISIX/conf/apisix.yaml`. Add the following to `config.yaml` at `/PATH/TO/APISIX/conf/config.yaml`:
    ```yaml
    apisix:
        config_center: yaml
        enable_admin: false
    ```

    If you deploy APISIX with docker compose, you need to add `apisix.yaml` to volumes. You could change docker-compose.yml and re-do `docker-compose up`
    ```yaml
    volumes:
      - ./apisix_log:/usr/local/apisix/logs
      - ./apisix_conf/config.yaml:/usr/local/apisix/conf/config.yaml:ro
      - ./apisix_conf/apisix.yaml:/usr/local/apisix/conf/apisix.yaml:ro
    ```

5. Reload APISIX to make declarative configuration work and now test with your new API Gateway
   ```shell
   /PATH/TO/APISIX/bin/apisix reload
   ```

## Demo

1. Make sure you have docker running, and then setup apisix and kong
    ```shell
    cd kong-to-apisix
    ./tools/setup.sh
    ```

2. Follow https://docs.konghq.com/getting-started-guide/2.4.x/overview/ to
   1. Expose services using Service and Route objects
   2. Set up rate limits and proxy caching
   3. Secure services with key authentication
   4. Set up load balancing
    ```shell
    ./examples/kong-example.sh
    ```

3. Dump kong configuration to `kong.yaml`
   ```shell
   go run ./cmd/dumpkong/main.go
   ```

4. Run migration tool, import `kong.yaml` and generate `apisix.yaml` for apisix to use
    ```shell
    export EXPORT_PATH=./repos/apisix-docker/example/apisix_conf
    go run ./cmd/kong-to-apisix/main.go
    ```

5. Verify migration succeeds
    ```shell
    ./examples/apisix-verification.sh
    ```

## Roadmap
- [ ] Improving and completing current apis, eg. support tcp/tls in kong to stream route in APISIX
- [ ] Provide migration report, to declare what has been migrated and those currently not supported
- [ ] Support sni, certificates, ca_certificates configuration migration
- [ ] Support 15+ common plugins
- [ ] Support customized plugin migration
- [ ] Support Incremental migration
