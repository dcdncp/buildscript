# buildscript

A small, experimental programming language created as a toy project. This project was intended as a replacement for cmake, but more for fun. The language currently has a lot of problems and bugs. The standard library has very little functionality and is implemented only for one specific task.

## Features
- Basic syntax and control structures
- Basic data types: integers, strings, booleans
- Simple arithmetic and logical operations
- Exceptions, modules, etc.

## Installation
To run the language, clone this repository and follow the instructions below:

```bash
# Clone this repository
git clone https://github.com/dcdncp/buildscript.git
cd bscript

# Compile (or set up) the language
go build
```

## Usage
Run the interpreter or compiler with a sample file:

```bash
./bscript example.bs
```

By default runs build.script file:

```bash
# equals to ./bscript build.script
./bscript
```

## Examples
Here’s a simple program in the buildscript:

```bscript
// Simple make file
const make = import("std/make")

make.set_compiler("clang")

make.src("src/main.c")
make.src("src/test.c", "src/test.h")

make.target("app")
make.build()
```

This program creates Makefile:

```Makefile
app.exe: src/main.c.o src/test.c.o 
	gcc -o app.exe src/main.c.o src/test.c.o 

src/test.c.o: src/test.c src/test.h 
	gcc -c -o src/test.c.o src/test.c 

src/main.c.o: src/main.c 
	gcc -c -o src/main.c.o src/main.c 
```

## Roadmap
- [ ] Fix all bugs
- [ ] Implement classes
- [ ] Better error handling
- [ ] Extend standart library

## License
This project is licensed under the MIT License.
