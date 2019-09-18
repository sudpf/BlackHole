# BlackHole

# vim IDE
	http://blog.xuezheyoushi.com/2017/09/07/Mac-OSXVim%E6%90%AD%E5%BB%BAGolang%E5%BC%80%E5%8F%91%E7%8E%AF%E5%A2%83/

# Depend
	go get -u -v github.com/kardianos/govendor

	mkdir -p $GOPATH/src/golang.org/x/
	cd $GOPATH/src/golang.org/x/
	git clone https://github.com/golang/tools.git
	go install golang.org/x/tools/cmd/goimports
