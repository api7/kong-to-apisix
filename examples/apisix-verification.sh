#!/usr/bin/env bash

set -ex

# something wrong with quoted line
sed -i -e 's#- http://etcd:2379#- "http://etcd:2379"#g' "repos/apisix-docker/example/apisix_conf/config.yaml"
# config.yaml is read only in docker, so we need to rebuild apisix with docker compose
docker-compose -f repos/apisix-docker/example/docker-compose.yml down
docker-compose -f repos/apisix-docker/example/docker-compose.yml up -d

# test key auth
code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:9080/mock)
if [ $code -eq 401 ]; then
    echo "key-auth take effect"
else
    echo "fail: key-auth not take effect"
    exit 1
fi

# test proxy cache
curl -k -i -s  -o /dev/null http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org"
hit=$(curl -i -s -X GET http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org" | grep "Apisix-Cache-Status" | awk '{print $2}' | tr -d '\r')

if [ "$hit" == "HIT" ]; then
    echo "proxy-cache take effect"
else
    echo "fail: proxy-cache not take effect"
    exit 1
fi

# test limit count
for i in {1..5}; do
    curl -s -o /dev/null -X GET http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org"
done

code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org")
if [ $code -eq 429 ]; then
    echo "limit count take effect"
else
    echo "fail: limit count not take effect"
    exit 1
fi

# it seems mockbin have adapted to kong, so kong don't need to explicitly set host when accessing mockbin
httpbin_num=0
mockbin_num=0
set +x
for i in {1..6}; do
    body=$(curl -k -i -s http://127.0.0.1:9080/mock -H "apikey: apikey" -H "Host: mockbin.org")
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
    echo "load balancing take effect"
else
    echo "fail: load balancing not take effect"
fi
