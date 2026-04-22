package icon

import (
	"fmt"
	"strings"
)

func Format(entry IconEntry, cfg ConfigSection, fallbackName string) string {
	icon := entry.Icon
	name := entry.DisplayName

	if entry.DisplayName == "" {
		icon = cfg.FallbackIcon
		args := strings.Fields(fallbackName)
		if len(args) > 0 {
			name = args[0]
		} else {
			name = fallbackName
		}
	}

	if cfg.ShowName && name != "" {
		if icon != "" {
			return fmt.Sprintf("%s %s", icon, name)
		}
		return name
	}

	if icon != "" {
		return icon
	}

	return name
}
