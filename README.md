# logprocess_influxdb 2.0

Read and write requsets into [influxdb](https://www.influxdata.com/) Concurrently.
Logs in Explore by [grafana](https://grafana.com/).

## influxdb
### install

``` shell
docker pull influxdb:2.0.6
```

### run
``` shell
docker run -p 8086:8086 \
-e DOCKER_INFLUXDB_INIT_USERNAME=kimiuser \
-e DOCKER_INFLUXDB_INIT_PASSWORD=kimipassword \
influxdb:2.0.6
```


### test
``` shell
curl -i -XPOST 'http://localhost:8086/api/v2/write?org=kimiORG&bucket=kk&precision=s'  -u kimiuser:kimipassword \
 --header "Authorization: Token rOdFRpyoGg9aGQKt6GhUlBHgkmFOCX5HNIirOmBU3fCFbiwAY4jviMtbxtBBJ9cry_OGiEjieHrEPnfSeRO1mw==" \
 --data-binary "kk,machine_id=1,region=tw value=0.5"
```
## grafana
### install & run
``` shell
docker run -d -p 3000:3000 grafana/grafana
```
default account/passwword is `admin/admin`


## Demo: docker compose
``` shell
# Background execution (docker run)
docker-compose up -d

# (docker ps)
docker-compose ps 

# show logs (docker logs)
docker-compose logs

# start (docker start)
docker-compose start

# Stop execution (docker stop)
docker-compose stop

# remove (docker rm)
docker-compose down

# Stop execution, then remove  
docker-compose rm -s
```

## System Info API
``` shell
curl 127.0.0.1:9193/monitor
```

## Screenshot
### influxDB
![influxDB](https://github.com/kimi0230/logprocess_influxdb/blob/master/screenshot/data_explore.png)
### grafana
![grafana](https://github.com/kimi0230/logprocess_influxdb/blob/master/screenshot/grafana.png)

# Reference
* https://hub.docker.com/_/influxdb 
* https://github.com/influxdata/influxdb
* https://github.com/influxdata/influxdb-client-go
* https://docs.influxdata.com/influxdb/v2.0/
* https://db-engines.com/en/ranking
* [https://www.imooc.com/learn/982](https://www.imooc.com/learn/982)
* [imooc_logprocess](https://github.com/itsmikej/imooc_logprocess)
