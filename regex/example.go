package regex

import (
	"regexp"
)

func splitLabelPrefixAndName(key string) (string, string) {
	prefix, name := "", ""
	cnt := 0
	for i := len(key) - 1; i >= 0; i-- {
		if key[i] == '/' {
			break
		}
		cnt++
	}
	if cnt == len(key) {
		name = key
	} else {
		prefix = key[:len(key)-cnt-1]
		name = key[len(key)-cnt:]
	}
	return prefix, name
}

func checkLabel(prefix string, name string, value string) bool {
	if len(prefix) > 253 || len(name) > 63 || len(name) <= 0 || len(value) > 63 {
		return false
	}
	var res string
	regexPrefix := regexp.MustCompile(`^[a-z0-9A-Z][a-z0-9A-Z\.]+[a-z0-9A-Z]${1,253}$`)
	regexName := regexp.MustCompile(`^[a-z0-9A-Z][a-z0-9A-Z_\-\.]+[a-z0-9A-Z]${1,63}$`)
	if prefix != "" {
		if prefix == "kubernetes.io" || prefix == "k8s.io" {
			return false
		}
		res = regexPrefix.FindString(prefix)
		if res == "" {
			return false
		}
	}
	res = regexName.FindString(name)
	if res == "" {
		return false
	}
	if value != "" {
		res = regexName.FindString(value)
		if res == "" {
			return false
		}
	}
	return true
}

var LabelPrefixSet = map[string]struct{}{
	"kubelet.kubernetes.io":                    {},
	"node.kubernetes.io":                       {},
	"beta.kubernetes.io/arch":                  {},
	"beta.kubernetes.io/instance-type":         {},
	"beta.kubernetes.io/os":                    {},
	"failure-domain.beta.kubernetes.io/region": {},
	"failure-domain.beta.kubernetes.io/zone":   {},
	"kubernetes.io/arch":                       {},
	"kubernetes.io/hostname":                   {},
	"kubernetes.io/os":                         {},
	"node.kubernetes.io/instance-type":         {},
	"topology.kubernetes.io/region":            {},
	"topology.kubernetes.io/zone":              {},
}
