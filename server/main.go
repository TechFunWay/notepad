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
	dataDir := flag.String("data-dir", "", "数据目录 (默认 ./data)")
	webDir := flag.String("web-dir", "", "前端静态文件目录")
	uploadDir := flag.String("upload-dir", "", "上传文件目录")
	shareDirs := flag.String("share-dirs", "", "共享目录列表，冒号分隔")
	flag.Parse()

	cmd.StartServer(*port, *dataDir, *webDir, *uploadDir, *shareDirs)
}
