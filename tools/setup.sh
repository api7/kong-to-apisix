#!/usr/bin/env bash

BASEDIR=$(dirname "$0")/..

fetch_docker_repos() {
    mkdir -p ${BASEDIR}/repos
    if [[ ! -d ${BASEDIR}"/repos/apisix-docker" ]]; then
        git clone https://github.com/apache/apisix-docker.git ${BASEDIR}/repos/apisix-docker --depth=1
        chmod 777 ${BASEDIR}/repos/apisix-docker/example/etcd_data
        cp ${BASEDIR}/examples/config.yaml ${BASEDIR}/repos/apisix-docker/example/apisix_conf/config.yaml
        cp ${BASEDIR}/examples/apisix.yaml ${BASEDIR}/repos/apisix-docker/example/apisix_conf/apisix.yaml
        if [ $(uname) = "Darwin" ]; then
            which gsed || brew install gnu-sed
            gsed -i '/apisix_conf/a \      - ./apisix_conf/apisix.yaml:/usr/local/apisix/conf/apisix.yaml:ro' ${BASEDIR}/repos/apisix-docker/example/docker-compose.yml
        else
            sed -i '/apisix_conf/a \      - ./apisix_conf/apisix.yaml:/usr/local/apisix/conf/apisix.yaml:ro' ${BASEDIR}/repos/apisix-docker/example/docker-compose.yml
        fi
    fi

    if [[ ! -d ${BASEDIR}"/repos/kong-docker" ]]; then
        git clone --depth=1 --branch 2.5.1 https://github.com/Kong/docker-kong.git ${BASEDIR}/repos/kong-docker
        mkdir -p ${BASEDIR}/repos/kong-docker/compose/kong_conf
        chmod -R 777 ${BASEDIR}/repos/kong-docker/compose/kong_conf
        sed -i '/user:/a \    working_dir: /config/kong' ${BASEDIR}/repos/kong-docker/compose/docker-compose.yml
        sed -i '/user:/a \    container_name: kong' ${BASEDIR}/repos/kong-docker/compose/docker-compose.yml
        sed -i '/security_opt:/i \      - ./kong_conf:/config/kong:rw' ${BASEDIR}/repos/kong-docker/compose/docker-compose.yml
    fi
}

fetch_docker_repos

docker ps > /dev/null
if [ $? -ne 0 ]; then
    echo "docker not working"
    exit 1
fi

retries=10
if [ $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:9080) -eq 404 ]; then
    echo "apisix work as expected"
else
    docker-compose -f ${BASEDIR}/repos/apisix-docker/example/docker-compose.yml up -d
    count=0
    while [ $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:9080) -ne 404 ];
    do
        echo "Waiting for apisix setup" && sleep 2;

        ((count=count+1))
        if [ $count -gt ${retries} ]; then
            echo $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:9080)
            echo "apisix not work as expected"
            docker ps -a
            exit 1
        fi
    done
    echo "apisix work as expected"
fi

if [ $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:8001) -eq 200 ]; then
    echo "kong work as expected"
else
    docker-compose -f ${BASEDIR}/repos/kong-docker/compose/docker-compose.yml up -d
    count=0
    while [ $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:8001) -ne 200 ];
    do
        echo "Waiting for kong setup" && sleep 2;

        ((count=count+1))
        if [ $count -gt ${retries} ]; then
            echo $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:8001)
            printf "kong not work as expected\n"
            docker ps -a
            exit 1
        fi
    done
    echo "kong work as expected"
fi

if [ ! -z "$1" ]; then
    echo "set upstream"
    docker container inspect upstream > /dev/null 2>&1 \
    || docker run -itd --name upstream -v $(pwd)/examples/conf:/etc/nginx/conf.d -p 7024:7024 -p 7025:7025 -p 7026:7026 openresty/openresty:alpine
fi
