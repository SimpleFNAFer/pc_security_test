//go:build windows

package preferences

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

func CheckInitAppPrefs(a fyne.App) {
	checkInitCommonAppPrefs(a)

	firstRun := a.Preferences().Bool(firstRunKey)

	// av.search_default
	checkInitBoolInAppPrefs(a, avSearchDefaultKey, AVSearchDefaultDefVal, firstRun)
	AVSearchDefault = binding.BindPreferenceBool(avSearchDefaultKey, a.Preferences())

	// fw.search_default
	checkInitBoolInAppPrefs(a, fwSearchDefaultKey, FWSearchDefaultDefVal, firstRun)
	FWSearchDefault = binding.BindPreferenceBool(fwSearchDefaultKey, a.Preferences())

	if !firstRun {
		a.Preferences().SetBool(firstRunKey, true)
	}
}

func SetDefaultAll() {
	setDefaultAllCommon()
	SetDefaultAVSearchDefault()
	SetDefaultFWSearchDefault()
}

// firstrun
const firstRunKey = "firstrun"

func checkInitBoolInAppPrefs(a fyne.App, k string, defV, firstRun bool) {
	if !firstRun {
		a.Preferences().SetBool(k, defV)
	}
}

// av.search_default
const avSearchDefaultKey = "av.search_default"

var (
	AVSearchDefaultDefVal = true
	AVSearchDefault       binding.Bool
)

func SetDefaultAVSearchDefault() {
	if err := AVSearchDefault.Set(AVSearchDefaultDefVal); err != nil {
		fyne.LogError("SetDefaultAVSearchDefault", err)
	}
}

// fw.search_default
const fwSearchDefaultKey = "fw.search_default"

var (
	FWSearchDefaultDefVal = true
	FWSearchDefault       binding.Bool
)

func SetDefaultFWSearchDefault() {
	if err := FWSearchDefault.Set(FWSearchDefaultDefVal); err != nil {
		fyne.LogError("SetDefaultFWSearchDefault", err)
	}
}

var (
	avBinariesDefVal  = []string{}
	avFilePathsDefVal = []string{}
	fwBinariesDefVal  = []string{}
	fwFilePathsDefVal = []string{}
)
