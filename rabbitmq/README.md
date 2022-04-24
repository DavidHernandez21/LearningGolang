Getting started with [Rabbitmq](https://www.rabbitmq.com/) with Golang.

Thanks to That Devops Guy for the inspirational tutorials [link1](https://www.youtube.com/watch?v=hfUIWe1tK8E&t=1036s) and [link2](https://www.youtube.com/watch?v=FzqjtU2x6YA). Haven't watch yet his 3rd video about clustering with kubernetes. Don't forget to enable the rabbitmq_federation plugin on all the nodes and also check out the [Basic Queue Mirroring](https://github.com/marcel-dempers/docker-development-youtube-series/tree/master/messaging/rabbitmq#basic-queue-mirroring) policy 

After watching some videos I had the courage to go through Rabbitmq getting started [tutorials](https://www.rabbitmq.com/getstarted.html) with Goland

For the HAProxy bit I have to give credit to [this](http://throughaglass.io/technology/RabbitMQ-cluster-with-Docker-and-Docker-Compose.html) blog
and [this](https://github.com/pardahlman/docker-rabbitmq-cluster/blob/master/haproxy.cfg) repo

I also recommend [video](https://youtube.com/watch?v=ez9kQEhHsnc)

For the docker-compose up command to work you would need to add a .env file with the value of the ERLANG_COOKIE variable
