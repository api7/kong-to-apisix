# Kong-To-APISIX

Kong-To-APISIX is a migration tool helping you migrate configuration data of your API gateway from Kong to Apache APISIX. It aims to help people to dip their toes in APISIX and also reduce the operations cost.

## How to use
1. Setup APISIX and Kong, if you don't have them

   Recommend to use `docker compose` to deploy APISIX or Kong:
   - APISIX docker compose guide: https://github.com/apache/apisix-docker#quickstart-via-docker-compose
   - Kong docker compose guide: https://github.com/Kong/docker-kong/tree/master/compose

2. Set address of APISIX and Kong Admin API. For example:

   ```shell
   export APISIX_ADMIN_ADDR="http://127.0.0.1:9080"
   export APISIX_ADMIN_TOKEN="edd1c9f034335f136f87ad84b625c8f1"
   export KONG_ADMIN_ADDR="http://127.0.0.1:8001"
   ```

3. Run Kong-To-APISIX

   ```shell
   make build
   ./bin/kta
   ```

## Demo

1. Make sure you have docker running, and then setup apisix and kong
    ```shell
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

3. Run migration tool
    ```shell
    go run ./cmd/kong-to-apisix/main.go
    ```

4. Verify migration succeeds
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
