.ONESHELL:

clean:
	rm -f -r build

build: clean
	mkdir build
	cp -R client/public/. build
	cd server && go build -o ../build

run:
	./build/cdn