package sources

import "strings"

func titleFromKeyword(keyword string) string {
	parts := strings.Fields(strings.TrimSpace(keyword))
	if len(parts) == 0 {
		return "Software Engineer"
	}

	for i := range parts {
		parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
	}

	return strings.Join(parts, " ") + " Engineer"
}
