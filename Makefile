CC=go
FLAGS=

progname=guard
bin=dist
target=$(bin)/$(progname)

main=cmd/guard/main.go

src=$(main)

all: $(target)

debug: $(target)

$(target): $(bin)
	$(CC) build -o $(target) $(FLAGS) $(src)

$(build):
	mkdir -p $(build)

$(bin):
	mkdir -p $(bin)

clean:
	rm -rf $(build) $(bin)

.PHONY: all clean
