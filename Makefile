.PHONY: all build clean debug

all:
	docker-compose up server frontend

build:
	docker-compose -f docker-compose.build.yml up frontend-build

clean:
	docker-compose -f docker-compose.yml -f docker-compose.debug.yml -f docker-compose.build.yml down

debug:
	docker-compose -f docker-compose.yml -f docker-compose.debug.yml up
