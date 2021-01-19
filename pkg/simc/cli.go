package simc

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/google/uuid"
)

type CLI struct {
	Executable string
}

var DefautlCLI CLI = CLI{
	Executable: "simc",
}

func (c CLI) Version() (Version, error) {
	cmd := exec.Command("simc")
	stdout, err := cmd.Output()
	if stdout == nil && err != nil {
		return "", err
	}
	out := string(stdout)
	versionLine := strings.Split(out, "\n")[0]
	return Version(versionLine), nil
}

func (c CLI) Simulate(conf Configuration) (*Result, error) {
	// id := uuid.NewHash(sha256.New(), )
	// hash := sha256.Sum256([]byte(conf))
	id := uuid.New().String()

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

	errCh := make(chan error)

	go func() {
		confFile, err := os.OpenFile(confPath, os.O_WRONLY|os.O_APPEND, 0)
		if err != nil {
			errCh <- err
			return
		}
		defer os.Remove(confFile.Name())
		defer confFile.Close()

		_, err = confFile.WriteString(string(conf)) // we have opened the pipe first as a writer. this call blocks until simc finishes reading
		if err != nil {
			errCh <- err
		}
	}()

	dataCh := make(chan map[string]interface{})
	go func() {
		jsonFile, err := os.OpenFile(jsonPath, os.O_RDONLY, 0)
		if err != nil {
			errCh <- err
			return
		}
		defer jsonFile.Close()
		decoder := json.NewDecoder(jsonFile)
		var data map[string]interface{} //TODO: FIX THIS, we need a schema... maybe use CUE?
		err = decoder.Decode(&data)     // we have opened the pipe first as a reader. this call will block until the file is finished writing to
		if err != nil {
			errCh <- err
		}
		dataCh <- data
	}()

	cmd := exec.Command(c.Executable, confPath) // this will block until confFile is written to, and blocks until output json is fully consumed

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	select {
	case data := <-dataCh:
		return &Result{
			Data: data,
		}, nil
	case err := <-errCh:
		close(dataCh)
		return nil, err
	}

}
