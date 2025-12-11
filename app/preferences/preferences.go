package preferences

import (
	"errors"
	"net/netip"
	"regexp"
	"runtime"
	"slices"
	"strconv"
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

func checkInitCommonAppPrefs(a fyne.App) {
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

	// ping.wait_duration
	checkInitDurationInAppPrefs(
		a, pingWaitDurationKey, pingWaitDurationDefVal, &PingWaitDurationMin, &PingWaitDurationMax,
	)
	PingWaitDuration = binding.BindPreferenceString(pingWaitDurationKey, a.Preferences())

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

	// fw.check.wait_duration
	checkInitDurationInAppPrefs(
		a, fwCheckWaitDurationKey, fwCheckWaitDurationDefVal, &FWCheckWaitDurationMin, &FWCheckWaitDurationMax,
	)
	FWCheckWaitDuration = binding.BindPreferenceString(fwCheckWaitDurationKey, a.Preferences())

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

type HostType int

const (
	HostTypeInvalid HostType = iota
	HostTypeIPv4
	HostTypeIPv6
	HostTypeDN
)

func DetectHostType(host string) (HostType, error) {
	if ip, err := netip.ParseAddr(host); err == nil {
		switch {
		case ip.Is4():
			return HostTypeIPv4, nil
		case ip.Is6():
			return HostTypeIPv6, nil
		}
	}
	if domainRegexp.MatchString(host) {
		return HostTypeDN, nil
	}

	return HostTypeInvalid, errors.New("невалидный адрес")
}

func HostValidator(v string) error {
	_, err := DetectHostType(v)
	return err
}

func MinMaxDurValidator(min, max time.Duration) func(v string) error {
	return func(v string) error {
		dur, err := time.ParseDuration(v)
		if err != nil {
			return errors.New("не является временным промежутком")
		}
		if dur < min ||
			dur > max {
			return errors.New("значение за пределами допустимого диапазона")
		}
		return nil
	}
}

func MinMaxIntValidator(min, max int) func(v string) error {
	return func(v string) error {
		d, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("не является целым числом")
		}
		if d < min ||
			d > max {
			return errors.New("значение за пределами допустимого диапазона")
		}
		return nil
	}
}

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

var (
	Protocols   = []string{string(TCP), string(UDP)}
	DefaultPort = map[Protocol]string{TCP: "443", UDP: "53"}
)

func setDefaultAllCommon() {
	SetDefaultAppearanceTheme()
	SetDefaultQueueWorkerNum()
	SetDefaultPingDefaultHost()
	SetDefaultPingWaitDuration()
	SetDefaultEICARMaxParallel()
	SetDefaultEICARWaitDuration()
	SetDefaultFWCheckWaitDuration()
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
	if err := AppearanceTheme.Set(appearanceThemeDefVal); err != nil {
		fyne.LogError("SetDefaultAppearanceTheme", err)
	}
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
	if err := QueueWorkerNum.Set(queueWorkerNumDefVal); err != nil {
		fyne.LogError("SetDefaultQueueWorkerNum", err)
	}
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

func SetDefaultPingDefaultHost() {
	if err := PingDefaultHost.Set(pingDefaultHostDefVal); err != nil {
		fyne.LogError("SetDefaultPingDefaultHost", err)
	}
}

// ping.wait_duration
const pingWaitDurationKey = "ping.wait_duration"

var (
	PingWaitDurationMin    = 1 * time.Second
	PingWaitDurationMax    = 10 * time.Second
	pingWaitDurationDefVal = 3 * time.Second
	PingWaitDuration       binding.String
)

func PingWaitDurationValidator(v string) error {
	return MinMaxDurValidator(PingWaitDurationMin, PingWaitDurationMax)(v)
}
func GetPingWaitDuration() time.Duration {
	strDur, err := PingWaitDuration.Get()
	if err != nil {
		return pingWaitDurationDefVal
	}
	dur, err := time.ParseDuration(strDur)
	if err != nil {
		return pingWaitDurationDefVal
	}
	return dur
}
func SetDefaultPingWaitDuration() {
	if err := PingWaitDuration.Set(pingWaitDurationDefVal.String()); err != nil {
		fyne.LogError("SetDefaultPingWaitDuration", err)
	}
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
	if err := EICARMaxParallel.Set(eicarMaxParallelDefVal); err != nil {
		fyne.LogError("SetDefaultEICARMaxParallel", err)
	}
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
	return MinMaxDurValidator(EICARWaitDurationMin, EICARWaitDurationMax)(v)
}
func GetEICARWaitDuration() time.Duration {
	strDur, err := EICARWaitDuration.Get()
	if err != nil {
		return eicarWaitDurationDefVal
	}
	dur, err := time.ParseDuration(strDur)
	if err != nil {
		return eicarWaitDurationDefVal
	}
	return dur
}
func SetDefaultEICARWaitDuration() {
	if err := EICARWaitDuration.Set(eicarWaitDurationDefVal.String()); err != nil {
		fyne.LogError("SetDefaultEICARWaitDuration", err)
	}
}

// fw.check.wait_duration
const fwCheckWaitDurationKey = "fw.check.wait_duration"

var (
	FWCheckWaitDurationMin    = 1 * time.Second
	FWCheckWaitDurationMax    = 10 * time.Second
	fwCheckWaitDurationDefVal = 3 * time.Second
	FWCheckWaitDuration       binding.String
)

func FWCheckWaitDurationValidator(v string) error {
	return MinMaxDurValidator(FWCheckWaitDurationMin, FWCheckWaitDurationMax)(v)
}
func GetFWCheckWaitDuration() time.Duration {
	strDur, err := FWCheckWaitDuration.Get()
	if err != nil {
		return fwCheckWaitDurationDefVal
	}
	dur, err := time.ParseDuration(strDur)
	if err != nil {
		return fwCheckWaitDurationDefVal
	}
	return dur
}
func SetDefaultFWCheckWaitDuration() {
	if err := FWCheckWaitDuration.Set(fwCheckWaitDurationDefVal.String()); err != nil {
		fyne.LogError("SetDefaultFWCheckWaitDuration", err)
	}
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
	if err := AVBinaries.Set(avBinariesDefVal); err != nil {
		fyne.LogError("SetDefaultAVBinaries", err)
	}
}
func SetDefaultAVFilePaths() {
	if err := AVFilePaths.Set(avFilePathsDefVal); err != nil {
		fyne.LogError("SetDefaultAVFilePaths", err)
	}
}
func SetDefaultFWBinaries() {
	if err := FWBinaries.Set(fwBinariesDefVal); err != nil {
		fyne.LogError("SetDefaultFWBinaries", err)
	}
}
func SetDefaultFWFilePaths() {
	if err := FWFilePaths.Set(fwFilePathsDefVal); err != nil {
		fyne.LogError("SetDefaultFWFilePaths", err)
	}
}
