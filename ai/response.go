/*
 * response.go
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

package ai

import (
	"log"
	
	"github.com/isangeles/flame/data/res"

	"github.com/isangeles/fire/response"
)

// handleResponse handles specified response from Fire server.
func (g *Game) handleResponse(resp response.Response) {
	if !resp.Logon && g.onLoginFunc != nil {
		g.onLoginFunc(g)
	}
	g.handleUpdateResponse(resp.Update)
	g.handleNewCharResponse(resp.NewChar)
	for _, r := range resp.Error {
		log.Printf("Game server error: %s", r)
	}
}

// handleNewCharResponse handles new characters from server response.
func (g *Game) handleNewCharResponse(resp []response.NewChar) {
	for _, r := range resp {
		char := g.Module().Chapter().Character(r.ID, r.Serial)
		if char == nil {
			log.Printf("Game server: handle new-char response: unable to find character in module: %s %s",
				r.ID, r.Serial)
			return
		}
		g.AddCharacter(NewCharacter(char, g))
	}
}

// handleUpdateRespone handles update response.
func (g *Game) handleUpdateResponse(resp response.Update) {
	res.Clear()
	g.Module().Apply(resp.Module)
}
