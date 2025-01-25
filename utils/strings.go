package utils

import "regexp"

func MatchUrl(url string, urlMap map[string][]string) ([]string, bool, string) {
	for pattern, allowedHosts := range urlMap {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(url); matches != nil {
			return allowedHosts, true, pattern
		}
	}
	return nil, false, ""
}
