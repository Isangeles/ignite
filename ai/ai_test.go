/*
 * ai_test.go
 *
 * Copyright 2022 Dariusz Sikora <ds@isangeles.dev>
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
	"testing"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/character"
	"github.com/isangeles/flame/data/res"

	"github.com/isangeles/ignite/config"
)

var charData = res.CharacterData{ID: "char", Level: 1, Attributes: res.AttributesData{5, 5, 5, 5, 5}}

// TestUpdateMoveAround test moving around by AI.
func TestUpdateMoveAronud(t *testing.T) {
	mod := flame.NewModule(res.ModuleData{})
	game := NewGame(mod)
	char := NewCharacter(character.New(charData), game)
	game.AddCharacter(char)
	ai := New(game)
	ai.Update(config.MoveFreq)
	posX, posY := char.Position()
	destX, destY := char.DestPoint()
	if posX == destX && posY == destY {
		t.Fatalf("Character was not moved")
	}
}
