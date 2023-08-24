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

package bans

import (
	"log"
	"strings"

	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/chat_status"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/error_handling"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/extraction"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/string_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
)

func ban(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userId, _ := extraction.ExtractUserAndText(message, args)
	if userId == 0 {
		_, err := message.ReplyText("Try targeting a user next time bud.")
		return err
	}

	member, err := chat.GetMember(userId)
	if err != nil {
		if err.Error() == "User not found" {
			_, err = message.ReplyText("This user is ded mate.")
		}
		return err
	}

	userMember, _ := chat.GetMember(user.Id)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to ban users!")
		return err
	}

	if chat_status.IsUserBanProtected(chat, userId, member) {
		_, err := message.ReplyText("One day I'll find out how to work around the bot API. Today is not that day.")
		return err
	}

	if userId == bot.Id {
		_, err := message.ReplyText("No u")
		return err
	}

	_, err = chat.KickMember(userId)
	if err != nil {
		return err
	}

	_, err = message.ReplyText("Banned!")
	return err
}

func tempBan(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userId, reason := extraction.ExtractUserAndText(message, args)
	if userId == 0 {
		_, err := message.ReplyText("Try targeting a user next time bud.")
		return err
	}

	member, err := chat.GetMember(userId)
	if err != nil {
		if err.Error() == "User not found" {
			_, err = message.ReplyText("This user is ded mate.")
		}
		return err
	}

	userMember, _ := chat.GetMember(user.Id)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to ban users!")
		return err
	}

	if chat_status.IsUserBanProtected(chat, userId, member) {
		_, err := message.ReplyText("One day I'll find out how to work around the bot API. Today is not that day.")
		return err
	}

	if userId == bot.Id {
		_, err := message.ReplyText("No u")
		return err
	}

	if reason == "" {
		_, err := message.ReplyText("I don't know how long I'm supposed to ban them for 🤔.")
		return err
	}

	splitReason := strings.SplitN(reason, " ", 2)
	timeVal := splitReason[0]
	banTime := string_handling.ExtractTime(message, timeVal)
	if banTime == -1 {
		return nil
	}
	newMsg := bot.NewSendableKickChatMember(chat.Id, userId)
	string_handling.ExtractTime(message, timeVal)
	newMsg.UntilDate = banTime
	_, err = newMsg.Send()
	if err != nil {
		_, err := message.ReplyText("Press F, I can't seem to ban this user.")
		error_handling.HandleErr(err)
	}
	_, err = message.ReplyHTMLf("Banned for <b>%s</b>!", timeVal)
	return err
}

func kick(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userId, _ := extraction.ExtractUserAndText(message, args)
	if userId == 0 {
		_, err := message.ReplyText("Try targeting a user next time bud.")
		return err
	}

	var member, err = chat.GetMember(userId)
	if err != nil {
		if err.Error() == "User not found" {
			_, err = message.ReplyText("This user is ded mate.")
		}
		return err
	}

	userMember, _ := chat.GetMember(user.Id)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to ban users!")
		return err
	}

	if chat_status.IsUserBanProtected(chat, userId, member) {
		_, err := message.ReplyText("One day I'll find out how to work around the bot API. Today is not that day.")
		return err
	}

	if userId == bot.Id {
		_, err := message.ReplyText("No u")
		return err
	}

	_, err = chat.UnbanMember(userId) // Apparently unban on current user = kick
	if err != nil {
		_, err = message.ReplyText("Hec, I can't seem to kick this user.")
		return err
	}
	_, err = message.ReplyText("Kicked!")
	return err
}

func kickme(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}

	if chat_status.IsUserAdmin(chat, user.Id) {
		_, err := message.ReplyText("Admin sir pls ;_;")
		error_handling.HandleErr(err)
		return gotgbot.EndGroups{}
	}
	bb, _ := chat.UnbanMember(user.Id)
	if bb {
		_, err := message.ReplyText("Sure thing boss.")
		return err
	} else {
		_, err := message.ReplyText("OwO I can't :/")
		return err
	}
}

func unban(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) && chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userId, _ := extraction.ExtractUserAndText(message, args)

	if userId == 0 {
		_, err := message.ReplyText("Try targeting a user next time bud.")
		return err
	}

	_, err := chat.GetMember(userId)
	if err != nil {
		_, err := message.ReplyText("This user is ded m8.")
		return err
	}

	userMember, _ := chat.GetMember(user.Id)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to unban users!")
		return err
	}

	if userId == bot.Id {
		_, err := message.ReplyText("What exactly are you attempting to do?.")
		return err
	}

	if chat_status.IsUserInChat(chat, userId) {
		_, err := message.ReplyText("This user is already in the group!")
		return err
	}

	_, err = chat.UnbanMember(userId)
	error_handling.HandleErr(err)
	_, err = message.ReplyText("Fine, I'll allow it, this time...")
	return err
}

func LoadBans(u *gotgbot.Updater) {
	defer log.Println("Loading module bans")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("tban", []rune{'/', '!'}, tempBan))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("ban", []rune{'/', '!'}, ban))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("kick", []rune{'/', '!'}, kick))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("kickme", []rune{'/', '!'}, kickme))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("unban", []rune{'/', '!'}, unban))
}
