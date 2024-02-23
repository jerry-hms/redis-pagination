package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	redis_pagination "github.com/jerry-hms/redis-pagination"
)

type User struct {
	ID   int    `redis:"id" json:"id"`
	Name string `redis:"name" json:"name"`
	Age  int    `redis:"age" json:"age"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func main() {
	db := redis_pagination.NewHashDB(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}).SetHashTable("your_hash_key")

	var err error
	user := &User{
		ID:   1,
		Name: "jerry",
		Age:  18,
	}
	err = db.SetHashField("your_hash_field").Store(user)
	if err != nil {
		// fail
	}
	// success
}
