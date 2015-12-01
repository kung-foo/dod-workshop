Demo 03 Consul
==============

Init
----

1. `cd demo-03-consul`
2. `docker-compose build`
3. `docker-compose rm -f && docker-compose up`
4. Consul UI is available on 8500

Consul Health
-------------

1. Create a new shell in the redis container: `docker exec -it demo03consul_redis_1 sh`
2. Find the pid of redis: `ps | grep redis-server`
3. Stop redis: `kill -s STOP $PID`
4. Go back to the Consul UI
5. List processes and look for a health check:
  ```
116 redis      0:01 redis-server /etc/redis.conf
528 root       0:00 sh
613 consul     0:00 /bin/sh -c redis-check
614 consul     0:00 sh /usr/bin/redis-check
615 consul     0:00 redis-cli ping
  ```
6. Start redis: `kill -s CONT $PID`
7. Verify in the UI and process list that the service is now healthly.

Consul DNS
----------

1. Create a new shell in the consul container: `docker exec -it demo03consul_consul_1 sh`
2. Query the local DNS for the redis service: `dig @127.0.0.1 -p 8600 redis.service.dc1.consul. ANY`
3. Scale up the redis container: `docker-compose scale redis=2`
4. Query the local DNS and see the new container. Also notice the order of entries dig returns.
