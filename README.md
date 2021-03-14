# Platformer Maker Godot Server
#### Mario Maker clone built in Godot

## Requirements

- Go version 1.whateverVersionAddedGoModules
- Either:
    - RethinkDB OR 
    - Docker + Docker Compose

I would suggest grabbing docker to to handle spinning up your RethinkDB instance.

## Run (Development)

From terminal / command line:

```bash
    cd Platformer-Maker-Godot-Server
    # USE ONLY THE COMMANDS YOU NEED; 
    # Im just listing them all here

    # start your rethinkdb instance
    docker-compose up -d rethinkdb

    # start your server instance
    docker-compose up -d pm_server

    # start RethinkDB + server instance
    docker-compose up -d

    # STOPPING THE SERVICES
    docker-compose stop rethinkdb
    docker-compose stop pm_server
    docker-compose stop
```

## Debugging

The project comes with a `launch.json` setup for all the default settings in the `.env` file. 

Therefore, it's recommended you use [Visual Studio Code](https://code.visualstudio.com/Download) with the Go Extensions installed and enabled.

From there, you can press F5 to debug.