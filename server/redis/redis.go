package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	REDIS_POOL_SIZE   = 10
	CACHE_EXPIRE_TIME = 60 * 60 * 24 * 10         //10 days
	FRESH_CACHE_LIMIT = CACHE_EXPIRE_TIME - 60*10 //10 minutes
)

func Dial(network, address string) (redis.Conn, error) {
	return redis.Dial(network, address)
}

func NewPool(network, address string) *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", address)
		if err != nil {
			return nil, err
		}
		return c, nil
	}, REDIS_POOL_SIZE)
}
func SetStats(c redis.Conn, repo string) error {
	c.Send("ZINCRBY", "top", 1, repo)
	c.Send("ZADD", "last", time.Now().Unix(), repo)
	return c.Flush()
}

func SetCache(c redis.Conn, repo, content, coverage string) error {
	c.Send("SETEX", repo, CACHE_EXPIRE_TIME, content)
	c.Send("SET", repo+".coverage", coverage)
	return c.Flush()
}

func GetRepo(c redis.Conn, repo string) (string, bool, error) {
	c.Send("GET", repo)
	c.Send("TTL", repo)
	if err := c.Flush(); err != nil {
		return "", false, nil
	}
	content, err := redis.String(c.Receive())
	if err == redis.ErrNil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	if ttl, err := redis.Int(c.Receive()); err == nil && ttl > FRESH_CACHE_LIMIT {
		return content, true, nil
	}

	return content, false, err
}

func Top(c redis.Conn, key string, count int) ([]string, error) {
	reply, err := redis.Strings(c.Do("ZREVRANGE", key, 0, count-1))
	if err != nil && err == redis.ErrNil {
		return nil, nil
	}
	return reply, err
}

func GetCoverage(c redis.Conn, repo string) (float64, error) {
	return redis.Float64(c.Do("GET", repo+".coverage"))
}
