package main

import (
	"runtime"

	. "github.com/kettek/gobl"
)

func main() {
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	builds := map[string][2]string{
		"build-server": {"server", "server" + exe},
		"build-client": {"client", "client" + exe},
	}
	for taskName, build := range builds {
		func(taskName string, build [2]string) {
			Task(taskName).
				Chdir(build[0]).
				Exec("go", "build", "-v", "-o", build[1])
		}(taskName, build)
	}
	var goroot string
	Task("build-client-wasm").
		Chdir("client").
		Exec("go", "env", "GOROOT").
		Result(func(r interface{}) {
			goroot = r.(string)[:len(r.(string))-1] + "/misc/wasm/wasm_exec.js"
		}).
		Exec("cp", &goroot, "../static/").
		Env("GOOS=js", "GOARCH=wasm").
		Exec("go", "build", "-v", "-o", "../static/client.wasm")

	Task("watch-server").
		Watch("server/*.go", "server/*/*.go").
		Signaler(SigQuit).
		Run("build-server").
		Run("run-server")
	Task("watch-client").
		Watch("client/*.go", "client/*/*.go", "net/*.go").
		Signaler(SigQuit).
		Run("build-client").
		Run("run-client")
	Task("watch-client-wasm").
		Watch("client/*.go", "client/*/*.go", "net/*.go").
		Signaler(SigQuit).
		Run("build-client-wasm")

	Task("run-server").
		Exec("./server/server"+exe, ":8080")
	Task("run-client").
		Exec("./client/client" + exe)

	Go()
}
