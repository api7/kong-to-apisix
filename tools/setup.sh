#!/usr/bin/env bash

BASEDIR=$(dirname $(dirname $(dirname "$0")))

fetch_docker_repos() {
    mkdir -p ${BASEDIR}/repos
    if [[ ! -d ${BASEDIR}"/repos/apisix-docker" ]]; then
        git clone https://github.com/apache/apisix-docker.git ${BASEDIR}/repos/apisix-docker --depth=1
    fi

    # fix image error for now
    sed -i -e 's#- http://etcd:2379#- "http://etcd:2379"#g' ${BASEDIR}"/repos/apisix-docker/example/apisix_conf/config.yaml"

    if [[ ! -d ${BASEDIR}"/repos/kong-docker" ]]; then
        git clone https://github.com/Kong/docker-kong.git ${BASEDIR}/repos/kong-docker --depth=1
    fi
}

setup_with_docker_compose() {
    fetch_docker_repos

    docker ps > /dev/null
    if [ $? -ne 0 ]; then
        echo "docker not working"
        exit 1
    fi

    docker-compose -f ${BASEDIR}/repos/apisix-docker/example/docker-compose.yml up -d
    docker-compose -f ${BASEDIR}/repos/kong-docker/compose/docker-compose.yml up -d

    code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:9080)
    if [ $code -eq 404 ]; then
        echo "apisix work as expected"
    else
        echo "apisix not work as expected"
    fi

    sleep 3
    code=$(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:8001)
    if [ $code -eq 404 ]; then
        echo "kong work as expected"
    else
        echo "kong not work as expected"
    fi
}

setup_with_docker_compose
