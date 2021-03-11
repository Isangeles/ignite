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

// main package loads configuration, connects to the game server,
// and starts AI process.
package main

import (
	"fmt"
	"time"
	"log"
	
	"github.com/isangeles/flame"
	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/module"

	"github.com/isangeles/fire/response"
	"github.com/isangeles/fire/request"
	
	"github.com/isangeles/ignite/ai"
	"github.com/isangeles/ignite/config"
	"github.com/isangeles/ignite/game"
)

var (
	charsAI *ai.AI
	server  *game.Server
)

// Main function.
func main() {
	fmt.Printf("*%s(%s)*\n", config.Name, config.Version)
	// Load config.
	err := config.Load()
	if err != nil {
		panic(fmt.Errorf("Unable to load config: %v", err))
	}
	// Connect to the server.
	serv, err := game.NewServer(config.ServerHost, config.ServerPort)
	if err != nil {
		panic(fmt.Errorf("Unable to create game server connection: %v",
			err))
	}
	server = serv
	server.SetOnResponseFunc(handleResponse)
	// Login to the server.
	loginReq := request.Login{config.UserID, config.UserPass}
	err = server.Send(request.Request{Login: []request.Login{loginReq}})
	if err != nil {
		panic(fmt.Errorf("Unable to send login request: %v", err))
	}
	update := time.Now()
	for {
		// Delta.
		dtNano := time.Since(update).Nanoseconds()
		delta := dtNano / int64(time.Millisecond) // delta to milliseconds
		// Update.
		if charsAI == nil {
			continue
		}
		charsAI.Update(delta)
		update = time.Now()
		time.Sleep(time.Duration(16) * time.Millisecond)
	}
}

// handleResponse handles response from the Fire server.
func handleResponse(resp response.Response) {
	if !resp.Logon {
		handleUpdateResponse(resp.Update)
		for _, r := range resp.NewChar {
			handleNewCharResponse(r)
		}
	}
	for _, r := range resp.Error {
		log.Printf("Server error response: %s", r)
	}
}

// handleUpdateResponse handles update response from the server.
func handleUpdateResponse(resp response.Update) {
	flameres.Clear()
	mod := module.New()
	mod.Apply(resp.Module)
	game := game.New(flame.NewGame(mod))
	game.SetServer(server)
	charsAI = ai.New(game)
}

// handleNewCharResponse handlres new character response from the server.
func handleNewCharResponse(resp response.NewChar) {
	if charsAI == nil {
		return
	}
	for _, c := range charsAI.Game().Module().Chapter().Characters() {
		if resp.ID == c.ID() && resp.Serial == c.Serial() {
			aiChar := game.NewCharacter(c, charsAI.Game())
			charsAI.Game().AddCharacter(aiChar)
		}
	}
}
