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

package warns

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/famhawiteinfosys/Melissa/melissa/modules/sql"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/chat_status"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

func setWarnLimit(_ ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	user := u.EffectiveUser

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}

	if len(args) > 0 {
		for _, elem := range args[0] {
			if unicode.IsDigit(elem) {
				num, err := strconv.Atoi(args[0])
				error_handling.HandleErr(err)

				if num < 3 {
					_, err := msg.ReplyText("The minimum warn limit is 3!")
					return err
				} else {
					go sql.SetWarnLimit(strconv.Itoa(chat.Id), num)
					_, err := msg.ReplyHTML(fmt.Sprintf("Updated the warn limit to <b>%v</b>", num))
					return err
				}
			}
		}
		_, err := msg.ReplyText("Give me a number as an argument!")
		return err
	} else {
		limit, softWarn := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		_, err := msg.ReplyHTML(fmt.Sprintf("The current warn limit is <b>%v</b>.\nThe soft warn setting is set to <b>%v</b>.", limit, softWarn))
		return err
	}
}

func setWarnStrength(_ ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	user := u.EffectiveUser

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}

	if len(args) > 0 {
		if strings.ToLower(args[0]) == "on" || strings.ToLower(args[0]) == "yes" {
			go sql.SetWarnStrength(strconv.Itoa(chat.Id), false)
			_, err := msg.ReplyText("Too many warns will now result in a ban!")
			return err
		} else if strings.ToLower(args[0]) == "off" || strings.ToLower(args[0]) == "no" {
			go sql.SetWarnStrength(strconv.Itoa(chat.Id), true)
			_, err := msg.ReplyText("Too many warns will now result in a kick! User will be able to join again after.")
			return err
		} else {
			_, err := msg.ReplyText("I only understand on/yes/no/off!")
			return err
		}
	} else {
		limit, softWarn := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		if softWarn {
			_, err := msg.ReplyHTML(fmt.Sprintf("Warns are currently set to <b>kick</b> users when they exceed <b>%v</b> warns.", limit))
			return err
		} else {
			_, err := msg.ReplyHTML(fmt.Sprintf("Warns are currently set to <b>ban</b> users when they exceed <b>%v</b> warns.", limit))
			return err
		}
	}
}
