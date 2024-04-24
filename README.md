# Userinfo System
This is a simple Go (Golang) backend system example that users can register accounts, login, check and update their profiles. This project contains these points that you can simply refer to:
#### gin, microservice (protobuf, grpc, go-micro, etcd), JWT, dependency injection (wire), project structure, docker, Mysql read/write split and master/slave replication, redis cache, nginx reverse proxy and load balancer, benchmark test (wrk, lua), customized logger.

For more information about the project design, please refer to [the system design doc](User Infomation System Design.pdf).

## Deployment
We can run this project on the local environment by docker, assuming you already have docker on you environment. To run it, we need to deploy all components, including frontend, nginx, userapi, userinfo, mysql, redis, and etcd.
### Clone
```shell
# move into your working dir.
cd your-workspace

# clone this project.
git clone git@github.com:lgangkai/userinfo-system.git 

# move into the root dir of this project, and keep in this dir of the whole deployment.
cd userinfo-system
```
### Docker network
Because we will run this project on a single machine, which is your own PC, we need to use docker custom networks to connect associated containers and isolate different clusters. The network framework is indicated as follows:  
```shell
# create docker networks.
docker network create --driver bridge userapi 
docker network create --driver bridge userinfo 
# specify an unoccupied subnet for redis cluster, because redis cluster can't config container name as ip address.
docker network create --driver bridge redis-cluster --subnet 172.23.0.0/16
docker network create --driver bridge mysql-cluster 
docker network create --driver bridge etcd-cluster 
```
#### Let's deploy independent middlewares at first, including MySql, redis, and etcd.
### MySql
We will build a mysql cluster with one master node and one slave node to realize read/write split and master/slave replication. 
```shell
# 1. run two mysql:8.0 containers and set one as master and one as slave.
# -v is necessary because we don't want our data being deleted if containers are deleted.
docker run --name mysql-master --net mysql-cluster -d -p 13306:3306 \
 -e MYSQL_ROOT_PASSWORD=qwer1234 --privileged=true -v ./data/mysql/master/log:/var/log/mysql \
 -v ./data/mysql/master/data:/var/lib/mysql -v ./data/mysql/master/conf:/etc/mysql/conf.d mysql:8.0
docker run --name mysql-slave --net mysql-cluster -d -p 23306:3306 \
 -e MYSQL_ROOT_PASSWORD=qwer1234 --privileged=true -v ./data/mysql/slave/log:/var/log/mysql \
 -v ./data/mysql/slave/data:/var/lib/mysql -v ./data/mysql/slave/conf:/etc/mysql/conf.d mysql:8.0

# 2. create database userinfo on both nodes.
docker exec -it mysql-master /bin/bash
mysql -uroot -p
qwer1234
CREATE DATABASE userinfo;
# do same operations for mysql-slave node.
```
To realize master/slave replication, we need to add configs to my.cnf file. 
```shell
docker cp mysql-master:/etc/my.cnf ./
```
Open my.cnf and add this code block under [mysqld] label:
```shell
[mysqld]
# add config here.
log-bin=master-bin
binlog-format=ROW
server-id=1
binlog-do-db=userinfo
#...
# remove leading #
default-authentication-plugin=mysql_native_password
```
Copy file back and restart mysql.
```shell
docker cp ./my.cnf mysql-master:/etc/
docker restart mysql-master
```
#### Do the same operations that change my.cnf to slave node too. The only difference is that the added configs are:
```shell
[mysqld]
# add config here.
log-bin=master-bin
binlog-format=ROW
server-id=2
```
Then we need to grant permissions on master node under mysql service.
```shell
install plugin validate_password soname 'validate_password.so'; # Since mysql8.x, it is required.
set global validate_password_policy=0;
set global validate_password_length=1;
CREATE USER 'repl'@'%' IDENTIFIED BY 'qwer1234';
ALTER USER 'repl'@'%' IDENTIFIED WITH mysql_native_password BY 'qwer1234';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';
FLUSH PRIVILEGES;
SHOW master status;
+-------------------+----------+--------------+------------------+-------------------+
| File              | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+-------------------+----------+--------------+------------------+-------------------+
| master-bin.000001 |     1109 | userinfo     |                  |                   |
+-------------------+----------+--------------+------------------+-------------------+
```
And config on slave node under mysql service.
```shell
CHANGE MASTER TO
MASTER_HOST='mysql-master',
MASTER_PORT=3306,
MASTER_USER='repl',
MASTER_PASSWORD='qwer1234',
MASTER_LOG_FILE='master-bin.000001', # the file name of master status
MASTER_LOG_POS=1109;                 # the Position of master status
start slave;
show slave status\G
*************************** 1. row ***************************
               Slave_IO_State: Waiting for source to send event
                  Master_Host: mysql-master
                  Master_User: repl
                  Master_Port: 3306
                Connect_Retry: 60
              Master_Log_File: master-bin.000001
          Read_Master_Log_Pos: 1109
               Relay_Log_File: 8ff5047ed91a-relay-bin.000002
                Relay_Log_Pos: 327
        Relay_Master_Log_File: master-bin.000001
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
                           ......
```
Slave_IO_Running and Slave_SQL_Running all Yes means slave service is running successfully. Now you can create test tables and insert some data on master node and check whether there is replication on slave node. 

