/*
 * ignite.go
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

// main package starts the AI process.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/data/res"

	"github.com/isangeles/fire/request"
	"github.com/isangeles/fire/response"

	"github.com/isangeles/ignite/ai"
	"github.com/isangeles/ignite/config"
)

var (
	AI     *ai.AI
	server *ai.Server
)

// Main function.
func main() {
	log.Printf("%s(%s)", config.Name, config.Version)
	// Load config.
	err := config.Load()
	if err != nil {
		panic(fmt.Errorf("Unable to load config: %v", err))
	}
	// Connect to the server.
	server, err = ai.NewServer(config.ServerHost, config.ServerPort)
	if err != nil {
		panic(fmt.Errorf("Unable to create game server connection: %v",
			err))
	}
	server.SetOnResponseFunc(handleResponse)
	// Login to the server.
	loginReq := request.Login{config.UserID, config.UserPass}
	err = server.Send(request.Request{Login: []request.Login{loginReq}})
	if err != nil {
		panic(fmt.Errorf("Unable to send login request: %v", err))
	}
	update := time.Now()
	for !server.Closed() {
		// Delta.
		dtNano := time.Since(update).Nanoseconds()
		delta := dtNano / int64(time.Millisecond) // delta to milliseconds
		// Update.
		if AI == nil {
			continue
		}
		AI.Update(delta)
		AI.Game().Update(delta)
		update = time.Now()
		// Update break.
		time.Sleep(time.Duration(16) * time.Millisecond)
	}
}

// handleResponse handles response from the server.
func handleResponse(resp response.Response) {
	if !resp.Logon {
		handleUpdateResponse(resp.Update)
		for _, r := range resp.Character {
			handleCharacterResponse(r)
		}
	}
	for _, r := range resp.Error {
		log.Printf("Server error response: %s", r)
	}
}

// handleUpdateResponse handles update response from the server.
func handleUpdateResponse(resp response.Update) {
	if AI != nil {
		return
	}
	res.Clear()
	mod := flame.NewModule(resp.Module)
	game := ai.NewGame(mod)
	game.SetServer(server)
	AI = ai.New(game)
}

// handleCharacterResponse handles character response from the server.
func handleCharacterResponse(resp response.Character) {
	if AI == nil {
		return
	}
	for _, c := range AI.Game().Characters() {
		if c.ID() == resp.ID && c.Serial() == resp.Serial {
			return
		}
	}
	for _, c := range AI.Game().Chapter().Characters() {
		if resp.ID == c.ID() && resp.Serial == c.Serial() {
			aiChar := ai.NewCharacter(c, AI.Game())
			AI.Game().AddCharacter(aiChar)
		}
	}
}
