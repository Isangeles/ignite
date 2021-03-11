/*
 * config.go
 *
 * Copyright 2021 Dariusz Sikora <dev@isangeles.pl>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
 * MA 02110-1301, USA.
 *
 *
 */

// Package for AI configuration values.
package config

import (
	"fmt"
	"os"

	"github.com/isangeles/flame/data/text"
)

const (
	Name, Version  = "Ignite", "0.1.0-dev"
	ConfigFileName = ".ignite"
)

var (
	ServerHost  = ""
	ServerPort  = ""
	UserID      = ""
	UserPass    = ""
)

// Load load server configuration file.
func Load() error {
	// Open config file.
	file, err := os.Open(ConfigFileName)
	if err != nil {
		return fmt.Errorf("Unale to open config file: %v", err)
	}
	defer file.Close()
	// Unmarshal config.
	conf, err := text.UnmarshalConfig(file)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal config: %v", err)
	}
	if len(conf["server"]) > 1 {
		ServerHost = conf["server"][0]
		ServerPort = conf["server"][1]
	}
	if len(conf["user"]) > 1 {
		UserID = conf["user"][0]
		UserPass = conf["user"][1]
	}
	return nil
}
