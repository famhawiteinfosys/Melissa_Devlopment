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
	"log"

	"github.com/famhawiteinfosys/Melissa/melissa"
	"github.com/famhawiteinfosys/Melissa/melissa/modules/utils/error_handling"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

var SESSION *gorm.DB

func init() {
	conn, err := pq.ParseURL(Melissa.BotConfig.SqlUri)
	error_handling.FatalError(err)

	db, err := gorm.Open("postgres", conn)
	error_handling.FatalError(err)

	if Melissa.BotConfig.DebugMode == "True" {
		SESSION = db.Debug()
		log.Println("[INFO][Database] Using database in debug mode.")
	} else {
		SESSION = db
	}

	db.DB().SetMaxOpenConns(100)

	log.Println("[INFO][Database] Database connected")

	// Create tables if they don't exist
	SESSION.AutoMigrate(&User{}, &Chat{}, &Warns{}, &WarnFilters{}, &WarnSettings{}, &BlackListFilters{}, &Federation{},
		&FedChat{}, &FedAdmin{}, &FedBan{}, &Note{}, &Button{}, &Welcome{}, &WelcomeButton{}, &MutedUser{}, &Rules{})
	log.Println("Auto-migrated database schema")

}
