# redis-pagination
redis-pagination 是基于redis实现的数据列表存储，支持分页、排序功能

### 安装
```shell
go get github.com/jerry-hms/redis-pagination
```

### 使用示例
#### 调用存储
```go
package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jerry-hms/redis-pagination"
)

type User struct {
	ID    int    `redis:"id" json:"id"`
	Name  string `redis:"name" json:"name"`
	Age   int `redis:"age" json:"age"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func main() {
	db := NewHashDB(&redis.Options{
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
```
#### 获取分页数据
```go
func main() {
	db := NewHashDB(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}).SetHashTable("your_hash_key")

	result, err = db.Paginate(1, 10, "desc")
	if err != nil {
		// ...
	}
	fmt.Printf("result:%+v\n", result)
}
```