package simc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type CLI struct {
	Executable string
}

var DefautlCLI CLI = CLI{
	Executable: "simc",
}

func (c CLI) Simulate(conf Configuration) (*Result, error) {

	hash := sha256.Sum256([]byte(conf))
	id := hex.EncodeToString(hash[:])

	// TODO: make this portable, as Mkfifo is only on linux and probably mac?
	confPath := filepath.Join(os.TempDir(), id+".simc")
	jsonPath := filepath.Join(os.TempDir(), id+".json")
	conf = Configuration(fmt.Sprintf("%s\njson2=%s", string(conf), jsonPath))

	if err := syscall.Mkfifo(confPath, 0666); err != nil {
		return nil, err
	}
	defer os.Remove(confPath)

	if err := syscall.Mkfifo(jsonPath, 0666); err != nil {
		return nil, err
	}
	defer os.Remove(jsonPath)

	confFile, err := os.OpenFile(confPath, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return nil, err
	}
	defer confFile.Close()

	jsonFile, err := os.OpenFile(jsonPath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	go func() {
		_, _ = confFile.WriteString(string(conf)) // we have opened the pipe first as a writer. this call blocks until simc finishes reading
		_ = confFile.Close()                      // close as early as we can
	}()

	decoder := json.NewDecoder(jsonFile)
	var data map[string]interface{} //TODO: FIX THIS, we need a schema... maybe use CUE?
	go func() {
		_ = decoder.Decode(&data) // we have opened the pipe first as a reader. this call will block until the file is finished writing to
		_ = jsonFile.Close()      // close early if we can
	}()

	defer confFile.Close()
	defer os.Remove(confFile.Name())

	cmd := exec.Command(c.Executable, confFile.Name()) // this will block until confFile is written to, and blocks until output json is fully consumed

	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	return &Result{
		Data: data,
	}, nil
}
