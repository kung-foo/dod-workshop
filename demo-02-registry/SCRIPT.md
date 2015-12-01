Demo 02 Registrator
===================

Init
----

1. `cd demo-02-registry`
2. `docker-compose build`
3. `docker-compose rm -f && docker-compose up`

Registrator
-----------

1. Open a new terminal and start tailing the registrator log: `docker-compose logs registrator`
2. Look for:
  ```
Using consulkv adapter: consulkv://consul:8500/registry
Listening for Docker events ...
Syncing services on 4 containers
added: 679e8e179961 registrator:demo02registry_web_1:5000
added: 0f847613e3e1 registrator:demo02registry_redis_1:6379
ignored: 1293eb7680d8 no published ports
ignored: c77c7bc6901a port 53 not published on host
  ```
3. Scale redis container: `docker-compose scale redis=2`
4. Look for new registrator log line:
  `added: 7ad5747baf3c registrator:demo02registry_redis_2:6379`

Registrator Backend
-------------------

1. Use curl to view registrator data that is stored in Consul KV: `curl -s "http://localhost:8500/v1/kv/registry?recurse&pretty"`
2. Note that the values are base64 encoded
