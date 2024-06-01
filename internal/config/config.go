// Package contains tools for parsing Agent and Server runtime configuration data

package config

import (
	"flag"
)

// Gets provided Command Line flags or and Enviroment Vars of configuration
func getProvidedFlags(visitFunc func(func(f *flag.Flag))) []string {
	res := make([]string, 0)
	visitFunc(func(f *flag.Flag) {
		res = append(res, f.Name)
	})
	return res
}

// Returns nil if parameter does not set, otherwise pointer to the parameter
func getParWithSetCheck[S any](val S, isSet bool) *S {
	if !isSet {
		return nil
	}
	return &val
}

// Sets dst value equal to src value if src is not nil, otherwise do nothing
func combineParameter[S any](dst *S, src *S) {
	if src == nil {
		return
	}
	*dst = *src
}
