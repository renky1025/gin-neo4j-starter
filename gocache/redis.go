package gocache

import (
	"context"
	"strconv"
	"time"

	"go-gin-restful-service/config"
	"go-gin-restful-service/log"

	redis "github.com/redis/go-redis/v9"
)

type RedisCache struct {
	RedisClient *redis.Client
}

var RedisDB *RedisCache
var ctx = context.Background()

// // GetRedis ...to be concurrence secure
// func InitRedis() *RedisCache {
// 	cfg := config.GetConfig()
// 	var redisConfig = cfg.Redis
// 	rdbClient := redis.NewClient(&redis.Options{
// 		Network:  "tcp",
// 		Addr:     redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port),
// 		Password: redisConfig.Password, // 密码
// 		DB:       redisConfig.DB,       // 默认数据库
// 		PoolSize: 50,                   // 连接池大小
// 	})

// 	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancelFunc()
// 	result, err := rdbClient.Ping(ctx).Result()
// 	log.Logger.Info("redis connection result : " + result)
// 	if err != nil {
// 		log.Logger.Panic(err)
// 	}
// 	RedisDB = &RedisCache{
// 		RedisClient: rdbClient,
// 	}
// 	return RedisDB
// }

func NewRedisCache() *RedisCache {
	cfg := config.GetConfig()
	var redisConfig = cfg.Redis
	rdbClient := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port),
		Password: redisConfig.Password, // 密码
		DB:       redisConfig.DB,       // 默认数据库
		PoolSize: 50,                   // 连接池大小
	})

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	result, err := rdbClient.Ping(ctx).Result()
	log.Logger.Info("redis connection result : " + result)
	if err != nil {
		log.Logger.Panic(err)
	}
	RedisDB = &RedisCache{
		RedisClient: rdbClient,
	}
	return RedisDB
}

func (rdb *RedisCache) getCacheKey(key string) string {
	return config.GetConfig().Redis.Prefix + key
}

// Set 设置键值对
func (rdb *RedisCache) Set(key string, value interface{}, timeSec int) bool {
	err := rdb.RedisClient.Set(ctx,
		rdb.getCacheKey(key), value, time.Duration(timeSec)*time.Second).Err()
	if err != nil {
		log.Logger.Errorf("redisUtil.Set err: err=[%+v]", err)
		return false
	}
	return true
}

// Incr 设置键值递增
func (rdb *RedisCache) Incr(key string) int64 {
	count := rdb.RedisClient.Incr(ctx, rdb.getCacheKey(key)).Val()
	err := rdb.RedisClient.Expire(ctx, rdb.getCacheKey(key), rdb.GetTodayLastSeconds()).Err()
	if err != nil {
		log.Logger.Errorf("redisUtil.Expire err: err=[%+v]", err)
		return 100
	}
	return count
}

func (rdb *RedisCache) GetTodayLastSeconds() time.Duration {
	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)
	remainSecond := time.Duration(todayLastTime.Unix()-time.Now().Local().Unix()) * time.Second
	return remainSecond
}

// Get 获取key的值
func (rdb *RedisCache) Get(key string) interface{} {
	res, err := rdb.RedisClient.Get(ctx, rdb.getCacheKey(key)).Result()
	if err != nil {
		log.Logger.Errorf("redisUtil.Get err: err=[%+v]", err)
		return ""
	}
	return res
}

func (rdb *RedisCache) Scan(key string) map[string]string {
	results := make(map[string]string, 0)

	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.RedisClient.Scan(ctx, cursor, rdb.getCacheKey(key), 20).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)

		var value string
		for _, key := range keys {
			value, _ = rdb.RedisClient.Get(ctx, key).Result()
			results[key] = value
		}
		if cursor == 0 {
			break
		}
	}
	return results
}

// Exists 判断多项key是否存在
func (rdb *RedisCache) Exists(keys ...string) int64 {
	fullKeys := rdb.toFullKeys(keys)
	cnt, err := rdb.RedisClient.Exists(ctx, fullKeys...).Result()
	if err != nil {
		log.Logger.Errorf("redisUtil.Exists err: err=[%+v]", err)
		return -1
	}
	return cnt
}

// Expire 指定缓存失效时间
func (rdb *RedisCache) Expire(key string, timeSec int) bool {
	err := rdb.RedisClient.Expire(ctx, rdb.getCacheKey(key), time.Duration(timeSec)*time.Second).Err()
	if err != nil {
		log.Logger.Errorf("redisUtil.Expire err: err=[%+v]", err)
		return false
	}
	return true
}

// TTL 根据key获取过期时间
func (rdb *RedisCache) TTL(key string) int {
	td, err := rdb.RedisClient.TTL(ctx, rdb.getCacheKey(key)).Result()
	if err != nil {
		log.Logger.Errorf("redisUtil.TTL err: err=[%+v]", err)
		return 0
	}
	return int(td / time.Second)
}

func (rdb *RedisCache) Delete(key string) error {
	rdb.Del(key)
	return nil
}

func (rdb *RedisCache) Clear() error {
	//清空当前数据库，因为连接的是索引为0的数据库，所以清空的就是0号数据库
	res, err := rdb.RedisClient.FlushDB(ctx).Result()
	if err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Info(res)
	return err
}

// Del 删除一个或多个键
func (rdb *RedisCache) Del(keys ...string) bool {
	fullKeys := rdb.toFullKeys(keys)
	err := rdb.RedisClient.Del(ctx, fullKeys...).Err()
	if err != nil {
		log.Logger.Errorf("redisUtil.Del err: err=[%+v]", err)
		return false
	}
	return true
}

// toFullKeys 为keys批量增加前缀
func (rdb *RedisCache) toFullKeys(keys []string) (fullKeys []string) {
	for _, k := range keys {
		fullKeys = append(fullKeys, rdb.getCacheKey(k))
	}
	return
}
