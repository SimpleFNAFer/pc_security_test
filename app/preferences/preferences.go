package preferences

import (
	"errors"
	"runtime"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

const (
	stringFallback = "UNDEF"
	intFallback    = -1
)

var (
	stringSliceFallback = []string{stringFallback}

	ErrorWrongArg = errors.New("wrong arg")
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
		fyne.CurrentApp().Preferences().SetString(k, defV)
	}
}
func checkInitIntInAppPrefs(a fyne.App, k string, defV int, min, max *int) {
	v := a.Preferences().IntWithFallback(k, intFallback)
	if v == intFallback || min != nil && v < *min || max != nil && v > *max {
		fyne.CurrentApp().Preferences().SetInt(k, defV)
	}
}
func checkInitStringSliceInAppPrefs(a fyne.App, k string, defV []string) {
	v := a.Preferences().StringListWithFallback(k, stringSliceFallback)

	if slicesEqual(v, stringSliceFallback) {
		fyne.CurrentApp().Preferences().SetStringList(k, defV)
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

// appearance.theme
const (
	appearanceThemeKey    = "appearance.theme"
	AppearanceThemeLight  = "light"
	AppearanceThemeDark   = "dark"
	appearanceThemeDefVal = AppearanceThemeDark
)

var AppearanceTheme binding.String

func AvailableAppearanceTheme() []string {
	return []string{AppearanceThemeLight, AppearanceThemeDark}
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

// ping.default_host
const (
	pingDefaultHostKey    = "ping.default_host"
	pingDefaultHostDefVal = "mail.ru"
)

var PingDefaultHost binding.String

// eicar.max_parallel
const eicarMaxParallelKey = "eicar.max_parallel"

var (
	EICARMaxParallelMin    = 1
	EICARMaxParallelMax    = 5
	eicarMaxParallelDefVal = EICARMaxParallelMin
	EICARMaxParallel       binding.Int
)

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
