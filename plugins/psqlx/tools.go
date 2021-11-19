package psqlx

import "strings"

func domainNormalize(name string) string {
	name = strings.ToLower(name)

	if strings.HasSuffix(name, ".") {
		name = strings.TrimSuffix(name, ".")
	}

	return name
}

func domainZoneNormalize(name string) string {
	if !strings.HasSuffix(name, ".") {
		name += "."
	}

	return name
}

func domainNameGuesses(name string) []string {
	var names []string
	var nameSplited string

	names = append(names, name)
	nameSplited = name

	for strings.Contains(nameSplited, ".") {
		result := strings.SplitN(nameSplited, ".", 2)

		nameSplited = result[1]
		names = append(names, result[1])
	}

	return names
}
