package main

import (
	"encoding/json"
	"fmt"
	rp "github.com/jerry-hms/redis-pagination"
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

func main() {
	// 存储调用演示
	testStore()
	// 分页数据调用演示
	testPaginate()
}

// 存储调用演示
func testStore() {
	var err error
	for i := 0; i < 50; i++ {
		user := User{
			ID:    i,
			Name:  fmt.Sprintf("%s%d", "user", i),
			Phone: "13888888888",
		}
		data := rp.NewData(&user, &rp.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
			Key:      "user_list:1",
			Field:    fmt.Sprintf("%s%d", "test", i),
		})
		err = data.Store()
	}

	if err != nil {
		fmt.Println("存储失败")
		return
	}
	fmt.Println("存储成功")
}

func testPaginate() {
	var user User
	data := rp.NewData(&user, &rp.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Key:      "user_list:1",
		Page:     1,
		PageSize: 10,
	})
	paginate, err := data.Paginate("desc")
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}

	fmt.Printf("paginate:%+v\n", paginate)

	//for _, item := range paginate.Rows {
	//	_ = json.Unmarshal([]byte(item), &user)
	//	fmt.Println("user", user)
	//}
}
