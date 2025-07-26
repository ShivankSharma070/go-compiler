run : build
	@bin/go-compiler

build :
	@go build -o bin/go-compiler
test:
	@go test -v -count=1
	@go test ./lexer -v -count=1
	
clean : 
	@rm -r bin
