BlackHole
==================

# Go learning
	https://chai2010.cn/advanced-go-programming-book/ch3-asm/ch3-01-basic.html
	https://github.com/bingohuang/effective-go-zh-en

# VIM IDE
- [Macos](http://blog.xuezheyoushi.com/2017/09/07/Mac-OSXVim%E6%90%AD%E5%BB%BAGolang%E5%BC%80%E5%8F%91%E7%8E%AF%E5%A2%83)

	git clone https://github.com/golang/tools  
	go install golang.org/x/tools/cmd/guru
	
	git clone https://github.com/Go-zh/tools.git  
	go install github.com/Go-zh/tools/cmd/gopls

# Debug
- [Pprof](http://io.upyun.com/2018/01/21/debug-golang-application-with-pprof-and-flame-graph/)

# Depend
	go get -u -v github.com/kardianos/govendor

	mkdir -p $GOPATH/src/golang.org/x/
	cd $GOPATH/src/golang.org/x/
	git clone https://github.com/golang/tools.git
	go install golang.org/x/tools/cmd/goimports
