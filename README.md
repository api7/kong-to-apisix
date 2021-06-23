# kong-to-apisix

## demo

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
    go run ./cmd/kong-to-apisix/demo.go
    ```

4. Verify migration succeeds
    ```shell
    ./examples/apisix-verification.sh
    ```
