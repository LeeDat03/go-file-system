package main

import "strings"

func IsAbsolutePath(path string) bool {
	if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "root/") {
		return true
	}

	return false
}
