package main

import "time"

type RedisCommand struct {
	command string
	args    []string
}

type RedisValues struct {
	Value string
}

var Dir string
var DbFileName string
var port string

var redis = make(map[string]RedisValues)

type ExpiryStruct struct {
	Key    string
	Expiry time.Duration
}

var expiryChannel = make(chan ExpiryStruct, 1000)
