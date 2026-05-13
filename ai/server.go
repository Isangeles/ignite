/*
 * server.go
 *
 * Copyright 2022-2026 Dariusz Sikora <ds@isangeles.dev>
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

	"github.com/gorilla/websocket"

	"github.com/isangeles/fire/request"
	"github.com/isangeles/fire/response"
)

// Struct for server connection.
type Server struct {
	closed     bool
	conn       *websocket.Conn
	onResponse func(r response.Response)
}

// NewServer creates new server connection struct with connection
// to the server with specified host and port number.
// TLS switches between ws and wss protocols.
func NewServer(host, port string, tls bool) (*Server, error) {
	s := new(Server)
	protocol := "ws"
	if tls {
		protocol = "wss"
	}
	url := fmt.Sprintf("%s://%s:%s/", protocol, host, port)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to dial server: %v", err)
	}
	s.conn = conn
	go s.handleResponses()
	return s, nil
}

// Close closes server connection.
func (s *Server) Close() error {
	err := s.conn.Close()
	if err != nil {
		return fmt.Errorf("Unable to close server connection: %v",
			err)
	}
	s.closed = true
	return nil
}

// Closed checks if server connection was closed.
func (s *Server) Closed() bool {
	return s.closed
}

// Address returns server address.
func (s *Server) Address() string {
	return s.conn.RemoteAddr().String()
}

// SetOnServerResponseFunc sets function triggered on server reponse.
func (s *Server) SetOnResponseFunc(f func(r response.Response)) {
	s.onResponse = f
}

// Send sends specified request to the server.
// If error will occure while writing data using server connection
// then the server connection will be closed and error returned.
func (s *Server) Send(req request.Request) error {
	text, err := request.Marshal(&req)
	if err != nil {
		return fmt.Errorf("Unable to marshal request: %v", err)
	}
	err = s.conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		s.Close()
		return fmt.Errorf("Unable to write request: %v", err)
	}
	return nil
}

// handleResponses handles responses from the server and
// triggers onServerResponse for each response.
func (s *Server) handleResponses() {
	for !s.Closed() {
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			log.Printf("Server response: Unable to read from the server: %v", err)
			return
		}
		resp, err := response.Unmarshal(string(msg))
		if err != nil {
			log.Printf("Server response: Unable to unmarshal server response: %v",
				err)
			continue
		}
		if s.onResponse != nil {
			go s.onResponse(resp)
		}
		if resp.Closed {
			err := s.Close()
			if err != nil {
				log.Printf("Server response: unable to close connection: %v",
					err)
			}
			return
		}
	}
}
