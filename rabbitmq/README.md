Getting started with [Rabbitmq](https://www.rabbitmq.com/) with Golang.

Thanks to That Devops Guy for the inspirational tutorials [link1](https://www.youtube.com/watch?v=hfUIWe1tK8E&t=1036s) and [link2](https://www.youtube.com/watch?v=FzqjtU2x6YA). Haven't watch yet his 3rd video about clustering with kubernetes.

After watching some videos I had the courage to go through Rabbitmq getting started [tutorials](https://www.rabbitmq.com/getstarted.html) with Goland

For the HAProxy bit I have to give credit to [this](http://throughaglass.io/technology/RabbitMQ-cluster-with-Docker-and-Docker-Compose.html) blog
and [this](https://github.com/pardahlman/docker-rabbitmq-cluster/blob/master/haproxy.cfg) repo

## UPDATES 7/24/2022
- using new docker image tag for haproxy
- cluster-haproxy docker composer now uses the .erlang.cookie file instead of the ENV variable

## Forgotten acknowledgment

`rabbitmqctl list_queues name messages_ready messages_unacknowledged`

## erlang cookie file
[docs](https://www.rabbitmq.com/clustering.html#erlang-cookie)
 The file must be only accessible to the owner (e.g. have UNIX permissions of 600 or similar). Every cluster node must have the same cookie.
