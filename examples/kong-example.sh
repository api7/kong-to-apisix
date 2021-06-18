#!/usr/bin/env bash

set -ex

# Follow https://docs.konghq.com/getting-started-guide/2.4.x/overview/
# to
#   1. Expose your services using Service and Route objects
#   2. Set up rate limits and proxy caching
#   3. Secure services with key authentication
#   4. Set up load balancing

set_service() {
    echo "set service and route"

    curl -i -X POST http://127.0.0.1:8001/services \
        --data name=example_service \
        --data url='http://mockbin.org'

    curl -i -X POST http://127.0.0.1:8001/services/example_service/routes \
        --data 'paths[]=/mock' \
        --data name=mocking

    code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:8000/mock/request)
    if [ $code -eq 200 ]; then
        echo "route succeeded"
    else
        echo "route failed"
        exit 1
    fi
}

set_limit_rate() {
    echo "set plugin - rate limit"
    curl -s -o /dev/null -X DELETE http://127.0.0.1:8001/rate-limiting

    curl -i -X POST http://127.0.0.1:8001/plugins \
        --data name=rate-limiting \
        --data config.second=5 \
        --data config.policy=local

    for i in {1..5}; do
        curl -s -o /dev/null -X GET http://127.0.0.1:8000/mock/request
    done

    code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:8000/mock/request)
    if [ $code -eq 429 ]; then
        echo "rate limit succeeded"
    else
        echo "rate limit failed"
        exit 1
    fi
}

set_proxy_cache() {
    echo "set plugin - proxy cache"
    curl -s -o /dev/null -X DELETE http://127.0.0.1:8001/proxy-cache

    curl -i -X POST http://127.0.0.1:8001/plugins \
        --data name=proxy-cache \
        --data config.content_type="application/json; charset=utf-8" \
        --data config.cache_ttl=1 \
        --data config.strategy=memory

    latency1=$(curl -i -s -X GET http://127.0.0.1:8000/mock/request | grep "X-Kong-Upstream-Latency" | awk '{print $2}' | tr -d '\r')
    latency2=$(curl -i -s -X GET http://127.0.0.1:8000/mock/request | grep "X-Kong-Upstream-Latency" | awk '{print $2}' | tr -d '\r')

    echo "first latency: "$latency1"; second latency: "$latency2
    if [ $latency1 -gt $latency2 ]; then
        echo "proxy cache succeeded"
    else
        echo "proxy cache failed"
        exit 1
    fi
}

set_key_auth() {
    echo "set plugin - key auth"

    curl -i -X POST http://127.0.0.1:8001/routes/mocking/plugins \
        --data name=key-auth

    code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:8000/mock)
    if [ $code -ne 401 ]; then
        echo "set key auth failed"
        exit 1
    fi

    curl -i -X POST http://127.0.0.1:8001/consumers/ \
        --data username=consumer \
        --data custom_id=consumer

    curl -i -X POST http://127.0.0.1:8001/consumers/consumer/key-auth \
        --data key=apikey

    code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:8000/mock/request -H 'apikey:apikey')
    if [ $code -eq 200 ]; then
        echo "key auth succeeded"
    else
        echo "key auth failed"
        exit 1
    fi
}

set_load_balancing() {
    echo "set load balancing"

    curl -X POST http://127.0.0.1:8001/upstreams \
        --data name=upstream

    curl -X PATCH http://127.0.0.1:8001/services/example_service \
        --data host='upstream'

    curl -X POST http://127.0.0.1:8001/upstreams/upstream/targets \
        --data target='mockbin.org:80'

    curl -X POST http://127.0.0.1:8001/upstreams/upstream/targets \
        --data target='httpbin.org:80'

    httpbin_num=0
    mockbin_num=0
    set +x
    for i in {1..6}; do
        body=$(curl -k -i -s http://127.0.0.1:8000/mock -H 'apikey:apikey')
        if [[ $body == *"httpbin"* ]]; then
            httpbin_num=$((httpbin_num+1))
        elif [[ $body == *"mockbin"* ]]; then
            mockbin_num=$((mockbin_num+1))
        fi
        sleep 1.1
    done
    set -x
    echo "httpbin number: "${httpbin_num}", mockbin number: "${mockbin_num}

    if [[ $httpbin_num -gt 0 && $mockbin_num -gt 0 ]]; then
        echo "load balancing succeeded"
    else
        echo "load balancing failed"
    fi
}

# "$@"

set_service
set_limit_rate
set_proxy_cache
set_key_auth
set_load_balancing
