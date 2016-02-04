package redis

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	REDIS_POOL_SIZE   = 10
	CACHE_EXPIRE_TIME = 60 * 60 * 24 * 10         //10 days
	FRESH_CACHE_LIMIT = CACHE_EXPIRE_TIME - 60*10 //10 minutes

	TOP_SCRIPT = "local keys = redis.call('ZREVRANGE', KEYS[1], '0', KEYS[2]) for key,value in pairs(keys) do keys[key] = keys[key] .. ':'.. redis.call('GET', value .. '.coverage') end return keys"
)

var SHA string

func Dial(network, address string) (redis.Conn, error) {
	return redis.Dial(network, address)
}

func NewPool(network, address, password string) (pool *redis.Pool, err error) {
	pool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", address)
		if err != nil {
			return nil, err
		}
		if password != "" {
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
		}
		return c, nil
	}, REDIS_POOL_SIZE)

	conn := pool.Get()
	defer conn.Close()
	SHA, err = redis.String(conn.Do("SCRIPT", "LOAD", TOP_SCRIPT))
	return pool, err
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

type top struct {
	Repo     string
	Coverage string
	Color    string
}

func Top(c redis.Conn, key string, count int) ([]top, error) {
	repos, err := redis.Strings(c.Do("EVALSHA", SHA, 2, key, count-1))
	if err != nil && err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	tops := []top{}
	for _, repo := range repos {
		parts := strings.Split(repo, ":")
		if len(parts) == 1 {
			tops = append(tops, top{Repo: repo})
		} else {
			color := ""
			if f, err := strconv.ParseFloat(parts[1], 64); err == nil {
				if f >= 75 {
					color = "green"
				} else if f >= 25 {
					color = "orange"
				} else {
					color = "red"
				}
			}

			tops = append(tops, top{Repo: parts[0], Coverage: fmt.Sprintf("%s%%", parts[1]), Color: color})
		}
	}
	return tops, err
}

func GetCoverage(c redis.Conn, repo string) (float64, error) {
	return redis.Float64(c.Do("GET", repo+".coverage"))
}
