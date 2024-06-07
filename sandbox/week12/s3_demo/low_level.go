package s3_demo

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// --- useful little functions ---

// firstN, startFrom = actualLimitOffset(limit, offset, listSize, 0, listSize)
func actualLimitOffset(limit, offset *int, dfltOffset, listSize int) (firstN, startFrom int) {
	startFrom = getOrDefaultPositiveIntValue(offset, dfltOffset)
	startFrom = min(startFrom, listSize)

	firstN = getOrDefaultPositiveIntValue(limit, listSize)
	// listSize >= startFrom + firstN
	firstN = min(firstN, listSize-startFrom)

	return
}

func getOrDefaultPositiveIntValue(v *int, dflt int) int {
	if v != nil && *v >= 0 {
		return *v
	} else {
		return dflt
	}
}

func loadIntFromMap(amap map[string]any, key string) (int, error) {
	vAny, vExists := amap[key]
	if !vExists {
		return 0, fmt.Errorf("loadIntFromMap failed, no key `%s` in map", key)
	}

	// userID = strconv.FormatUint(uint64(uidAny.(float64)), 36) // float64: json package to blame
	vFloat, isFloat := vAny.(float64)
	if !isFloat {
		return 0, fmt.Errorf("loadIntFromMap failed, value `%#v` is not a float", vAny)
	}

	return int(vFloat), nil
}

func loadStringFromMap(amap map[string]any, key string) (string, error) {
	vAny, vExists := amap[key]
	if !vExists {
		return "", fmt.Errorf("loadStringFromMap failed, no key `%s` in map", key)
	}

	vString, isString := vAny.(string)
	if !isString {
		return "", fmt.Errorf("loadStringFromMap failed, value `%#v` is not a string", vAny)
	}

	return vString, nil
}

var atomicCounter = new(atomic.Uint64)

func nextID_36() string {
	return strconv.FormatInt(int64(atomicCounter.Add(1)), 36)
}

func nextID_10() string {
	return strconv.FormatInt(int64(atomicCounter.Add(1)), 10)
}

func cutPrefix(s, prefix string) string {
	res, _ := strings.CutPrefix(s, prefix)
	return res
}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

func strRef(in string) *string {
	return &in
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const (
		RFC3339      = "2006-01-02T15:04:05Z07:00"
		RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	)
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		// line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}

// show_noop: do nothing, makes govet happy
func show_noop(msg string, xs ...any) {
	if len(msg)+len(xs) < 0 {
		fmt.Println("impossible")
	}
}
