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
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/famhawiteinfosys/Melissa/melissa/modules/sql"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/chat_status"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/error_handling"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/extraction"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

func addWarnFilter(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	user := u.EffectiveUser
	var keyword string
	var content string

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}

	args := strings.SplitN(msg.Text, " ", 2)

	if len(args) < 2 {
		return nil
	}

	extracted := helpers.SplitQuotes(args[1])

	if len(extracted) >= 2 {
		keyword = strings.ToLower(extracted[0])
		content = extracted[1]
	} else {
		return nil
	}

	go sql.AddWarnFilter(strconv.Itoa(chat.Id), keyword, content)
	_, err := msg.ReplyText(fmt.Sprintf("Warn handler added for '%v'!", keyword))
	error_handling.HandleErr(err)
	return gotgbot.EndGroups{}
}

func removeWarnFilter(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}

	args := strings.SplitN(msg.Text, " ", 2)

	if len(args) < 2 {
		return gotgbot.EndGroups{}
	}

	extracted := helpers.SplitQuotes(args[1])

	if len(extracted) < 1 {
		return nil
	}

	toRemove := extracted[0]

	chatFilters := sql.GetChatWarnTriggers(strconv.Itoa(chat.Id))

	if chatFilters == nil {
		_, err := msg.ReplyText("No warning filters are active here!")
		return err
	}

	for _, filt := range chatFilters {
		if filt.Keyword == toRemove {
			go sql.RemoveWarnFilter(strconv.Itoa(chat.Id), toRemove)
			_, err := msg.ReplyText("Yep, I'll stop warning people for that.")
			error_handling.HandleErr(err)
			return gotgbot.EndGroups{}
		}
	}
	_, err := msg.ReplyText("That's not a current warning filter - run /warnlist for all active warning filters.")
	return err
}

func listWarnFilters(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	allHandlers := sql.GetChatWarnTriggers(strconv.Itoa(chat.Id))

	if allHandlers == nil {
		_, err := u.EffectiveMessage.ReplyText("No warning filters are active here!")
		return err
	}

	filterList := "<b>Current warning filters in this chat:</b>\n"
	for _, handler := range allHandlers {
		entry := fmt.Sprintf(" - %v\n", html.EscapeString(handler.Keyword))
		if len(entry)+len(filterList) > 4096 {
			_, err := u.EffectiveMessage.ReplyHTML(filterList)
			error_handling.HandleErr(err)
			filterList = entry
		} else {
			filterList += entry
		}
	}
	if filterList != "<b>Current warning filters in this chat:</b>\n" {
		_, err := u.EffectiveMessage.ReplyHTML(filterList)
		return err
	}
	return nil
}

func replyFilter(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	message := u.EffectiveMessage

	if !chat_status.IsBotAdmin(chat, nil) {
		return gotgbot.EndGroups{}
	}
	if chat_status.IsUserAdmin(chat, u.EffectiveUser.Id) {
		return gotgbot.ContinueGroups{}
	}

	chatWarnFilters := sql.GetChatWarnTriggers(strconv.Itoa(chat.Id))
	toMatch := extraction.ExtractText(message)
	if toMatch == "" {
		return nil
	}

	for _, handler := range chatWarnFilters {
		pattern := `( |^|[^\w])` + regexp.QuoteMeta(handler.Keyword) + `( |$|[^\w])`
		re, err := regexp.Compile("(?i)" + pattern)
		error_handling.HandleErr(err)

		if re.MatchString(toMatch) {
			user := u.EffectiveUser
			warnFilter := sql.GetWarnFilter(strconv.Itoa(chat.Id), handler.Keyword)
			return warn(user, chat, warnFilter.Reply, message)
		}
	}
	return gotgbot.ContinueGroups{}
}
