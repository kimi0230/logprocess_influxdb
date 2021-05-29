# logprocess_influxdb

Read and write requsets into influxdb Concurrently.
from imooc [https://www.imooc.com/learn/982](https://www.imooc.com/learn/982) and [imooc_logprocess](https://github.com/itsmikej/imooc_logprocess)

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
 --header "Authorization: Token I7XePMLi3cx-j4PjnpkMC1_jImyEikWCSL9ar7hNC4Ji4IEucISYgULyl2AORJdPaTrf2PixpZz2euSAzQLfCw==" \
 --data-binary "kk,machine_id=1,region=tw value=0.5"
```
## grafana
### install & run
``` shell
docker run -d -p 3000:3000 grafana/grafana
```
default account/passwword is `admin/admin`


## docker compose
``` shell
# 背景執行
docker-compose up -d
# 停止執行
docker-compose stop
# 停用再移除  
docker-compose rm -s
```


![image](https://github.com/kimi0230/logprocess_influxdb/blob/master/screenshot/data_explore.png)

# Reference
* https://hub.docker.com/_/influxdb 
* https://github.com/influxdata/influxdb
* https://github.com/influxdata/influxdb-client-go
* https://docs.influxdata.com/influxdb/v2.0/
* https://db-engines.com/en/ranking

