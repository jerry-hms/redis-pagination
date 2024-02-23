package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	redis_pagination "github.com/jerry-hms/redis-pagination"
)

func main() {
	db := redis_pagination.NewHashDB(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}).SetHashTable("your_hash_key")

	result, err := db.Paginate(1, 10, "desc")
	if err != nil {
		//...
	}
	fmt.Printf("result:%+v\n", result)
}
