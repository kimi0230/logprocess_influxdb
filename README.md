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
curl -i -XPOST 'http://localhost:8086/api/v2/write?org=kimiORG&bucket=kk&precision='  -u kimiuser:kimipassword \
 --header "Authorization: Token dOrI2xnBcY1A62JZTmcPEoTf30K9I5iro10fwHvSU6xDJK8aXFo_QncuAlxTGruIsQWeu9bq2WEylszu4lTP4A==" \
 --data-binary "kk,machine_id=1,region=tw value=0.5"
```

# Reference
* https://hub.docker.com/_/influxdb 
* https://github.com/influxdata/influxdb
* https://github.com/influxdata/influxdb-client-go
* https://docs.influxdata.com/influxdb/v2.0/
* https://db-engines.com/en/ranking

