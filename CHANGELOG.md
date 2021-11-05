---
title: Changelog
---

## Table of Contents

- [0.1.0](#010)

## 0.1.0

### Core

- Support both `Kong Deck` and `Kong Config` to export data [#28](https://github.com/api7/kong-to-apisix/pull/28) [#30](https://github.com/api7/kong-to-apisix/pull/30) [#31](https://github.com/api7/kong-to-apisix/pull/31) [#32](https://github.com/api7/kong-to-apisix/pull/32) [#33](https://github.com/api7/kong-to-apisix/pull/33) 
- Support conversion of `kong service` basic data to `APISIX service` data (including: ID, name, retry, protocol, timeout, path, port, host (default upstream)) [#28](https://github.com/api7/kong-to-apisix/pull/28)
- Support conversion of `kong route` basic data to `APISIX route` data (including: ID, name, methods, hosts, paths(Path handling algorithms), regex_priority) [#30](https://github.com/api7/kong-to-apisix/pull/30)
- Support conversion of `kong upstream` basic data to `APISIX upstream` data (including: ID, name, algorithm, upstream.target.target, upstream.target.weight) [#33](https://github.com/api7/kong-to-apisix/pull/33)
- Support conversion of `kong consumer` basic data to `APISIX consumer` data (including: ID, username, custom_id, plugins.keyauth_credentials, plugins.basicauth_credentials, plugins.hmacauth_credentials, plugins.jwt_secrets) [#32](https://github.com/api7/kong-to-apisix/pull/32)
- Support conversion of `kong plugin basic` data to `APISIX plugin` data (including: key-auth, rate-limiting, proxy-cache) [#32](https://github.com/api7/kong-to-apisix/pull/32)
- Support conversion of `kong global_plugin` basic data to `APISIX global rule` data (including: key-auth, rate-limiting, proxy-cache) [#31](https://github.com/api7/kong-to-apisix/pull/31)

### Test
- Improve `kong deck` and `kong config` test framework [#37](https://github.com/api7/kong-to-apisix/pull/37)
- Supplement `e2e` and `unit` test case coverage [#37](https://github.com/api7/kong-to-apisix/pull/37) [#38](https://github.com/api7/kong-to-apisix/pull/38) [#39](https://github.com/api7/kong-to-apisix/pull/39) [#40](https://github.com/api7/kong-to-apisix/pull/40)
