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
	"fmt"
	"log"

	"github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/character"

	"github.com/isangeles/fire/request"
	"github.com/isangeles/fire/response"
)

// handleResponse handles specified response from Fire server.
func (g *Game) handleResponse(resp response.Response) {
	if !resp.Logon && g.onLoginFunc != nil {
		g.onLoginFunc(g)
	}
	g.handleUpdateResponse(resp.Update)
	for _, r := range resp.Character {
		g.handleCharacterResponse(r)
	}
	for _, r := range resp.Trade {
		err := g.handleTradeResponse(r)
		if err != nil {
			log.Printf("Game server: unable to handle trade response: %v",
				err)
		}
	}
	for _, r := range resp.Error {
		log.Printf("Game server error: %s", r)
	}
}

// handleUpdateRespone handles update response.
func (g *Game) handleUpdateResponse(resp response.Update) {
	res.Clear()
	g.Apply(resp.Module)
}

// handleCharacterResponse handles character response from the server.
func (g *Game) handleCharacterResponse(resp response.Character) {
	for _, c := range g.Characters() {
		if c.ID() == resp.ID && c.Serial() == resp.Serial {
			return
		}
	}
	char := g.Chapter().Character(resp.ID, resp.Serial)
	if char == nil {
		log.Printf("Game server: handle characher response: unable to find character in module: %s %s",
			resp.ID, resp.Serial)
	}
	g.AddCharacter(NewCharacter(char, g))
}

// handleTradeResponse handles trade response from the server.
func (g *Game) handleTradeResponse(resp response.Trade) error {
	// Find seller & buyer.
	object := g.Object(resp.SellerID, resp.SellerSerial)
	if object == nil {
		return fmt.Errorf("Seller not found: %s %s", resp.SellerID,
			resp.SellerSerial)
	}
	seller, ok := object.(*character.Character)
	if !ok {
		return fmt.Errorf("Seller is not a character: %s %s", resp.SellerID,
			resp.SellerSerial)
	}
	object = g.Object(resp.BuyerID, resp.BuyerSerial)
	if object == nil {
		return fmt.Errorf("Buyer not found: %s %s", resp.BuyerID,
			resp.BuyerSerial)
	}
	buyer, ok := object.(*character.Character)
	if !ok {
		return fmt.Errorf("Buyer is not a character: %s %s", resp.BuyerID,
			resp.BuyerSerial)
	}
	// Validate trade.
	buyValue := 0
	for id, serials := range resp.ItemsBuy {
		for _, serial := range serials {
			it := seller.Inventory().Item(id, serial)
			if it != nil {
				buyValue += it.Value()
			}
		}
	}
	sellValue := 0
	for id, serials := range resp.ItemsSell {
		for _, serial := range serials {
			it := buyer.Inventory().Item(id, serial)
			if it != nil {
				sellValue += it.Value()
			}
		}
	}
	if sellValue < buyValue {
		return nil
	}
	// Send accept request.
	req := request.Request{Accept: []int{resp.ID}}
	err := g.Server().Send(req)
	if err != nil {
		return fmt.Errorf("Unable to send accept request: %v", err)
	}
	return nil
}
