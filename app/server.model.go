package main

type RedisCommand struct {
	command string
	args    []string
}

var redis = make(map[string]string)
