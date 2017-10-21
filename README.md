## This repo presents a few basic concepts in multiplayer 2d games in browser ##
## Build ##
install deps
``` bash
go get github.com/gorilla/websockets
```
compile (written for golang 1.8+)
``` bash
go build
```

## Run ##
After building do
``` bash
./realtime-move-sync
```
run with -help to get list of runtime params. 

Notable runtime params:
 * **addr** - interface:port to be bound to
 * **collisions** - pass this flag in order to turn the collision detection on
 * **bots** - number of "bots" to be spawned, they will just circle around constantly
## Acknowledgments ##
hub.go and client.go files were Initially based on code from Gorilla Websockets chat example (Apache Foundation licensed).
