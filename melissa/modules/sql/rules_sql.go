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

package sql

import (
	"encoding/json"
	"fmt"

	"github.com/wI2L/jettison"

	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/caching"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/error_handling"
)

type Rules struct {
	ChatId string `gorm:"primary_key" json:"chat_id"`
	Rules  string `json:"rules"`
}

func GetChatRules(chatId string) *Rules {
	ruleJson, err := caching.CACHE.Get(fmt.Sprintf("rules_%v", chatId))
	var rules *Rules
	if err != nil {
		rules = cacheRules(chatId)
	}
	_ = json.Unmarshal(ruleJson, &rules)
	return rules
}

func SetChatRules(chatId, rules string) {
	SESSION.Save(&Rules{ChatId: chatId, Rules: rules})
	cacheRules(chatId)
}

func cacheRules(chatId string) *Rules {
	rules := &Rules{}
	SESSION.Where("chat_id = ?", chatId).Find(&rules)
	ruleJson, _ := jettison.Marshal(&rules)
	err := caching.CACHE.Set(fmt.Sprintf("rules_%v", chatId), ruleJson)
	error_handling.HandleErr(err)
	return rules
}
