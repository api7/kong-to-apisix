# Kong-To-APISIX


[![Go Report Card](https://goreportcard.com/badge/github.com/api7/kong-to-apisix)](https://goreportcard.com/report/github.com/api7/kong-to-apisix)
[![Build Status](https://github.com/api7/kong-to-apisix/actions/workflows/test-ci.yml/badge.svg)](https://github.com/api7/kong-to-apisix/actions)
[![Codecov](https://codecov.io/gh/api7/kong-to-apisix/branch/main/graph/badge.svg)](https://codecov.io/gh/api7/kong-to-apisix)

Kong-To-APISIX is a migration tool helping you migrate configuration data of your API gateway from Kong to Apache APISIX. It aims to help people to dip their toes in APISIX and also reduce the operations cost.

Only tested with APISIX 2.8 and Kong 2.4 for now.

## How to use

1. Dump Kong Configuration with `Deck` or `CLI`, for detailed usage, please refer to:
   - Kong Deck : https://docs.konghq.com/deck/1.7.x/guides/backup-restore/
   - Kong CLI Config Export: https://docs.konghq.com/gateway-oss/2.4.x/cli/#kong-config

2. Build `Kong to APISIX`, go version require `1.16+`。

3. Run Kong-To-APISIX, and it would generate `apisix.yaml` as declarative configuration file for APISIX.

   ```shell
   $ make build
   $ ./bin/kong-to-apisix migrate -i kong.yaml -o apisix.yaml
   migrate succeed
   ```

4. Configure APISIX with `apisix.yaml`, see https://apisix.apache.org/docs/apisix/stand-alone for details.

If more help needed, you could refer [detail steps](docs/how-to-deploy.md)


## Support features

1. Kong service is converted to APISIX service (including: ID, name, retry, protocol, timeout, path, port, host (default upstream))
2. Kong route is converted to APISIX route (including: ID, name, methods, hosts, paths(Path handling algorithms), regex_priority)
3. Kong upstream is converted to APISIX upstream (including: ID, name, algorithm, upstream.target.target, upstream.target.weight)
4. Kong consumer is converted to APISIX consumer (including: ID, username, custom_id, plugins.keyauth_credentials, plugins.basicauth_credentials, plugins.hmacauth_credentials, plugins.jwt_secrets)
5. Kong plugin is converted to APISIX plugin (including: key-auth, rate-limiting, proxy-cache)
6. Kong global plugin is converted to APISIX global_rule (including: key-auth, rate-limiting, proxy-cache)


## Roadmap
- [ ] Improving and completing current apis, eg. support tcp/tls in kong to stream route in APISIX
- [ ] Provide migration report, to declare what has been migrated and those currently not supported
- [ ] Support sni, certificates, ca_certificates configuration migration
- [ ] Support 15+ common plugins
- [ ] Support customized plugin migration
- [ ] Support Incremental migration
