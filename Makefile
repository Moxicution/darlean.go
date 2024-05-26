BUILDDIR=.\build
EXAMPLE-DIR1=.\examples\embed-c

.PHONY: build

build: build-utils build-base build-core build-embedlib build-examples-runner build-examples-embed-c

build-utils:
	@echo building utils
	@go -C .\utils\internal\main build 

build-base:
	@echo building base 
	@go -C .\base\internal\main build

build-core:
	@echo building core
	@go -C .\core\internal\main build

build-embedlib:
	@echo building embedlib
	@go -C .\embedlib build -o $(BUILDDIR)\embedlib\embedlib.dll -buildmode=c-shared types.go actorstub.go api.go embedlib.go	

build-examples-runner:
	@echo building examples-runner
	@go -C .\examples\runner build


build-examples-embed-c:
	@echo building examples-embed-c
	copy $(BUILDDIR)\embedlib\embedlib.dll .\examples\embed-c\embedlib.dll
	copy $(BUILDDIR)\embedlib\embedlib.h .\examples\embed-c\embedlib.h
#cd examples/embed-c
#gcc -o test test.c embedlib.dll
	gcc -o $(EXAMPLE-DIR1)\test.exe $(EXAMPLE-DIR1)\test.c $(EXAMPLE-DIR1)\embedlib.dll
#cd ../..