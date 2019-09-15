.PHONY: all clean debug

all:
	docker-compose up

clean:
	docker-compose down

debug:
	docker-compose -f docker-compose.yml -f docker-compose.debug.yml up
