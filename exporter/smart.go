package exporter

import (
	"regexp"
	"strconv"
	"strings"
)

type Smart struct {
	info  map[string]string
	attrs map[int]*SmartAttr
}

type SmartAttr struct {
	id         int64
	name       string
	flag       int64
	value      int64
	worst      int64
	thresh     int64
	attrType   string
	updated    string
	whenFailed string
	rawValue   int64
}

func (s *Smart) GetInfo(name ...string) string {
	for _, k := range name {
		val, ok := s.info[k]
		if ok {
			return val
		}
	}
	return ""
}

func (s *Smart) GetAttr(id ...int) *SmartAttr {
	for _, k := range id {
		attr, ok := s.attrs[k]
		if ok {
			return attr
		}
	}
	return &SmartAttr{}
}

func ParseSmart(s string) *Smart {
	smartInfo := map[string]string{}
	smartAttrs := map[int]*SmartAttr{}
	for _, section := range strings.Split(s, "\n\n") {
		if strings.Index(section, "=== START OF INFORMATION SECTION ===") > -1 {
			smartInfo = parseSmartInfo(section)
		} else if strings.Index(section, "=== START OF READ SMART DATA SECTION ===") > -1 {
			smartAttrs = parseSmartAttrs(section)
		}
	}
	return &Smart{
		info:  smartInfo,
		attrs: smartAttrs,
	}
}

func parseSmartInfo(s string) map[string]string {
	info := make(map[string]string)
	for _, line := range strings.Split(s, "\n") {
		kv := strings.Split(line, ": ")
		if len(kv) == 2 {
			info[trim(kv[0])] = trim(kv[1])
		}
	}
	return info
}

func parseSmartAttrs(s string) map[int]*SmartAttr {
	attrs := make(map[int]*SmartAttr)
	reSpaces := regexp.MustCompile(`\s+`)
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		// ID  ATTRIBUTE_NAME  FLAG  VALUE  WORST  THRESH  TYPE  UPDATED  WHEN_FAILED  RAW_VALUE
		vals := reSpaces.Split(trim(line), -1)
		if len(vals) < 10 {
			continue
		}
		id := parseInt(vals[0], 10)
		attrs[int(id)] = &SmartAttr{
			id:         id,
			name:       trim(vals[1]),
			flag:       parseInt(vals[2], 16),
			value:      parseInt(vals[3], 10),
			worst:      parseInt(vals[4], 10),
			thresh:     parseInt(vals[5], 10),
			attrType:   trim(vals[6]),
			updated:    trim(vals[7]),
			whenFailed: trim(vals[8]),
			rawValue:   parseInt(vals[9], 10),
		}
	}
	return attrs
}

func trim(s string) string {
	return strings.Trim(s, " \t")
}

func parseInt(s string, base int) int64 {
	if v, err := strconv.ParseInt(trim(s), base, 64); err == nil {
		return v
	}
	parts := strings.SplitN(trim(s), " ", 2)
	first := parts[0]
	if v, err := strconv.ParseInt(first, base, 64); err == nil {
		return v
	}
	return 0
}
