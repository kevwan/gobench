package main

import (
	"time"

	"github.com/kevwan/gobench"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	config = `Name: bench
Mode: file
Encoding: json
Path: logs
Rotation: size
MaxSize: 1024`
	text = `The licenses for most software and other practical works are designed to take away your freedom to share and change the works. By contrast, the GNU General Public License is intended to guarantee your freedom to share and change all versions of a program--to make sure it remains free software for all its users. We, the Free Software Foundation, use the GNU General Public License for most of our software; it applies also to any other work released this way by its authors. You can apply it to your programs, too.`
)

func main() {
	var c logx.LogConf
	logx.Must(conf.LoadFromYamlBytes([]byte(config), &c))
	logx.MustSetup(c)

	b := gobench.NewBenchWithConfig(gobench.Config{
		Duration: time.Minute * 5,
	})
	b.Run(120000, func() {
		logx.Info(text)
	})
}
