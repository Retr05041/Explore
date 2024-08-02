OBJS = main.go

run: $(OBJS)
	if ! [ -d ./databases ]; then mkdir databases; fi
	go run $(OBJS) 

test: $(OBJS)
	if [ -d ./databases ]; then rm -r databases; fi
	if ! [ -d ./databases ]; then mkdir databases; fi
	go run $(OBJS) 

build: $(OBJS)
	if [ -d ./bin ]; then rm -r bin; fi
	if ! [ -d ./bin ]; then mkdir bin; fi
	go build -o ./bin/explore
	mkdir ./bin/databases
	cp -r ./maps ./bin
