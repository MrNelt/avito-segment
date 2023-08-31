ROOT_FOLDER=.

# Commands
env:  ##@Environment Create .env file with variables
	@$(eval SHELL:=/bin/bash)
	@cp .env.example .env

db:  ##@Database Create database with docker-compose
	docker-compose -f docker-compose.yml up postgres -d --remove-orphans

run:  ##@Application Run application server
	docker-compose -f docker-compose.yml up --build api -d

stop: ##@Application Stop application server
	docker-compose -f docker-compose.yml stop

down:  ##@Application Down and clean application server
	docker-compose -f docker-compose.yml down

# open_db:  ##@Database Open database console inside docker-image
# 	docker exec -it postgres psql -d $(POSTGRES_DB) -U $(POSTGRES_USER)

clean_db: ##@Database Stop database and clean its data
	docker stop postgres && docker rm postgres

test:  ##@Testing Test application with gotest (be careful: tests drops all data in database from env)
	make db && go test $(ROOT_FOLDER)/...

lint:
	golangci-lint run --issues-exit-code 0 --print-issued-lines=false --out-format code-climate:gl-code-quality-report.json,line-number $(ROOT_FOLDER)/pkg/...

format:
	go fmt $(ROOT_FOLDER)/pkg/...
