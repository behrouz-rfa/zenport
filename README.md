# zenport

### Golang [zenport]() is a simple Even Driven Architector

#### 🙌👨‍💻🚀 Docker-compose files easy lunch:

    docker-compose.yml
    sudo docker compose up 

```
curl -X 'POST' \
'http://localhost:8080/api/times' \
-H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"ask": "What time is it?"
}'

```

#### Recommendation for local development most comfortable usage:

    make run   // run the application
    make race // run the application with  race flag
    make integration-test // run integration test

## Using

supports the following options:

```
support .env file on docker dircotry
PG_CONN= //postgres url connection
NATS_URL=  nats url
RB_URL= rabbitmq url
RABBITMQC_ISENABLE= true by default is true // we have alost nats support by set to false you we abl;e to work with nats
```





