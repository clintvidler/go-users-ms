logs:
	docker-compose logs -f main

test: 
	docker exec users-ms go test -run ./... users-ms/data -v

db-prod:
	docker exec -it go-users-ms-db-prod-1 psql -U root datastore