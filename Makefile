OBJS = main.go

run : $(OBJS)
	if ! [ -d ./databases ]; then mkdir databases; fi
	go run $(OBJS) 
