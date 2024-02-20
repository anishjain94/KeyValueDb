package main

type RedisCommand struct {
	command string
	args    []string
}

type RedisValues struct {
	Value string
}

var redis = make(map[string]RedisValues)
