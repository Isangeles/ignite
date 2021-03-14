/*
 * ai.go
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

// ai package provides AI structs for controling
// module characters.
package ai

import (
	"fmt"

	"github.com/isangeles/flame/module/character"
	"github.com/isangeles/flame/module/effect"
	"github.com/isangeles/flame/module/skill"
	"github.com/isangeles/flame/rng"

	"github.com/isangeles/ignite/config"
)

// Struct for controlling non-player characters.
type AI struct {
	game      *Game
	moveTimer int64
	chatTimer int64
}

// New creates new AI for specified game.
func New(g *Game) *AI {
	ai := new(AI)
	ai.game = g
	return ai
}

// Update updates AI.
func (ai *AI) Update(delta int64) {
	ai.moveTimer += delta
	ai.chatTimer += delta
	// NPCs.
	for _, npc := range ai.Game().Characters() {
		// Move around.
		if ai.moveTimer >= config.MoveFreq {
			if npc.Casted() != nil || npc.Moving() || npc.Fighting() || npc.Agony() {
				continue
			}
			posX, posY := npc.Position()
			defX, defY := npc.DefaultPosition()
			if posX != defX || posY != defY {
				npc.SetDestPoint(defX, defY)
				continue
			}
			ai.moveAround(npc)
		}
		// Random chat.
		if ai.chatTimer >= config.ChatFreq {
			if npc.Casted() != nil || npc.Moving() || npc.Fighting() || npc.Agony() {
				continue
			}
			ai.saySomething(npc)
		}
		// Combat.
		if len(npc.Targets()) < 1 || npc.AttitudeFor(npc.Targets()[0]) != character.Hostile {
			// Look for hostile target.
			var tar effect.Target
			area := ai.Game().Module().Chapter().CharacterArea(npc.Character)
			if area == nil {
				continue
			}
			for _, t := range area.NearTargets(npc, npc.SightRange()) {
				if t == npc {
					continue
				}
				if npc.AttitudeFor(t) == character.Hostile {
					tar = t
					break
				}
			}
			if tar == nil {
				continue
			}
			npc.SetTarget(tar)
		}
		if npc.Fighting() {
			skill := ai.combatSkill(npc, npc.Targets()[0])
			if skill == nil {
				continue
			}
			npc.Use(skill)
		}
		break
	}
	// Reset timers.
	if ai.moveTimer >= config.MoveFreq {
		ai.moveTimer = 0
	}
	if ai.chatTimer >= config.ChatFreq {
		ai.chatTimer = 0
	}
}

// Game returns AI game.
func (ai *AI) Game() *Game {
	return ai.game
}

// moveAround moves specified character in random direction.
func (ai *AI) moveAround(npc *Character) {
	dir := rng.RollInt(1, 4)
	posX, posY := npc.Position()
	switch dir {
	case 1:
		posY += 1
	case 2:
		posX += 1
	case 3:
		posY -= 1
	case 4:
		posX -= 1
	}
	npc.SetDestPoint(posX, posY)
}

// saySomething sends random text on NPC chat channel.
func (ai *AI) saySomething(npc *Character) {
	if npc.Race() == nil {
		return
	}
	textID := fmt.Sprintf("random_chat_%s", npc.Race().ID())
	npc.AddChatMessage(textID)
}

// combatSkill selects NPC skill to use in combat or nil if specified
// NPC has no skills to use.
func (ai *AI) combatSkill(npc *Character, tar effect.Target) *skill.Skill {
	if len(npc.Skills()) < 1 {
		return nil
	}
	return npc.Skills()[0]
}
