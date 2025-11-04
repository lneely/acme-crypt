install:V:
	go build -o $HOME/bin/CryptGet ./cmd/CryptGet
	go build -o $HOME/bin/CryptPut ./cmd/CryptPut

clean:V:
	rm -f $HOME/bin/CryptGet $HOME/bin/CryptPut