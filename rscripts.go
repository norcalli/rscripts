package rscripts

import (
	"github.com/garyburd/redigo/redis"
)

var incrAddSha string

func AddScripts(client redis.Conn) error {
	var err error
	incrAddSha, err = redis.String(client.Do("SCRIPT", "LOAD", "local x=redis.call('incr', KEYS[1]);redis.call('sadd', KEYS[2], x);return x"))
	if err != nil {
		return err
	}
	return nil
}

func IncrementAndAdd(client redis.Conn, idKey, setKey string) (int64, error) {
	return redis.Int64(client.Do("EVALSHA", incrAddSha, 2, idKey, setKey))
}

func Init(client redis.Conn) error {
	return AddScripts(client)
}
