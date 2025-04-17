/*
 * character.go
 *
 * Copyright 2021-2025 Dariusz Sikora <ds@isangeles.dev>
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
	"math"

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
	game        *Game
	onUseEvents []func(o useaction.Usable)
}

// NewCharacter creates new game character.
func NewCharacter(char *character.Character, game *Game) *Character {
	c := Character{
		Character: char,
		game:      game,
	}
	return &c
}

// AddOnUseEvent adds function to trigger after using an usable object.
func (c *Character) AddOnUseEvent(event func(o useaction.Usable)) {
	c.onUseEvents = append(c.onUseEvents, event)
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
	c.ChatLog().Add(objects.NewMessage(message, false))
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
		// If no server then trigger onUse event and return.
		// With server this event should be triggered after
		// use response from the server.
		for _, event := range c.onUseEvents {
			event(ob)
		}
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

// MoveCloseTo moves character to the position at minimal range
// to the specified position.
func (c *Character) MoveCloseTo(x, y, minRange float64) {
	charX, charY := c.Position()
	switch {
	case x > charX:
		x -= minRange
	case x < charX:
		x += minRange
	case y > charY:
		y -= minRange
	case y < charY:
		y += minRange
	}
	c.SetDestPoint(x, y)
}

// Retruns distance from the character default position.
func (c *Character) DefPosDistance() float64 {
	posX, posY := c.Position()
	defX, defY := c.DefaultPosition()
	return math.Hypot(posX-defX, posY-defY)
}

// hasHostileTarget checks if character first target is
// hostile.
func (c *Character) hasHostileTarget() bool {
	return len(c.Targets()) > 0 && c.AttitudeFor(c.Targets()[0]) == character.Hostile
}

// meetTargetRangeReqs check if all target range requirements are meet.
// Returns true, if none of specified requirements is a target range
// requirement.
func (c *Character) meetTargetRangeReqs(reqs ...req.Requirement) bool {
	tarRangeReqs := make([]req.Requirement, 0)
	for _, r := range reqs {
		if r, ok := r.(*req.TargetRange); ok {
			tarRangeReqs = append(tarRangeReqs, r)
		}
	}
	if !c.MeetReqs(tarRangeReqs...) {
		return false
	}
	return true
}
