/*
 * game.go
 *
 * Copyright 2021 Dariusz Sikora <dev@isangeles.pl>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either Version 2 of the License, or
 * (at your option) any later Version.
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

package ai

import (
	"log"

	"github.com/isangeles/flame"
)

// Struct for game wrapper.
type Game struct {
	*flame.Module
	server      *Server
	characters  map[string]*Character
	onLoginFunc func(g *Game)
}

// NewGame creates new AI game wrapper for specified module.
func NewGame(module *flame.Module) *Game {
	g := Game{
		Module:     module,
		characters: make(map[string]*Character),
	}
	return &g
}

// AddCharacter adds character to control by the game AI.
func (g *Game) AddCharacter(c *Character) {
	g.characters[c.ID()+c.Serial()] = c
}

// RemoveCharacter removes character from game AI control.
func (g *Game) RemoveCharacter(c *Character) {
	delete(g.characters, c.ID()+c.Serial())
}

// Character returns game characters.
func (g *Game) Characters() (chars []*Character) {
	for _, char := range g.characters {
		chars = append(chars, char)
	}
	return chars
}

// SetServer sets remote game server.
func (g *Game) SetServer(server *Server) {
	g.server = server
	g.Server().SetOnResponseFunc(g.handleResponse)
	err := g.Server().Update()
	if err != nil {
		log.Printf("Game: unable to send update request to the server: %v",
			err)
	}
}

// Server retruns game server.
func (g *Game) Server() *Server {
	return g.server
}

// SetOnLoginFunc sets function triggered on login.
func (g *Game) SetOnLoginFunc(f func(g *Game)) {
	g.onLoginFunc = f
}
