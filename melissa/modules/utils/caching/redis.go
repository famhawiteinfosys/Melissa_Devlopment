/*
 *    Copyright © 2019 - 2023 Famhawite Infosys Project Arise
 *    This file is part of Famhawite Infosys
 *
 *    Melissa is free software: you can redistribute it and/or modify
 *    it under the terms of the Raphielscape Public License as published by
 *    the Devscapes Open Source Holding GmbH., version 1.d
 *
 *    Melissa is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    Devscapes Raphielscape Public License for more details.
 *
 *    You should have received a copy of the Devscapes Raphielscape Public License
 */

package caching

import (
	"time"

	"github.com/famhawiteinfosys/Melissa/melissa"
	"github.com/go-redis/redis"
)

var REDIS *redis.Client

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:         Melissa.BotConfig.RedisAddress,
		Password:     Melissa.BotConfig.RedisPassword,
		DB:           0,
		DialTimeout:  time.Second,
		MinIdleConns: 0,
	})
	REDIS = client
}
