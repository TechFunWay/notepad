package main

import (
	"flag"
	"os"
	"strings"

	"notepad/cmd"
)

var (
	Version   = "v1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	cmd.SetVersionInfo(Version, BuildTime, GitCommit)

	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		cmd.ExecuteCLI(os.Args[1:])
		return
	}

	port := flag.Int("port", 0, "服务端口 (默认 8904)")
	data := flag.String("data", "", "数据目录 (默认 ./data)")
	flag.Parse()

	cmd.StartServer(*port, *data)
}
