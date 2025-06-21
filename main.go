package main

import (
	"embed"
	"io/fs"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default: // linux, freebsd, openbsd, netbsd
		cmd = "xdg-open"
		args = []string{url}
	}
	exec.Command(cmd, args...).Start()
}

func waitForPort(port string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", port, 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

//go:embed public/*
var staticFiles embed.FS

func main() {
	gin.SetMode("release")
	app := gin.Default()
	publicFS, err := fs.Sub(staticFiles, "public")
	if err != nil {
		panic("出错了：" + err.Error())
	}
	app.NoRoute(func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(publicFS))
	})
	// 启动端口监听检测
	go func() {
		if waitForPort("localhost:38080", 10*time.Second) {
			openBrowser("http://localhost:38080")
		}
	}()
	app.Run(":38080")
}
