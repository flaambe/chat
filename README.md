# Chat server
Simple chat server with **Go** and **MongoDB**

## Run & Environments
### Docker
Set environments
```bash
export MONGO_URI="mongodb://mongo:27017"
export DB_NAME="chat" 
export PORT=9000
```

Run docker-compose
```bash
docker-compose up
```
### Local
Set environments

```bash
export MONGO_URI="mongodb://localhost:27017"
export DB_NAME="chat" 
export PORT=9000
```

Run
```bash
make build
make run
```

## Test

```bash
make test
```
