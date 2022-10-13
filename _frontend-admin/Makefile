.PHONY: clean start

build/index.html: node_modules src/* src/*/* src/*/*/* .env.production.local
	npm run build

clean:
	rm --recursive --force build

start:
	npm start

node_modules:
	npm install
