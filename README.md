# penultima
A multiplayer, browser-based, open-world game inspired heavily by the classic from the 80s. A game set in yesteryear, at a time when the world was 2-D, 8bit, and torroidal.

React in the front, Go in the back. Postgres for the DB. Websockets. Docker-ized.

It started out by modifying a basic Go+React chat app to move tiles around on the screen in the browser. This got out of hand, so I rebuilt it using https://github.com/floralbit/dungeon
as a framework for the backend. As such, the server side is essentially a fork of floralbit's excellent project. It also leans heavily on https://github.com/norendren/go-fov for the nice field-of-view calculations. Another resource I must credit is the tileset from the OG classic game.


This is very much a work-in-progress!
