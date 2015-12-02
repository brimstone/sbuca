package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func SetupTest() string {
	dir, err := ioutil.TempDir(os.TempDir(), "sbuca-")

	if err != nil {
		panic(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	return dir
}

func TearDown(dir string) {
	os.RemoveAll(dir)
}

func Test_main(t *testing.T) {
	defer TearDown(SetupTest())
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"sbuca"}
	main()
}

func Test_main_genkey(t *testing.T) {
	defer TearDown(SetupTest())
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"sbuca",
		"oneshot",
		"--key",
		"server.key",
		"--crt",
		"server.crt",
		"--ca",
		"ca.crt",
		"--host",
		"localhost:8600",
	}
	main()
}
