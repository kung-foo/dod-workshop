Demo 04 Everything
==================

Init
----

1. `cd demo-04-full`
2. `docker-compose build`
3. `docker-compose rm -f && docker-compose up`

Consul Template
---------------

1. Show the current config:

        $ docker exec -it demo04full_proxy_1 cat /root/config.toml
        Master = "redis-master"
        Slaves = [
            "192.168.200.8",
        ]
2. Scale up: `docker-compose scale slave=2`
3. Show the current config:

        $ docker exec -it demo04full_proxy_1 cat /root/config.toml
        Master = "redis-master"
        Slaves = [
            "192.168.200.8","192.168.200.9",
        ]

4. For fun, try `kill -s STOP`ing a redis slave and looking for a change in the proxy's config.toml.
5. Scale back to one: `docker-compose scale slave=1`

Add Some Data
-------------

1. Find address/port: `docker-compose port proxy 8888`
2. `POST` some data (returns length of key):
        $ curl -X POST http://localhost:8888/db/hello -d 'world'
        5
3. Verify:

        $ curl -i http://localhost:8888/db/hello
        HTTP/1.1 200 OK
        X-Rrproxy-Server: 192.168.200.4:6379
        Date: Tue, 01 Dec 2015 10:40:32 GMT
        Content-Length: 5
        Content-Type: text/plain; charset=utf-8

        world

Benchmark
---------

1. `ab -k -n 10000 -c 8 http://localhost:8888/db/hello`
2. My laptop: `Requests per second: 808.12 [#/sec] (mean)`


Scale
-----

1. Add a new redis slave: `docker-compse scale slave=2`
2. Verify in Consul UI that the new container is healthy
3. Look for log message indicating that the new slave is added:

        proxy_1     | time="2015-12-01T10:33:16Z" level=info msg="Connected 1 master and 1 slaves"
4. Use curl to verify that the new slave is in the rotation:

        $ curl -si http://localhost:8888/db/hello | grep X-Rrproxy-Server
        X-Rrproxy-Server: 192.168.200.7:6379
        $ curl -si http://localhost:8888/db/hello | grep X-Rrproxy-Server
        X-Rrproxy-Server: 192.168.200.4:6379

5. Benchmark again: `ab -k -n 10000 -c 8 http://localhost:8888/db/hello`
6. My laptop: `Requests per second: 1482.58 [#/sec] (mean)`
7. With 4 slaves: `Requests per second: 2395.43 [#/sec] (mean)`
