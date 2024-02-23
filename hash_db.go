package redis_pagination

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"sync"
)

var RedisIns *redis.Client
var RedisOnce sync.Once

func NewHashDB(opt *redis.Options) *HashDb {
	return &HashDb{
		rdb: connRedis(opt),
	}
}

func connRedis(opt *redis.Options) *redis.Client {
	RedisOnce.Do(func() {
		RedisIns = redis.NewClient(&redis.Options{
			Addr:     opt.Addr,
			Password: opt.Password,
			DB:       opt.DB,
		})
	})
	return RedisIns
}

type HashDb struct {
	rdb   *redis.Client
	key   string
	field string
}

func (d *HashDb) SetHashTable(key string) *HashDb {
	d.key = key
	return d
}

func (d *HashDb) GetHashTable() (string, error) {
	if d.key == "" {
		return "", errors.New("the hash table is not set, please call the SetHashTable method to to set it")
	}
	return d.key, nil
}

func (d *HashDb) SetHashField(field string) *HashDb {
	d.field = field
	return d
}

func (d *HashDb) GetHashField() (string, error) {
	if d.field == "" {
		return "", errors.New("the storage field is not set, please call the SetHashField method to set it")
	}
	return d.field, nil
}

// Store 存储数据
func (d *HashDb) Store(data interface{}) error {
	var err error
	jsonStr, _ := json.Marshal(data)

	hashTable, err := d.GetHashTable()
	if err != nil {
		return err
	}
	field, err := d.GetHashField()
	if err != nil {
		return err
	}

	storeScript := redis.NewScript(`
local hashTable = KEYS[1]
local field = ARGV[1]
local jsonStr = ARGV[2]

-- 检查字段是否存在
local exists = redis.call('hexists', hashTable, field)
if exists == 0 then
    -- 如果字段不存在，设置哈希表的值，并将字段添加到有序集合中
    local length = redis.call('hlen', hashTable)
    redis.call('zadd', hashTable .. ':sort', length, field)
end
-- 存储或更新hash表
redis.call('hset', hashTable, field, jsonStr)
return true
`)
	err = storeScript.Run(context.Background(), d.rdb, []string{hashTable}, field, string(jsonStr)).Err()
	return err
}

// First 查询一条数据
func (d *HashDb) First() (string, error) {
	hashTable, err := d.GetHashTable()
	if err != nil {
		return "", err
	}
	field, err := d.GetHashField()
	if err != nil {
		return "", err
	}
	result, _ := d.rdb.HGet(context.Background(), hashTable, field).Result()
	return result, nil
}

// Paginate 分页获取数据
func (d *HashDb) Paginate(page int, limit int, sort string) (*Pagination, error) {
	var rows, keys []string

	paginate := newPagination(page, limit)
	switch sort {
	case "asc":
		keys, _ = d.rdb.ZRange(context.Background(), d.GetSortKey(), int64(paginate.GetOffset()), int64(paginate.GetEnd())).Result()
	case "desc":
		keys, _ = d.rdb.ZRevRange(context.Background(), d.GetSortKey(), int64(paginate.GetOffset()), int64(paginate.GetEnd())).Result()
	}

	hashTable, err := d.GetHashTable()
	if err != nil {
		return nil, err
	}
	paginate.countTotalPages(d.rdb.HLen(context.Background(), hashTable).Val())
	for _, key := range keys {
		// 从redis中取出会话数据
		res := d.rdb.HGet(context.Background(), hashTable, key).Val()
		if res != "" {
			rows = append(rows, res)
		}
	}
	paginate.Rows = rows
	return paginate, nil
}

// Del 删除数据
func (d *HashDb) Del() error {
	delScript := redis.NewScript(`
local hashTable = KEYS[1]
local field = ARGV[1]
local jsonStr = ARGV[2]

-- 检查字段是否存在
local exists = redis.call('hexists', hashTable, field)
if exists == 1 then
    -- 字段存在执行删除
	redis.call('hdel', hashTable, field)
    redis.call('zrem', hashTable .. ':sort', field)
end
return true
`)
	hashTable, err := d.GetHashTable()
	if err != nil {
		return err
	}
	field, err := d.GetHashField()
	if err != nil {
		return err
	}

	err = delScript.Run(context.Background(), d.rdb, []string{hashTable}, field).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetSortKey 获取排序key
func (d *HashDb) GetSortKey() string {
	hashTable, _ := d.GetHashTable()
	return hashTable + ":sort"
}
