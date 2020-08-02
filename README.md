![Build](https://github.com/flaambe/chat/workflows/Build/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/flaambe/chat/branch/master/graph/badge.svg)](https://codecov.io/gh/flaambe/chat)

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
