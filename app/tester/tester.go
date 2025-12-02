package tester

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"pc_security_test/preferences"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	probing "github.com/prometheus-community/pro-bing"
)

func Ping(host string) (bool, error) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return false, errors.Wrap(err, "failed to create pinger")
	}

	pinger.Count = 3
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = pinger.RunWithContext(ctxTimeout)
	if err != nil {
		return false, errors.Wrap(err, "error running pinger")
	}

	return pinger.Statistics().PacketsRecv > 0, nil
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

	err := os.WriteFile(testFile, []byte(eicar), 0644)
	if err != nil {
		return false, errors.Wrap(err, "error creating file")
	}

	time.Sleep(5 * time.Second)

	_, err = os.Stat(testFile)
	if os.IsNotExist(err) {
		return true, nil
	} else if err != nil {
		return false, errors.Wrap(err, "error checking file")
	}
	_ = os.Remove(testFile)
	return false, nil
}
