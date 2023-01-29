# zenport

### Golang [zenport]() is a simple Even Driven Architector

#### 👨‍💻 Full list what has been used:

* [fasthttp](github.com/valyala/fasthttp/) - fasthttp
* [color](github.com/fatih/color/) - fatih color
* [Console Std Out](github.com/apoorvam/gotermina) - gotermina

#### Recommendation for local development most comfortable usage:

    make run   // run the application
    make race // run the application with  race flag
    make integration-test // run integration test

#### 🙌👨‍💻🚀 Docker-compose files:

    docker-compose.yml

## Using

supports the following options:

```
support .env file on docker dircotry
PG_CONN= //postgres url connection
NATS_URL=  nats url
RB_URL= rabbitmq url
RABBITMQC_ISENABLE= true by default is true // we have alost nats support by set to false you we abl;e to work with nats
```





