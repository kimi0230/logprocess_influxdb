# imooc_logprocess_practise

from imooc [https://www.imooc.com/learn/982](https://www.imooc.com/learn/982) and [imooc_logprocess](https://github.com/itsmikej/imooc_logprocess)

Example read and write requsets into influxdb Concurrently.

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

# Reference
* https://hub.docker.com/_/influxdb 
* https://github.com/influxdata/influxdb1-client/


curl --request POST 'http://127.0.0.1:8086/write?db=kk&precision=s' \
 --header "Authorization: Token iLDSvb1q2I5C2e5eyElYhM37n5Y2TQbwJhQyVNgaK3tSDWbJMM_m1kUbwoRvVsiAt5ytJdiWuHTCrd8Cf08Q5A==" \
 --data-urlencode "Method=POST,Scheme=ss,Status=sst" \
  --data-urlencode "bucket=example-bucket"
  
   --data-binary 'kk,Method=POST,Scheme=ss,Status=sst' --header "Authorization: Token iLDSvb1q2I5C2e5eyElYhM37n5Y2TQbwJhQyVNgaK3tSDWbJMM_m1kUbwoRvVsiAt5ytJdiWuHTCrd8Cf08Q5A=="