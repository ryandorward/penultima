# penultima
A multiplayer browser-based open-world game inspired heavily from a classic from the 80s. 

React in the front, Go in the back. Postgres for the DB. Dockerized. 

I wanted to build this in order to become more proficient with the Go language, learn about concurrency, and getting better at using Docker containers.

It started out by modifying a basic Go+React chat app to move tiles around on the screen in the browser. This got out of hand, so I rebuilt it using https://github.com/floralbit/dungeon
as the framework. As such, this is essentially a fork of floralbit's excellent project. It also leans heavily on https://github.com/norendren/go-fov for the nice field-of-view calculations. 
Another resource I must credit is the tileset from the OG classic game.

This is very much a work-in-progress! It is working on my local environment, but I would not recommend forking this project in its current state!
