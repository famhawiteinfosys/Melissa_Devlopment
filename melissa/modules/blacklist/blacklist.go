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

package blacklist

import (
	"fmt"
	"html"
	"log"
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
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
)

func blacklist(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	allBlacklisted := sql.GetChatBlacklist(strconv.Itoa(chat.Id))
	filterList := "Current <b>blacklisted</b> words:\n"

	if len(args) > 0 && strings.ToLower(args[0]) == "copy" {
		filterList = ""
		for _, filter := range allBlacklisted {
			filterList += fmt.Sprintf("<code>%v</code>\n", html.EscapeString(filter.Trigger))
		}
	} else {
		for _, filter := range allBlacklisted {
			filterList += fmt.Sprintf(" - <code>%v</code>\n", html.EscapeString(filter.Trigger))
		}
	}
	splitText := helpers.SplitMessage(filterList)
	for _, text := range splitText {
		if text == "Current <b>blacklisted</b> words:\n" {
			_, err := msg.ReplyText("There are no blacklisted messages here!")
			return err
		}
		_, err := msg.ReplyHTML(text)
		error_handling.HandleErr(err)
	}
	return nil
}

func addBlacklist(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	// Permission checks
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, msg, u.EffectiveUser.Id) {
		return gotgbot.EndGroups{}
	}

	words := strings.SplitN(msg.Text, " ", 2)

	if len(words) > 1 {
		text := words[1]
		var toBlacklist []string
		for _, trigger := range strings.Split(text, "\n") {
			toBlacklist = append(toBlacklist, strings.TrimSpace(trigger))
		}
		go func(chatId string, toBlacklist []string) {
			for _, trigger := range toBlacklist {
				sql.AddToBlacklist(chatId, strings.ToLower(trigger))
			}
		}(strconv.Itoa(chat.Id), toBlacklist) // Use a goroutine to insert all blacklists in the background

		if len(toBlacklist) == 1 {
			_, err := msg.ReplyHTMLf("Added <code>%v</code> to the blacklist!", html.EscapeString(toBlacklist[0]))
			return err
		} else {
			_, err := msg.ReplyHTMLf("Added <code>%v</code> triggers to the blacklist!", len(toBlacklist))
			return err
		}
	} else {
		_, err := msg.ReplyText("Tell me which words you would like to add to the blacklist.")
		return err
	}
}

func unblacklist(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	// Permission checks
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, msg, u.EffectiveUser.Id) {
		return gotgbot.EndGroups{}
	}

	words := strings.SplitN(msg.Text, " ", 2)
	if len(words) > 1 {
		text := words[1]

		var toUnBlacklist []string
		for _, trigger := range strings.Split(text, "\n") {
			toUnBlacklist = append(toUnBlacklist, strings.TrimSpace(trigger))
		}

		successful := 0
		for _, trigger := range toUnBlacklist {
			success := sql.RmFromBlacklist(strconv.Itoa(chat.Id), strings.ToLower(trigger))
			if success {
				successful++
			}
		}

		if len(toUnBlacklist) == 1 {
			if successful > 0 {
				_, err := msg.ReplyHTMLf("Removed <code>%v</code> from the blacklist!", html.EscapeString(toUnBlacklist[0]))
				return err
			} else {
				_, err := msg.ReplyText("This isn't a blacklisted trigger!")
				return err
			}
		} else if successful == len(toUnBlacklist) {
			_, err := msg.ReplyHTMLf("Removed <code>%v</code> triggers from the blacklist!", successful)
			return err
		} else if successful == 0 {
			_, err := msg.ReplyText("None of these triggers exist, so they weren't removed.")
			return err
		} else {
			_, err := msg.ReplyHTMLf("Removed <code>%v</code> triggers from the blacklist."+
				" %v did not exist, so were not removed", successful, len(toUnBlacklist)-successful)
			return err
		}
	} else {
		_, err := msg.ReplyText("Tell me which words you would like to remove from the blacklist.")
		return err
	}
}

func delBlacklist(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat_status.IsUserAdmin(chat, u.EffectiveUser.Id) {
		return gotgbot.EndGroups{}
	}

	toMatch := extraction.ExtractText(msg)
	if toMatch == "" {
		return gotgbot.EndGroups{}
	}
	chatFilters := sql.GetChatBlacklist(strconv.Itoa(chat.Id))

	for _, trigger := range chatFilters {
		pattern := `( |^|[^\w])` + regexp.QuoteMeta(trigger.Trigger) + `( |$|[^\w])`
		re, err := regexp.Compile("(?i)" + pattern)
		error_handling.HandleErr(err)

		if re.MatchString(toMatch) {
			_, err := msg.Delete()
			if err != nil {
				if err.Error() != "Bad Request: message to delete not found" || err.Error() != "Bad Request: reply message not found" {
					error_handling.HandleErr(err)
				} else {
					log.Println("Error while deleting blacklist message")
				}
			}
			break
		}
	}
	return nil
}

var customFilter handlers.FilterFunc = func(message *ext.Message) bool {
	return (Filters.Text(message) || Filters.Command(message) || Filters.Sticker(message) || Filters.Photo(message)) && (Filters.Group(message))
}
var blacklistMessage = handlers.NewMessage(customFilter, delBlacklist)

func LoadBlacklist(u *gotgbot.Updater) {
	defer log.Println("Loading module blacklist")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("blacklist", []rune{'/', '!'}, blacklist))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("addblacklist", []rune{'/', '!'}, addBlacklist))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("rmblacklist", []rune{'/', '!'}, unblacklist))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("unblacklist", []rune{'/', '!'}, unblacklist))
	blacklistMessage.AllowEdited = true
	u.Dispatcher.AddHandlerToGroup(blacklistMessage, 11)
}
