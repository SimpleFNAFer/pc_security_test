package preferences

import (
	"errors"
	"net"
	"regexp"
	"runtime"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

const (
	stringFallback = "UNDEF"
	intFallback    = -1
)

var (
	stringSliceFallback = []string{stringFallback}
	ErrorWrongArg       = errors.New("wrong arg")
)

func CheckInitAppPrefs(a fyne.App) {
	// appearance.theme
	checkInitStringInAppPrefs(
		a, appearanceThemeKey, appearanceThemeDefVal, AppearanceThemeLight, AppearanceThemeDark,
	)
	AppearanceTheme = binding.BindPreferenceString(appearanceThemeKey, a.Preferences())

	// queue.worker_num
	checkInitIntInAppPrefs(
		a, queueWorkerNumKey, queueWorkerNumDefVal, &QueueWorkerNumMin, &QueueWorkerNumMax,
	)
	QueueWorkerNum = binding.BindPreferenceInt(queueWorkerNumKey, a.Preferences())

	// ping.default_host
	checkInitStringInAppPrefs(a, pingDefaultHostKey, pingDefaultHostDefVal)
	PingDefaultHost = binding.BindPreferenceString(pingDefaultHostKey, a.Preferences())

	// eicar.max_parallel
	checkInitIntInAppPrefs(
		a, eicarMaxParallelKey, eicarMaxParallelDefVal, &EICARMaxParallelMin, &EICARMaxParallelMax,
	)
	EICARMaxParallel = binding.BindPreferenceInt(eicarMaxParallelKey, a.Preferences())

	// eicar.wait_duration
	checkInitDurationInAppPrefs(
		a, eicarWaitDurationKey, eicarWaitDurationDefVal, &EICARWaitDurationMin, &EICARWaitDurationMax,
	)
	EICARWaitDuration = binding.BindPreferenceString(eicarWaitDurationKey, a.Preferences())

	// av.binaries
	checkInitStringSliceInAppPrefs(
		a, avBinariesKey, avBinariesDefVal,
	)
	AVBinaries = binding.BindPreferenceStringList(avBinariesKey, a.Preferences())

	// av.filepaths
	checkInitStringSliceInAppPrefs(
		a, avFilePathsKey, avFilePathsDefVal,
	)
	AVFilePaths = binding.BindPreferenceStringList(avFilePathsKey, a.Preferences())

	// fw.binaries
	checkInitStringSliceInAppPrefs(
		a, fwBinariesKey, fwBinariesDefVal,
	)
	FWBinaries = binding.BindPreferenceStringList(fwBinariesKey, a.Preferences())

	// fw.filepaths
	checkInitStringSliceInAppPrefs(
		a, fwFilePathsKey, fwFilePathsDefVal,
	)
	FWFilePaths = binding.BindPreferenceStringList(fwFilePathsKey, a.Preferences())
}

func checkInitStringInAppPrefs(a fyne.App, k, defV string, avail ...string) {
	v := a.Preferences().StringWithFallback(k, stringFallback)
	if v == stringFallback || len(avail) != 0 && !slices.Contains(avail, v) {
		a.Preferences().SetString(k, defV)
	}
}
func checkInitIntInAppPrefs(a fyne.App, k string, defV int, min, max *int) {
	v := a.Preferences().IntWithFallback(k, intFallback)
	if v == intFallback || min != nil && v < *min || max != nil && v > *max {
		a.Preferences().SetInt(k, defV)
	}
}
func checkInitStringSliceInAppPrefs(a fyne.App, k string, defV []string) {
	v := a.Preferences().StringListWithFallback(k, stringSliceFallback)

	if slicesEqual(v, stringSliceFallback) {
		a.Preferences().SetStringList(k, defV)
	}
}
func checkInitDurationInAppPrefs(a fyne.App, k string, defV time.Duration, min, max *time.Duration) {
	v := a.Preferences().StringWithFallback(k, stringFallback)
	dur, err := time.ParseDuration(v)
	if v == stringFallback ||
		err != nil ||
		min != nil && dur < *min ||
		max != nil && dur > *max {
		a.Preferences().SetString(k, defV.String())
	}
}

func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func SetDefaultAll() {
	SetDefaultAppearanceTheme()
	SetDefaultQueueWorkerNum()
	SetDefaultPingDefaultHost()
	SetDefaultEICARMaxParallel()
	SetDefaultEICARWaitDuration()
	SetDefaultAVBinaries()
	SetDefaultAVFilePaths()
	SetDefaultFWBinaries()
	SetDefaultFWFilePaths()
}

// appearance.theme
const (
	appearanceThemeKey    = "appearance.theme"
	AppearanceThemeLight  = "light"
	AppearanceThemeDark   = "dark"
	AppearanceThemeSystem = "system"
	appearanceThemeDefVal = AppearanceThemeSystem
)

var AppearanceTheme binding.String

func AvailableAppearanceTheme() []string {
	return []string{AppearanceThemeLight, AppearanceThemeDark, AppearanceThemeSystem}
}
func SetDefaultAppearanceTheme() {
	_ = AppearanceTheme.Set(appearanceThemeDefVal)
}

// queue.worker_num
const queueWorkerNumKey = "queue.worker_num"

var (
	QueueWorkerNumMin    = 1
	queueWorkerNumDefVal = QueueWorkerNumMin
	QueueWorkerNumMax    = func() int {
		if runtime.NumCPU() > 16 {
			return 16
		}
		return runtime.NumCPU()
	}()
	QueueWorkerNum binding.Int
)

func SetDefaultQueueWorkerNum() {
	_ = QueueWorkerNum.Set(queueWorkerNumDefVal)
}

// ping.default_host
const (
	pingDefaultHostKey    = "ping.default_host"
	pingDefaultHostDefVal = "mail.ru"
)

var (
	domainRegexp    = regexp.MustCompile(`^(?i)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,63}$`)
	PingDefaultHost binding.String
)

func PingDefaultHostValidator(v string) error {
	if ip := net.ParseIP(v); ip != nil {
		return nil
	}
	if isDN := domainRegexp.MatchString(v); isDN {
		return nil
	}

	return errors.New("invalid ip or address")
}
func SetDefaultPingDefaultHost() {
	_ = PingDefaultHost.Set(pingDefaultHostDefVal)
}

// eicar.max_parallel
const eicarMaxParallelKey = "eicar.max_parallel"

var (
	EICARMaxParallelMin    = 1
	EICARMaxParallelMax    = 5
	eicarMaxParallelDefVal = EICARMaxParallelMin
	EICARMaxParallel       binding.Int
)

func SetDefaultEICARMaxParallel() {
	_ = EICARMaxParallel.Set(eicarMaxParallelDefVal)
}

// eicar.wait_duration
const eicarWaitDurationKey = "eicar.wait_duration"

var (
	EICARWaitDurationMin    = 5 * time.Second
	EICARWaitDurationMax    = 20 * time.Second
	eicarWaitDurationDefVal = 10 * time.Second
	EICARWaitDuration       binding.String
)

func EICARWaitDurationValidator(v string) error {
	dur, err := time.ParseDuration(v)
	if err != nil {
		return err
	}
	if dur < EICARWaitDurationMin ||
		dur > EICARWaitDurationMax {
		return errors.New("out of range")
	}
	return nil
}
func GetEICARWaitDuration() time.Duration {
	strDur, _ := EICARWaitDuration.Get()
	dur, err := time.ParseDuration(strDur)
	if err != nil {
		return eicarWaitDurationDefVal
	}
	return dur
}
func SetDefaultEICARWaitDuration() {
	_ = EICARWaitDuration.Set(eicarWaitDurationDefVal.String())
}

// av.binaries
// av.filepaths
// fw.binaries
// fw.filepaths
const (
	avBinariesKey  = "av.binaries"
	avFilePathsKey = "av.file_paths"
	fwBinariesKey  = "fw.binaries"
	fwFilePathsKey = "fw.file_paths"
)

var (
	AVBinaries  binding.StringList
	AVFilePaths binding.StringList
	FWBinaries  binding.StringList
	FWFilePaths binding.StringList
)

func SetDefaultAVBinaries() {
	_ = AVBinaries.Set(avBinariesDefVal)
}
func SetDefaultAVFilePaths() {
	_ = AVFilePaths.Set(avFilePathsDefVal)
}
func SetDefaultFWBinaries() {
	_ = FWBinaries.Set(fwBinariesDefVal)
}
func SetDefaultFWFilePaths() {
	_ = FWFilePaths.Set(fwFilePathsDefVal)
}
