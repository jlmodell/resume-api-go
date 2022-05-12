project_name = resume-api-go
image_name = resume-api-go:latest

run-local:
	go run .

requirements:
	go mod tidy

clean-packages:
	go clean -modcache

build:
	docker-compose up -d --build

rebuild:
	git pull
	make stop
	make remove
	make build

shell:
	docker exec -it $(project_name) bash

stop:
	docker stop $(project_name)

remove:
	docker rm $(project_name)

start:
	docker start $(project_name)

logs:
	docker logs $(project_name)
