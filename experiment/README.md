# Cache Design Experiment

## Basic Feature
- Emulate data access for two functions fn1 and fn2 running in two different docker container
  - Local Hit: If data exist in it's store, it would just fetch it from there
  - Remote Hit: Else, it would try to access the data from the peer container

## Usage
- emulate invoking two functions
```sh
  docker compose up
```
- emulate invoking single function
```sh
  go run main.go  "http://localhost:<peer port>" "<serving port>"
```
