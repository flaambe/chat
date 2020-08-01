![Build](https://github.com/flaambe/chat/workflows/Go/badge.svg?branch=master)

# Chat server
Simple chat server with **Go** and **MongoDB**

## Run & Environments

Set environments
```bash
export MONGO_URI="mongodb://mongo:27017"
export DB_NAME="chat" 
export PORT=9000
```

Run
```bash
docker-compose up
```
## Build

```bash
make build
```

## Test

```bash
make test
```
