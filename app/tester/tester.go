package tester

import (
	"bufio"
	"context"
	"fmt"
	"net"
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
	ctx, cancel := context.WithTimeout(context.Background(), preferences.GetPingWaitDuration())
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
		if strings.Contains(line, "time=") {
			available = true
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				fyne.LogError("Ping.cmd.Process.Signal(os.Interrupt)", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	if err := cmd.Wait(); err != nil {
		fyne.LogError("Ping.cmd.Wait", err)
	}

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
	var err error
	binaries := []string{}
	paths := []string{}

	switch st {
	case SourceTypeFW:
		if binaries, err = preferences.FWBinaries.Get(); err != nil {
			fyne.LogError("FindBinariesAndPaths.FWBinaries", err)
		}
		if paths, err = preferences.FWFilePaths.Get(); err != nil {
			fyne.LogError("FindBinariesAndPaths.FWFilePaths", err)
		}

	case SourceTypeAV:
		if binaries, err = preferences.AVBinaries.Get(); err != nil {
			fyne.LogError("FindBinariesAndPaths.AVBinaries", err)
		}
		if paths, err = preferences.AVFilePaths.Get(); err != nil {
			fyne.LogError("FindBinariesAndPaths.AVFilePaths", err)
		}
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
	m, err := preferences.EICARMaxParallel.Get()
	if err != nil {
		fyne.LogError("checkMaxParallelEICARTests.EICARMaxParallel.Get", err)
	}
	if currParallelEICARTests > m {
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

func FWTest(host, port string, proto preferences.Protocol) (bool, error) {
	if port == "" {
		port = preferences.DefaultPort[proto]
	}

	hostType, err := preferences.DetectHostType(host)
	if err != nil {
		return false, err
	}
	if hostType == preferences.HostTypeIPv6 {
		proto += "6"
	}

	var addr strings.Builder
	addr.WriteString(host)
	addr.WriteString(":" + port)

	conn, err := net.DialTimeout(string(proto), addr.String(), preferences.GetFWCheckWaitDuration())
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fyne.LogError("FWTest.conn.Close", err)
		}
	}()

	return false, nil
}
