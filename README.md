# BlackHole

# Depend
	go get -u -v github.com/kardianos/govendor

	mkdir -p $GOPATH/src/golang.org/x/
	cd $GOPATH/src/golang.org/x/
	git clone https://github.com/golang/tools.git
	go install golang.org/x/tools/cmd/goimports
