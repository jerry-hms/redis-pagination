package redis_pagination

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"sync"
)

var RedisIns *redis.Client
var RedisOnce sync.Once

func NewData(dto interface{}, opt *Options) *Data {
	return &Data{
		rdb: connRedis(opt),
		opt: opt,
		Dto: dto,
	}
}

func connRedis(opt *Options) *redis.Client {
	RedisOnce.Do(func() {
		RedisIns = redis.NewClient(&redis.Options{
			Addr:     opt.Addr,
			Password: opt.Password,
			DB:       opt.DB,
		})
	})
	return RedisIns
}

type Data struct {
	rdb *redis.Client
	opt *Options
	Dto interface{}
}

// Store 存储数据
func (d *Data) Store() error {
	var err error
	jsonStr, _ := json.Marshal(d.Dto)
	ets := d.rdb.HExists(context.Background(), d.opt.Key, d.opt.Field).Val()

	if err = d.rdb.HSet(context.Background(), d.opt.Key, d.opt.Field, string(jsonStr)).Err(); err != nil {
		return err
	}
	if !ets {
		sort := d.rdb.HLen(context.Background(), d.opt.Key).Val()
		item := &redis.Z{
			Score:  float64(sort),
			Member: d.opt.Field,
		}
		err = d.rdb.ZAdd(context.Background(), d.GetSortKey(), item).Err()
	}
	return err
}

// Paginate 分页获取数据
func (d *Data) Paginate(sort string) (*pagination, error) {
	var rows, keys []string
	paginate := newPagination(d.opt)
	// 计算分页总页数
	paginate.countTotalPages(d.rdb.HLen(context.Background(), d.opt.Key).Val())
	// 获取会话的排序,并且分页
	switch sort {
	case "asc":
		keys, _ = d.rdb.ZRange(context.Background(), d.GetSortKey(), int64(paginate.GetOffset()), int64(paginate.GetEnd())).Result()
	case "desc":
		keys, _ = d.rdb.ZRevRange(context.Background(), d.GetSortKey(), int64(paginate.GetOffset()), int64(paginate.GetEnd())).Result()
	}
	for _, key := range keys {
		// 从redis中取出会话数据
		res := d.rdb.HGet(context.Background(), d.opt.Key, key).Val()
		if res != "" {
			rows = append(rows, res)
		}
	}
	paginate.Rows = rows
	return paginate, nil
}

// GetSortKey 获取排序key
func (d *Data) GetSortKey() string {
	return d.opt.Key + ":sort"
}
