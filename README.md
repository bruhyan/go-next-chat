## go-next-chat
chat application using go (backend) and next (frontend)

### Docker image
1) Pull postgres image
````
docker pull postgres:15-alpine
````
2) Run the image in desired port
````
docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine
````
OR

1) use the make command:
````
make postgresinit
````

### Create demo database
1) Exec into pg container and create database
````
docker exec -it postgres15 createdb --username=root --owner=root go-chat
````
2) check by going into container 
````
$ docker exec -it postgres15 psql
$ \l (list databases)
$ exit
````

### Managing migrations
1) Install golang-migrate if not installed:
````
brew install golang-migrate
````
2) create up/down migration scripts:
````
migrate create -ext sql -dir db/migrations add_users_table 
````
3) create table with the up migration script OR
````
migrate -path db/migrations -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose up 
````
use the make command:
````
make migrateup
````

### Start up
1) install deps
````
go mod tidy
````
2) start server
````
go run cmd/main.go