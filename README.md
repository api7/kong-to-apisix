# Kong-To-APISIX


[![Go Report Card](https://goreportcard.com/badge/github.com/api7/kong-to-apisix)](https://goreportcard.com/report/github.com/api7/kong-to-apisix)
[![Build Status](https://github.com/api7/kong-to-apisix/actions/workflows/e2e-test.yml/badge.svg)](https://github.com/api7/kong-to-apisix/actions)
[![Codecov](https://codecov.io/gh/api7/kong-to-apisix/branch/master/graph/badge.svg)](https://codecov.io/gh/api7/kong-to-apisix)

Kong-To-APISIX is a migration tool helping you migrate configuration data of your API gateway from Kong to Apache APISIX. It aims to help people to dip their toes in APISIX and also reduce the operations cost.

Only tested with APISIX 2.7 and Kong 2.4 for now.

## How to use
1. Dump Kong Configuration with Deck. See https://docs.konghq.com/deck/1.7.x/guides/backup-restore/ for details.

2. Build `Kong to APISIX`, go version require `1.16+`。

3. Run Kong-To-APISIX, and it would generate `apisix.yaml` as declarative configuration file for APISIX.

   ```shell
   $ make build
   $ ./bin/kong-to-apisix migrate -i kong.yaml -o apisix.yaml
   migrate succeed
   ```

4. Configure APISIX with `apisix.yaml`, see https://apisix.apache.org/docs/apisix/stand-alone for details.

If more help needed, you could refer [detail steps](docs/how-to-deploy.md)

## Roadmap
- [ ] Improving and completing current apis, eg. support tcp/tls in kong to stream route in APISIX
- [ ] Provide migration report, to declare what has been migrated and those currently not supported
- [ ] Support sni, certificates, ca_certificates configuration migration
- [ ] Support 15+ common plugins
- [ ] Support customized plugin migration
- [ ] Support Incremental migration
