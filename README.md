## Introduction
Ignite is a [Flame](https://github.com/isangeles/flame) AI client program for the [Fire](https://github.com/isangeles/fire) game server.

The program connects to the game server and controls game NPCs(Non-Player Characters).

Currently in a early development stage.
## Build & Run
Get sources from git:
```
go get -u github.com/isangeles/ignite
```
Build executable:
```
go build github.com/isangeles/ignite
```
Configure host address and port of the game server and user credentials in `.ignite` file(create if it doesn't already exist):
```
server:[host];[port]
user:[id];[password]
```
Run program:
```
./ignite
```
After this, the program should establish a connection with the game server and control game characters assigned to the AI user by the server.
## Configuration
Configuration is stored in `.ignite` file placed in the program executable directory.
### Configuration values:
```
server:[address];[port]
```
Value for game server host and port.
```
user:[user ID];[password]
```
Value for game server user ID and password.
```
move-freq:[milliseconds]
```
Value for AI random move frequency in milliseconds, 3000 by default.
```
char-freq:[milliseconds]
```
Value for AI random chat frequency in milliseconds, 5000 by default.
## Documentation
Source code documentation could be easily browsed with the `go doc` command.

Besides that `doc` directory contains some useful documentation pages.

Documentation pages are in Troff format and could be easily displayed with `man` command.

For example to display documentation page for the AI configuration:
```
man doc/config
```
## Contributing
You are welcome to contribute to project development.

If you looking for things to do, then check the TODO file or contact maintainer(dev@isangeles.pl).

When you find something to do, create a new branch for your feature.
After you finish, open a pull request to merge your changes with master branch.
## Contact
* Isangeles <<dev@isangeles.pl>>
## License
Copyright 2021 Dariusz Sikora <<dev@isangeles.pl>>

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
MA 02110-1301, USA.
