# glassysky
this is a web project developed with golang
to start the server,use command:
'./restart' in production enviroment or './restart -m debug' for local enviroment(the default port is 9090)
if you start the project in local,you can visit the homepage by URL: http://lcoalhost:9090/

configure of redis: redis is used to storage the messages of chat root, it should work in a docker container
step1
docker search redis
docker pull redis
step2
edit the config file of redis, the file is put to config/redis path which is called redis.conf
step3
run the container with command like this:
sudo docker run --restart=always --log-opt max-size=100m --log-opt max-file=2 -p 6379:6379 --name myredis -v /home/ubuntu/redis/myredis/myredis.conf:/etc/redis/redis.conf -v /home/redis/myredis/data:/data -d redis redis-server /etc/redis/redis.conf  --appendonly yes  --requirepass 123456
you should change the requirepass to your own password and make the directories if they are not exist on your machine
step4
test the redis by client, use the following commands
docker exec -it myredis redis-cli
auth your-password
config get requirepass
ping
docker logs  myredis

