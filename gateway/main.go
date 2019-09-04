package main

import (
	"github.com/yamakiller/game/gateway/source"
	"github.com/yamakiller/magicNet/core/frame"
	"github.com/yamakiller/magicNet/core/launch"
)

func main() {

	launch.Launch(func() frame.Framework {
		return &source.Gateway{}
	})
}
