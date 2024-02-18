package main

import "time"

type RedisCommand struct {
	command string
	args    []string
}

type RedisValues struct {
	Value     string
	ExpiresAt *time.Time
}

var redis = make(map[string]RedisValues)
