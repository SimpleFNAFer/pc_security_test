package tester

import (
	"context"
	"os"
	"time"

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

func EICARTest() (bool, error) {
	eicar := `X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`
	testFile := "eicar.txt"

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
