elec = bin/elec
all: $(elec)

.DEFAULT: all

.PHONY: all

clean:
	rm -rf bin

$(elec): 
	go build -o $(elec) 