package tester

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pc_security_test/preferences"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func Ping(host string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	available := false

	cmd := exec.CommandContext(ctx, "ping", host)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	if err = cmd.Start(); err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(outPipe)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		if strings.Contains(line, "time=") {
			available = true
			_ = cmd.Process.Signal(os.Interrupt)
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	_ = cmd.Wait()

	return available, nil
}

func findBinariesPaths(binaries []string) map[string]string {
	res := make(map[string]string)

	for _, binary := range binaries {
		path, err := exec.LookPath(binary)
		if errors.Is(err, exec.ErrDot) {
			res[binary] = "Найден в текущем каталоге или подкаталогах"
		} else if err == nil {
			res[binary] = path
		}
	}

	return res
}

func findFilePaths(paths []string) map[string]struct{} {
	res := make(map[string]struct{})

	for _, path := range paths {
		_, err := os.Stat(path)
		if err == nil {
			res[path] = struct{}{}
		}
	}

	return res
}

type sourceType int

const (
	SourceTypeFW sourceType = iota
	SourceTypeAV
)

func FindBinariesAndPaths(st sourceType) map[string]string {
	binaries := []string{}
	paths := []string{}

	switch st {
	case SourceTypeFW:
		binaries, _ = preferences.FWBinaries.Get()
		paths, _ = preferences.FWFilePaths.Get()
	case SourceTypeAV:
		binaries, _ = preferences.AVBinaries.Get()
		paths, _ = preferences.AVFilePaths.Get()
	}

	binariesFound := findBinariesPaths(binaries)
	pathsFound := findFilePaths(paths)
	for path := range pathsFound {
		binariesFound[path] = ""
	}
	return binariesFound
}

var (
	eicarCounterMu         = &sync.RWMutex{}
	currParallelEICARTests = 0
)

func checkMaxParallelEICARTests() error {
	eicarCounterMu.RLock()
	defer eicarCounterMu.RUnlock()
	if m, _ := preferences.EICARMaxParallel.Get(); currParallelEICARTests > m {
		return errors.New("Дождитесь завершения предыдущего EICAR теста")
	}
	return nil
}

func incParallelEICARTests() {
	eicarCounterMu.Lock()
	defer eicarCounterMu.Unlock()

	currParallelEICARTests++
}

func decParallelEICARTests() {
	eicarCounterMu.Lock()
	defer eicarCounterMu.Unlock()

	currParallelEICARTests--
}

func EICARTest() (bool, error) {
	incParallelEICARTests()
	defer decParallelEICARTests()

	if err := checkMaxParallelEICARTests(); err != nil {
		return false, err
	}

	eicar := `X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`
	testFile := fmt.Sprintf("eicar-%s.txt", uuid.New().String())

	err := os.WriteFile(testFile, []byte(eicar), 0600)
	if err != nil {
		return false, errors.Wrap(err, "error creating file")
	}

	defer func() {
		if err := os.Remove(testFile); err != nil {
			fyne.LogError("error rm eicar file", err)
		}
	}()

	time.Sleep(preferences.GetEICARWaitDuration())

	_, err = os.Stat(testFile)
	if os.IsNotExist(err) {
		return true, nil
	} else if err != nil {
		return false, errors.Wrap(err, "error checking file")
	}

	return false, nil
}

func RemoveEICARs() {
	files, err := filepath.Glob("eicar-*.txt")
	if err != nil {
		fyne.LogError("error finding remaining eicar files", err)
		return
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			fyne.LogError("error rm eicar file", err)
		}
	}
}
