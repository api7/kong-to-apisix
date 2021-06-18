# kong-to-apisix

## demo

Make sure you have docker running, and then setup apisix and kong
```shell
./tools/setup.sh
```

Follow https://docs.konghq.com/getting-started-guide/2.4.x/overview/ to
1. Expose services using Service and Route objects
2. Set up rate limits and proxy caching
3. Secure services with key authentication
4. Set up load balancing
```shell
./examples/kong-example.sh set_service
set_limit_rate
set_proxy_cache
set_key_auth
set_load_balancing
```

Run migration tool
```shell
go run ./examples/demo.go
```

Verify migration succeeds
```shell
./examples/apisix-verification.sh
```
