#!/usr/bin/env bash

BASEDIR=$(dirname $(dirname $(dirname "$0")))

fetch_docker_repos() {
    mkdir -p ${BASEDIR}/repos
    if [[ ! -d ${BASEDIR}"/repos/apisix-docker" ]]; then
        git clone https://github.com/apache/apisix-docker.git ${BASEDIR}/repos/apisix-docker --depth=1
        chmod 777 ${BASEDIR}/repos/apisix-docker/example/etcd_data
        cp ${BASEDIR}/examples/config.yaml ${BASEDIR}/repos/apisix-docker/example/apisix_conf/config.yaml
        touch ${BASEDIR}/repos/apisix-docker/example/apisix_conf/apisix.yaml
        sed -i '/apisix_conf/a \      - ./apisix_conf/apisix.yaml:/usr/local/apisix/conf/apisix.yaml:ro' ${BASEDIR}/repos/apisix-docker/example/docker-compose.yml
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

    if [ $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:8088/get) -eq 200 ]; then
        echo "upstream work as expected"
    else
        docker run --name httpbin -d -p 8088:80 kennethreitz/httpbin
        count=0
        while [ $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:8088/get) -ne 200 ];
        do
            echo "Waiting for httpbin setup" && sleep 2;

            ((count=count+1))
            if [ $count -gt ${retries} ]; then
                echo $(curl -k -i -m 20 -o /dev/null -s -w %{http_code} http://localhost:8088/get)
                printf "upstream not work as expected\n"
                docker ps -a
                exit 1
            fi
        done
    fi
}

setup_with_docker_compose
