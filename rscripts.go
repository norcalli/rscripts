package rscripts

import (
	"github.com/garyburd/redigo/redis"
)

// r[i] = {a[i], redis.call('hgetall', KEYS[2] .. a[i])}
// r[i] = redis.call('hgetall', KEYS[2] .. a[i])

var scripts = struct {
	HGetAllMembers  string
	GetAllMembers   string
	IncrementAndAdd string
}{
	`local a=redis.call('smembers',KEYS[1])
local r={}
for i=1,#a do
  r[i*2-1] = a[i]
	r[i*2] = redis.call('hgetall', KEYS[2] .. a[i])
end
return r`,
	`local a=redis.call('smembers',KEYS[1])
local r={}
for i=1,#a do
  r[i*2-1] = a[i]
	r[i*2] = redis.call('get', KEYS[2] .. a[i])
end
return r`,
	"local x=redis.call('incr',KEYS[1]);redis.call('sadd', KEYS[2], x);return x",
}

var shas = struct {
	HGetAllMembers  string
	GetAllMembers   string
	IncrementAndAdd string
}{}

func AddScripts(client redis.Conn) error {
	client.Send("MULTI")
	client.Send("SCRIPT", "LOAD", scripts.HGetAllMembers)
	client.Send("SCRIPT", "LOAD", scripts.GetAllMembers)
	client.Send("SCRIPT", "LOAD", scripts.IncrementAndAdd)
	reply, err := redis.Strings(client.Do("EXEC"))
	if err != nil {
		return err
	}
	shas.HGetAllMembers = reply[0]
	shas.GetAllMembers = reply[1]
	shas.IncrementAndAdd = reply[2]
	return nil
}

func IncrementAndAdd(client redis.Conn, idKey, setKey string) (int64, error) {
	return redis.Int64(client.Do("EVALSHA", shas.IncrementAndAdd, 2, idKey, setKey))
}

type Member struct {
	ID    int64
	Value string
}

func GetAllMembers(client redis.Conn, idKey, prefix string) ([]Member, error) {
	// return client.Do("EVALSHA", shas.GetAllMembers, 2, idKey, prefix)
	reply, err := redis.Values(client.Do("EVALSHA", shas.GetAllMembers, 2, idKey, prefix))
	if err != nil {
		return nil, err
	}
	var result []Member
	if err := redis.ScanSlice(reply, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type HashMember struct {
	ID    int64
	Value []string
}

func HGetAllMembers(client redis.Conn, idKey, prefix string) ([]HashMember, error) {
	reply, err := redis.Values(client.Do("EVALSHA", shas.HGetAllMembers, 2, idKey, prefix))
	if err != nil {
		return nil, err
	}

	n := len(reply) / 2
	result := make([]HashMember, n)
	for i := 0; i < n; i++ {
		x := &result[i]
		if reply, err = redis.Scan(reply, &x.ID, &x.Value); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func Init(client redis.Conn) error {
	return AddScripts(client)
}
