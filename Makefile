run : build
	@bin/go-interpreter

build :
	@go build -o bin/go-interpreter
test:
	@go test -v -count=1
	@go test ./lexer -v -count=1
	
clean : 
	@rm -r bin
