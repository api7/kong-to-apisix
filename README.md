# Kong-To-APISIX


[![Go Report Card](https://goreportcard.com/badge/github.com/api7/kong-to-apisix)](https://goreportcard.com/report/github.com/api7/kong-to-apisix)
[![Build Status](https://github.com/api7/kong-to-apisix/actions/workflows/e2e.yml/badge.svg)](https://github.com/api7/kong-to-apisix/actions)
[![Codecov](https://codecov.io/gh/api7/kong-to-apisix/branch/master/graph/badge.svg)](https://codecov.io/gh/api7/kong-to-apisix)

Kong-To-APISIX is a migration tool helping you migrate configuration data of your API gateway from Kong to Apache APISIX. It aims to help people to dip their toes in APISIX and also reduce the operations cost.

Only tested with APISIX 2.7 and Kong 2.4 for now.

## How to use
1. Dump Kong Configuration with Deck. See https://docs.konghq.com/deck/1.7.x/guides/backup-restore/ for details.

2. Run Kong-To-APISIX, and it would generate `apisix.yaml` as declarative configuration file for APISIX.

   ```shell
   $ make build
   $ ./bin/kong-to-apisix migrate -i kong.yaml -o apisix.yaml
   migrate succeed
   ```

3. Configure APISIX with `apisix.yaml`, see https://apisix.apache.org/docs/apisix/stand-alone for details.

## Demo

1. Make sure you have docker running, and then setup apisix and kong
    ```shell
    $ cd kong-to-apisix
    $ ./tools/setup.sh
    ```

2. Follow https://docs.konghq.com/getting-started-guide/2.4.x/overview/ to
   1. Expose services using Service and Route objects
   2. Set up rate limits and proxy caching
   3. Secure services with key authentication
   4. Set up load balancing
    ```shell
    $./examples/kong-example.sh
    ```

3. Dump kong configuration to `kong.yaml`
   ```shell
   $ make build
   $ ./bin/kong-to-apisix dump -o kong.yaml
   generated kong configuration file at kong.yaml
   ```

4. Run migration tool, import `kong.yaml` and generate `apisix.yaml` for apisix to use
    ```shell
    $ ./bin/kong-to-apisix migrate -i kong.yaml -o ./repos/apisix-docker/example/apisix_conf/apisix.yaml -c ./repos/apisix-docker/example/apisix_conf/config.yaml
    migrate succeed
    ```

5. Verify migration succeeds
    1. test key auth
    ```shell
    curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:9080/mock
    # output: 401
    ```
    2. test proxy cache
    ```shell
    # access for the first time
    curl -k -I -s  -o /dev/null http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org"
    # see if got cached
    curl -I -s -X GET http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org"
    # output:
    #   HTTP/1.1 200 OK
    #   ...
    #   Apisix-Cache-Status: HIT
    ```

    3. test limit count
    ```shell
    for i in {1..5}; do
        curl -s -o /dev/null -X GET http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org"
    done
    curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org"
    # output: 429
    ```

    4. test load balancing
    ```shell
    httpbin_num=0
    mockbin_num=0
    for i in {1..8}; do
        body=$(curl -k -i -s http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org")
        if [[ $body == *"httpbin"* ]]; then
            httpbin_num=$((httpbin_num+1))
        elif [[ $body == *"mockbin"* ]]; then
            mockbin_num=$((mockbin_num+1))
        fi
        sleep 1.5
    done
    echo "httpbin number: "${httpbin_num}", mockbin number: "${mockbin_num}
    # output:
    #   httpbin number: 6, mockbin number: 2
    ```

## Roadmap
- [ ] Improving and completing current apis, eg. support tcp/tls in kong to stream route in APISIX
- [ ] Provide migration report, to declare what has been migrated and those currently not supported
- [ ] Support sni, certificates, ca_certificates configuration migration
- [ ] Support 15+ common plugins
- [ ] Support customized plugin migration
- [ ] Support Incremental migration