Let's initialize the database for this project on master node:
```shell
# Run script file defined at ./script/init-db.sql
docker cp ./script/init-db.sql mysql-master:/
docker exec -it mysql-master /bin/bash
mysql -uroot -p
qwer1234
source init-db.sql
```


### Redis
Similarly, for redis, we will deploy a redis cluster with typically 6-node structure containing 3 master nodes and 3 slave nodes, in order to reach a high availability and load balance. 
```shell
# 1. Use shell script to write 6 redis config files.
for node in $(seq 0 5); \
do \
mkdir -p ./conf/redis
touch ./conf/redis/redis-node${node}.conf
cat << EOF >./conf/redis/redis-node${node}.conf
port 6379
bind 0.0.0.0
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 5000
cluster-announce-ip 172.23.0.1${node}
cluster-announce-port 6379
cluster-announce-bus-port 16379
appendonly yes
EOF
done

# 2. run 6 redis containers based on config files.
docker run -d -p 6371:6379 -p 16371:16379 --name redis-node0 --net redis-cluster --ip 172.23.0.10 -v ./data/redis/node0:/data \
 -v ./conf/redis/redis-node0.conf:/etc/redis/redis.conf redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf
docker run -d -p 6372:6379 -p 16372:16379 --name redis-node1 --net redis-cluster --ip 172.23.0.11 -v ./data/redis/node1:/data \
 -v ./conf/redis/redis-node1.conf:/etc/redis/redis.conf redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf
docker run -d -p 6373:6379 -p 16373:16379 --name redis-node2 --net redis-cluster --ip 172.23.0.12 -v ./data/redis/node2:/data \
 -v ./conf/redis/redis-node2.conf:/etc/redis/redis.conf redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf
docker run -d -p 6374:6379 -p 16374:16379 --name redis-node3 --net redis-cluster --ip 172.23.0.13 -v ./data/redis/node3:/data \
 -v ./conf/redis/redis-node3.conf:/etc/redis/redis.conf redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf
docker run -d -p 6375:6379 -p 16375:16379 --name redis-node4 --net redis-cluster --ip 172.23.0.14 -v ./data/redis/node4:/data \
 -v ./conf/redis/redis-node4.conf:/etc/redis/redis.conf redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf
docker run -d -p 6376:6379 -p 16376:16379 --name redis-node5 --net redis-cluster --ip 172.23.0.15 -v ./data/redis/node5:/data \
 -v ./conf/redis/redis-node5.conf:/etc/redis/redis.conf redis:5.0.9-alpine3.11 redis-server /etc/redis/redis.conf

# 3. create redis cluster.
docker exec -it redis-node0 /bin/sh 
redis-cli --cluster create 172.23.0.10:6379 172.23.0.11:6379 172.23.0.12:6379 172.23.0.13:6379 \
 172.23.0.14:6379 172.23.0.15:6379 --cluster-replicas 1
```
Now redis cluster is successfully deployed. Let's check whether it works.
```shell
redis-cli -c
set foo bar
# -> Redirected to slot [12182] located at 172.23.0.12:6379
# we can see our set request is redirected to node2
# let's test what will happen if node2 dead.
docker stop redis-node2 
docker exec -it redis-node0 /bin/sh 
redis-cli -c
get foo
# -> Redirected to slot [12182] located at 172.23.0.13:6379
# "bar"
# now the request is redirected to node3, which is the slave node of node2 before.
# if you restart node2, you will find node2 now becoming slave node.
```
### etcd
We use etcd to register and find our microservice. Here we deploy the etcd cluster with three nodes.
```shell
# run an etcd container using the config file
docker run -d -p 2379:2379 -p 2380:2380 --net etcd-cluster -v ./data/etcd/node0:/etcd-data -v ./conf/etcd:/etcdconf \
--name etcd0 gcr.io/etcd-development/etcd:v3.5.13 /usr/local/bin/etcd --config-file=/etcdconf/etcd0.yaml
docker run -d -p 12379:2379 -p 12380:2380 --net etcd-cluster -v ./data/etcd/node1:/etcd-data -v ./conf/etcd:/etcdconf \
--name etcd1 gcr.io/etcd-development/etcd:v3.5.13 /usr/local/bin/etcd --config-file=/etcdconf/etcd1.yaml
docker run -d -p 22379:2379 -p 22380:2380 --net etcd-cluster -v ./data/etcd/ndoe2:/etcd-data -v ./conf/etcd:/etcdconf \
--name etcd2 gcr.io/etcd-development/etcd:v3.5.13 /usr/local/bin/etcd --config-file=/etcdconf/etcd2.yaml
```
To check whether it runs successfully, run the following command.
```shell
# check the node health
docker exec etcd0 /usr/local/bin/etcdctl endpoint health --cluster -w table
```
And you may see the response like this:
```shell
+-------------------+--------+------------+-------+
|     ENDPOINT      | HEALTH |    TOOK    | ERROR |
+-------------------+--------+------------+-------+
| http://etcd2:2379 |   true | 1.477459ms |       |
| http://etcd0:2379 |   true | 1.481584ms |       |
| http://etcd1:2379 |   true |  1.51775ms |       |
+-------------------+--------+------------+-------+
```
And run these commands to check connection between nodes:
```shell
# put a key-value pair to node0/
docker exec etcd0 /usr/local/bin/etcdctl put foo bar # OK
# get value by key from node0
docker exec etcd0 /usr/local/bin/etcdctl get foo # foo bar
# get value by key from node1
docker exec etcd1 /usr/local/bin/etcdctl get foo # foo bar
# get value by key from node2
docker exec etcd2 /usr/local/bin/etcdctl get foo # foo bar
```
### userinfo
Let's deploy 2 instances for userinfo microservice.
```shell
# 1. build the image of userinfo from dockerfile.
docker build -t userinfo -f ./userinfo/Dockerfile .

# 2. create containers.
docker create --name userinfo-node0 --net userinfo -p 8081:8081 userinfo  
docker create --name userinfo-node1 --net userinfo -p 18081:8081 userinfo 
 
# 3. connect containers to dependent networks.
docker network connect etcd-cluster userinfo-node0
docker network connect etcd-cluster userinfo-node1
docker network connect mysql-cluster userinfo-node0
docker network connect mysql-cluster userinfo-node1
docker network connect redis-cluster userinfo-node0
docker network connect redis-cluster userinfo-node1

# 4. run containers. 
docker start userinfo-node0
docker start userinfo-node1

# 5. check whether services are registered successfully. 
docker exec etcd0 /usr/local/bin/etcdctl get / --prefix --keys-only=true
# The response should look like this:
#/micro/registry/api.lgk.com.userinfo/api.lgk.com.userinfo-d57a4941-3ad3-4ebe-a213-eb7ff957b4d5
#/micro/registry/api.lgk.com.userinfo/api.lgk.com.userinfo-d592449a-ee30-4758-b1f8-a079a008153b

```
### userapi
Let's deploy 2 instances for api service too.
```shell
# 1. build the image of userapi from dockerfile.
docker build -t userapi -f ./userapi/Dockerfile .

# 2. create containers.
docker create --name userapi-node0 --net userapi -p 8080:8080 userapi
docker create --name userapi-node1 --net userapi -p 18080:8080 userapi

# 3. connect containers to dependent networks.
docker network connect userinfo userapi-node0
docker network connect userinfo userapi-node1
docker network connect etcd-cluster userapi-node0
docker network connect etcd-cluster userapi-node1

# 4. run containers. 
docker start userapi-node0
docker start userapi-node1
```
### Frontend & Nginx
Here we use nginx to deploy the frontend webpages. We also use nginx as reverse proxy and load balancer for userapi nodes.
* Run nginx.
```shell
# run an nginx container, and specify the frontend dir as html source.
docker run --name=nginx -d -p 80:80 -v ./frontend:/usr/share/nginx/html nginx 

# connect containers to userapi network.
docker network connect userapi nginx
```
* Config nginx to perform dynamic static resource separation, and load balance for the api server.
```shell
 # move into the nginx container.
 docker exec -it nginx /bin/bash
 
 # (optional) if there is no vim on the container, install vim.
 apt-get update
 apt-get install vim
 
 # open the nginx config file.
 vim /etc/nginx/conf.d/default.conf
```
&emsp;&emsp;&emsp;&emsp;add these code to default.conf file.
```shell
# Add this code block at the root level of the file.
# config 2 api server for load balance, with weighted round-robin method.
upstream userapi {
  server userapi-node0:8080 weight=1;
  server userapi-node1:8080 weight=2;
}

server {
  # ...
  #access_log /var/log/nginx/host.access.log main
  
  # add this code block in the server block.
  # if there already exists the same code block, then don't add it.
  # for static resource(html)
  location / {
    root /usr/share/nginx/html;
    index index.html index.htm;
  }
  
  # for dynamic contents(api server)
  location /api {
    proxy_pass http://userapi;
  }
  
  # ...
}
```
&emsp;&emsp;&emsp;&emsp;reload nginx config.
```shell
nginx -s reload
```
docker build -t userinfo -f ./userinfo/Dockerfile .
docker build -t userapi .

docker run --name userinfo --net etcd-cluster userinfo-node0


docker exec etcd0 /usr/local/bin/etcdctl get / --prefix --keys-only=true

dock

tips：the address must be set to 0.0.0.0 or you can not access container from host browser

go install github.com/asim/go-micro/cmd/protoc-gen-micro/v3@latest

ab -n 10000 -c 100 -C 'access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM5NDIwNzUsInVzZXJfaWQiOjEsImVtYWlsIjoiMTIzQGdtYWlsLmNvbSJ9.NY7_6E3nt_P2X9fTNofscnUPi9nwVeC6WWqZiLiAaz0; Path=/; HttpOnly; Expires=Wed, 24 Apr 2024 07:01:15 GMT;' http://localhost/api/user/profile