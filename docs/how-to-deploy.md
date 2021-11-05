# How to deploy

1. Setup APISIX and Kong, if you don't have them

   Recommend to use `docker compose` to deploy APISIX or Kong:
   - APISIX docker compose guide: https://github.com/apache/apisix-docker#quickstart-via-docker-compose
   - Kong docker compose guide: https://github.com/Kong/docker-kong/tree/master/compose

2. Dump Kong Configuration with `Deck` or `CLI`, for detailed usage, please refer to:
   - Kong Deck : https://docs.konghq.com/deck/1.7.x/guides/backup-restore/
   - Kong CLI Config Export: https://docs.konghq.com/gateway-oss/2.4.x/cli/#kong-config

3. Run Kong-To-APISIX, and it would generate `apisix.yaml` as declarative configuration file for APISIX.

   ```shell
   make build
   export KONG_YAML_PATH="/PATH/TO/YOUR/Kong.yaml"
   ./bin/kong-to-apisix
   ```

4. Configure APISIX using `apisix.yaml` by move it to `/PATH/TO/APISIX/conf/apisix.yaml`. Add the following to `config.yaml` at `/PATH/TO/APISIX/conf/config.yaml`:
    ```yaml
    apisix:
        config_center: yaml
        enable_admin: false
    ```

    If you deploy APISIX with docker compose, you need to add `apisix.yaml` to volumes. You could change docker-compose.yml and re-do `docker-compose up`
    ```yaml
    volumes:
      - ./apisix_log:/usr/local/apisix/logs
      - ./apisix_conf/config.yaml:/usr/local/apisix/conf/config.yaml:ro
      - ./apisix_conf/apisix.yaml:/usr/local/apisix/conf/apisix.yaml:ro
    ```

5. Reload APISIX to make declarative configuration work and now test with your new API Gateway
   ```shell
   /PATH/TO/APISIX/bin/apisix reload
   ```
