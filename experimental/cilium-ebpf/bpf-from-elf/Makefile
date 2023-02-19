all: main

main: clean
	go build main.go

run:
	sudo ./main ~/libbpf-bootstrap/examples/c/.output/bootstrap.bpf.o

clean:
	@rm -f main
