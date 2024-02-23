package redis_pagination

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

type User struct {
	ID    int    `redis:"id" json:"id"`
	Name  string `redis:"name" json:"name"`
	Phone string `redis:"phone" json:"phone"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func TestHashDb_Store(t *testing.T) {
	db := NewHashDB(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}).SetHashTable("test")

	var err error
	count := 50
	for i := 0; i < count; i++ {
		user := &User{
			ID:    i,
			Name:  fmt.Sprintf("%s%d", "user", i),
			Phone: "13888888888",
		}
		err = db.SetHashField(fmt.Sprint("field", i)).Store(user)
		if err != nil {
			t.Fatal("Store() 插入失败:", err)
		}
	}
}

func TestHashDb_Paginate(t *testing.T) {
	db := NewHashDB(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}).SetHashTable("test")
	count := 10
	paginate, err := db.Paginate(1, count, "desc")
	if err != nil {
		t.Fatal("Paginate() 获取数据出错:", err)
	}
	if count != len(paginate.Rows) {
		t.Fatal("Paginate() 获取数据条数与期望不符")
	}
}
