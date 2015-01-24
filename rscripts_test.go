package rscripts

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"testing"
)

func TestGetAllHMembers(t *testing.T) {
	reply, err := HGetAllMembers(client, "restaurants", "restaurant:")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
	t.Log(err)
}

func TestGetAllMembers(t *testing.T) {
	for i := 1; i <= 3; i++ {
		client.Do("SET", fmt.Sprintf("test:%d", i), fmt.Sprintf("value %d", i))
	}
	reply, err := GetAllMembers(client, "restaurants", "test:")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
	t.Log(err)
}

var client redis.Conn

func init() {
	var err error
	client, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalln(err)
	}
	Init(client)
}
