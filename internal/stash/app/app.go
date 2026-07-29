package app

import "BlackHole/internal/stash/service"

func Run() {
	service.Init()
	service.Run()
}
