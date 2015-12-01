Demo 01 Basics
==============

1. `cd demo-01-basics`
2. `docker-compose build`
3. `docker-compose rm -f && docker-compose up`
4. Open a new terminal. `docker-compose ps web` or `docker-compose port web 5000` and look for mapped port.
5. Open a browser to http://IP_ADDR_OF_DOCKER_BRIDGE:PORT/ (henceforth `$APP`) and you should see "Hello World!"

    _Linux only: you can also read the log output of `docker-compose up` and see the container's IP address in a message like:_

    `web_1   |  * Running on http://192.168.200.3:5000/`
6. Access `http://$APP/env` to see a filtered list of environmental variables exposed in the container.
7. Access `http://$APP/hosts` to see a filtered list of hostsnames that can be resolved in the container.
8. Add a second redis container `docker-compose scale redis=2`.
9. Verify that none of the values changed in the app container.
10. Add a second app container `docker-compose scale web=2`.
11. Verify that none of the values changed in the first app container.
12. See the new values in the second app container (find the address using the same method as in step 4)

New In 1.9
----------

1. Kill any existing containers: `docker-compose kill` or `Ctrl-C`
2. `docker-compose --x-networking up`
3. Access `http://$APP/hosts` and see a list of hostnames. For example:
  ```
172.18.0.2	web
172.18.0.3	demo01basics_redis_1
172.18.0.3	demo01basics_redis_1.demo01basics
  ```
4. Open a new terminal. `docker-compose --x-networking scale redis=2`
5. Access `http://$APP/hosts` and see an updated list of hostnames. For example:
  ```
172.18.0.2	web
172.18.0.3	demo01basics_redis_1
172.18.0.3	demo01basics_redis_1.demo01basics
172.18.0.4	demo01basics_web_2
172.18.0.4	demo01basics_web_2.demo01basics
  ```
6. Run `docker network ls` to see a list of network managed by Docker:
```
NETWORK ID          NAME                DRIVER
4213bab7fabc        demo01basics        bridge
3990af1a8ec7        bridge              bridge
63454c509063        none                null
fa4d3f714ea7        host                host
```

  Note the `demo01basics` network that was created dynamically by `docker-compose`.

  You can also get more bridge details:

  ```
$ ifconfig br-4213bab7fabc
br-4213bab7fabc Link encap:Ethernet  HWaddr 02:42:ce:6b:81:ee  
        inet addr:172.18.0.1  Bcast:0.0.0.0  Mask:255.255.0.0
        inet6 addr: fe80::42:ceff:fe6b:81ee/64 Scope:Link
        UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
        RX packets:87 errors:0 dropped:0 overruns:0 frame:0
        TX packets:171 errors:0 dropped:0 overruns:0 carrier:0
        collisions:0 txqueuelen:0
        RX bytes:6586 (6.5 KB)  TX bytes:26347 (26.3 KB)

$ brctl show br-4213bab7fabc
bridge name         bridge id           STP enabled     interfaces
br-4213bab7fabc     8000.0242ce6b81ee   no              veth312d030
                                                        veth49fdb0c
                                                        veth8b8feb5
  ```
