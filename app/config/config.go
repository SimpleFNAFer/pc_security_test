package config

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	mu      = &sync.Mutex{}
	cfgPath = ""
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get caller information")
	}
	cfgPath = filepath.Join(filepath.Dir(filename), "config.yaml")
}

func readCfgFile() map[string]any {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil
	}

	var cfg map[string]any
	yaml.Unmarshal(data, &cfg)

	return cfg
}

func get[T any](path string, defVal T) T {
	cfg := readCfgFile()
	if cfg == nil {
		return defVal
	}

	pathSplit := strings.Split(path, ".")
	currentCfg := cfg

	for i := range pathSplit {
		if i == len(pathSplit)-1 {
			return getTypedValue(currentCfg[pathSplit[i]], defVal)
		}

		newCfgInterface, ok := currentCfg[pathSplit[i]]
		if !ok {
			return defVal
		}
		newCfg, ok := newCfgInterface.(map[string]any)
		if !ok {
			return defVal
		}

		currentCfg = newCfg
	}

	return defVal
}

func getTypedValue[T any](i any, defVal T) T {
	iVal := reflect.ValueOf(i)
	defValType := reflect.TypeOf(defVal)

	if iVal.Kind() == reflect.Slice && defValType.Kind() == reflect.Slice {
		return getSliceTypedValue(i, defVal)
	}

	if val, ok := i.(T); ok {
		return val
	}
	return defVal
}

func getSliceTypedValue[T any](i any, defVal T) T {
	iVal := reflect.ValueOf(i)
	defValVal := reflect.ValueOf(defVal)

	if iVal.Kind() != reflect.Slice || defValVal.Kind() != reflect.Slice {
		return defVal
	}

	elemType := defValVal.Type().Elem()
	result := reflect.MakeSlice(defValVal.Type(), 0, iVal.Len())

	for j := 0; j < iVal.Len(); j++ {
		elem := iVal.Index(j)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}
		if elem.Type().AssignableTo(elemType) {
			result = reflect.Append(result, elem)
		} else {
			return defVal
		}
	}

	return result.Interface().(T)
}

type StringSlice struct {
	mu      sync.Mutex
	cfgPath string
	value   []string
}

func (s *StringSlice) Get() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = get(s.cfgPath, s.value)
	return s.value
}

func NewStringSlice(cfgPath string, defVal []string) *StringSlice {
	return &StringSlice{mu: sync.Mutex{}, cfgPath: cfgPath, value: get(cfgPath, defVal)}
}

type Int struct {
	mu      sync.Mutex
	cfgPath string
	value   int
}

func (i *Int) Get() int {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.value = get(i.cfgPath, i.value)
	return i.value
}

func NewInt(cfgPath string, defVal int) *Int {
	return &Int{mu: sync.Mutex{}, cfgPath: cfgPath, value: get(cfgPath, defVal)}
}
