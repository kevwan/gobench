package main

import (
	"github.com/kevwan/gobench"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

const config = `Name: bench
Mode: file
Encoding: json
Path: logs
Rotation: size
MaxSize: 10`

func main() {
	var c logx.LogConf
	logx.Must(conf.LoadFromYamlBytes([]byte(config), &c))
	logx.MustSetup(c)

	b := gobench.NewBench()
	b.Run(100000, func() {
		logx.Info("hello world")
	})
}
