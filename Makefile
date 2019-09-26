.PHONY: all build clean debug

all:
	docker-compose up server frontend

build:
	docker-compose up frontend-build

clean:
	docker-compose down

debug:
	docker-compose -f docker-compose.yml -f docker-compose.debug.yml up
