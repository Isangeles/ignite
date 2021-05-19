/*
 * character.go
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

	"github.com/isangeles/flame/character"
	"github.com/isangeles/flame/effect"
	"github.com/isangeles/flame/objects"
	"github.com/isangeles/flame/req"
	"github.com/isangeles/flame/serial"
	"github.com/isangeles/flame/useaction"

	"github.com/isangeles/fire/request"
)

// Wrapper struct for AI character.
type Character struct {
	*character.Character
	game *Game
}

// NewCharacter creates new game character.
func NewCharacter(char *character.Character, game *Game) *Character {
	c := Character{
		Character: char,
		game:      game,
	}
	return &c
}

// SetDestPoint sets a specified XY position as current
// as a character destination point.
func (c *Character) SetDestPoint(x, y float64) {
	c.Character.SetDestPoint(x, y)
	if c.game.Server() == nil {
		return
	}
	moveReq := request.Move{c.ID(), c.Serial(), x, y}
	req := request.Request{Move: []request.Move{moveReq}}
	err := c.game.Server().Send(req)
	if err != nil {
		log.Printf("Character: %s %s: unable to send move request: %v",
			c.ID(), c.Serial(), err)
	}
}

// AddChatMessage adds new message to character chat log.
func (c *Character) AddChatMessage(message string) {
	c.ChatLog().Add(objects.Message{Text: message})
	if c.game.Server() == nil {
		return
	}
	chatReq := request.Chat{c.ID(), c.Serial(), message, false}
	req := request.Request{Chat: []request.Chat{chatReq}}
	err := c.game.Server().Send(req)
	if err != nil {
		log.Printf("Character: %s %s: unable to send chat request: %v",
			c.ID(), c.Serial(), err)
	}
}

// SetTarget sets specified targetable object as current target.
func (c *Character) SetTarget(tar effect.Target) {
	c.Character.SetTarget(tar)
	if c.game.Server() == nil {
		return
	}
	targetReq := request.Target{
		ObjectID:     c.ID(),
		ObjectSerial: c.Serial(),
	}
	if tar != nil {
		targetReq.TargetID, targetReq.TargetSerial = tar.ID(), tar.Serial()
	}
	req := request.Request{Target: []request.Target{targetReq}}
	err := c.game.Server().Send(req)
	if err != nil {
		log.Printf("Character: %s %s: unable to send target request to the server: %v",
			c.ID(), c.Serial(), err)
	}
}

// Use uses specified usable object.
func (c *Character) Use(ob useaction.Usable) {
	err := c.Character.Use(ob)
	if err != nil {
		return
	}
	if c.game.Server() == nil {
		return
	}
	useReq := request.Use{
		UserID:     c.ID(),
		UserSerial: c.Serial(),
		ObjectID:   ob.ID(),
	}
	if ob, ok := ob.(serial.Serialer); ok {
		useReq.ObjectSerial = ob.Serial()
	}
	req := request.Request{Use: []request.Use{useReq}}
	err = c.game.Server().Send(req)
	if err != nil {
		log.Printf("Character: %s %s: unable to send use request: %v",
			c.ID(), c.Serial(), err)
	}
}

// meetTargetRangeReqs check if all target range requirements are meet.
// Returns true, if none of specified requirements is a target range
// requirement.
func (c *Character) meetTargetRangeReqs(reqs ...req.Requirement) bool {
	for _, r := range reqs {
		if r, ok := r.(*req.TargetRange); ok {
			if !c.MeetReq(r) {
				return false
			}
		}
	}
	return true
}
