#!/usr/bin/env bash

BASEDIR=$(dirname $(dirname $(dirname "$0")))

fetch_docker_repos() {
    mkdir -p ${BASEDIR}/repos
    if [[ ! -d ${BASEDIR}"/repos/apisix-docker" ]]; then
        git clone https://github.com/apache/apisix-docker.git ${BASEDIR}/repos/apisix-docker --depth=1
        chmod 777 ${BASEDIR}/repos/apisix-docker/example/etcd_data
    fi

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

    setup_app "apisix" "9080" "404" "docker-compose -f ${BASEDIR}/repos/apisix-docker/example/docker-compose.yml up -d"
    setup_app "kong" "8001" "200" "docker-compose -f ${BASEDIR}/repos/kong-docker/compose/docker-compose.yml up -d"
    setup_app "httpbin" "8088" "200" "docker run --name httpbin1 -d -p 8088:80 kennethreitz/httpbin"
}

setup_app() {
    local name=$1
    local port=$2
    local expect_code=$3
    local command=$4

    retries=10

    if [ $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:${port}) -eq $expect_code ]; then
        echo "${name} work as expected"
    else
        $command
        count=0
        while [ $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:${port}) -ne $expect_code ];
        do
            echo "Waiting for ${name} setup" && sleep 2;

            ((count=count+1))
            if [ $count -gt ${retries} ]; then
                echo $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:${port})
                printf "${name} not work as expected\n"
                docker ps -a
                exit 1
            fi
        done
    fi
}

setup_with_docker_compose
