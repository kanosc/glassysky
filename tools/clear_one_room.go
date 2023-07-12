package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "929319", // no password set
	//Password: "myredis6379", // no password set
	DB:       0,
})
var ctx = context.Background()

func main() {

	rooms, err := redisClient.LRange(ctx, "chat:rooms", 0, -1).Result()
	log.Println("chat rooms ", rooms)
	if err != nil {
		log.Println(err.Error())
	}
	_, _ = redisClient.Expire(ctx, "chat:roomname:"+"家"+":messages", 2*time.Second).Result()

}
