package roles

import (
	"strings"
)

type Permission struct {
	Name        string   `yaml:"name"`
	Resource    string   `yaml:"resource"`
	Method      []string `yaml:"method"`
	GroupKey    string   `yaml:"group_key"`
	SectionKey  string   `yaml:"section_key"`
	DisplayName string   `yaml:"display_name"`
	Description string   `yaml:"description"`
	SortOrder   int      `yaml:"sort_order"`
}

type Role struct {
	Name        string   `yaml:"name"`
	ID          int      `yaml:"id"`
	Description string   `yaml:"description"`
	Permissions []string `yaml:"permissions"`
}

type Config struct {
	Permissions []Permission `yaml:"permissions"`
	Roles       []Role       `yaml:"roles"`
}

type PermissionMetadata struct {
	GroupKey    string
	SectionKey  string
	DisplayName string
	Description *string
	SortOrder   int32
}

func (p Permission) Normalize() PermissionMetadata {
	groupKey := normalizeKey(p.GroupKey)
	sectionKey := normalizeKey(p.SectionKey)

	parts := splitPermissionName(p.Name)
	if groupKey == "" {
		groupKey = deriveGroupKey(parts)
	}
	if sectionKey == "" {
		sectionKey = deriveSectionKey(parts)
	}

	displayName := strings.TrimSpace(p.DisplayName)
	if displayName == "" {
		displayName = deriveDisplayName(parts)
	}

	var description *string
	if trimmed := strings.TrimSpace(p.Description); trimmed != "" {
		description = &trimmed
	}

	return PermissionMetadata{
		GroupKey:    groupKey,
		SectionKey:  sectionKey,
		DisplayName: displayName,
		Description: description,
		SortOrder:   int32(p.SortOrder),
	}
}

func splitPermissionName(name string) []string {
	rawParts := strings.Split(strings.TrimSpace(name), ".")
	parts := make([]string, 0, len(rawParts))
	for _, part := range rawParts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func deriveGroupKey(parts []string) string {
	if len(parts) == 0 {
		return "general"
	}
	return normalizeKey(parts[0])
}

func deriveSectionKey(parts []string) string {
	if len(parts) <= 2 {
		return "general"
	}
	return normalizeKey(strings.Join(parts[1:len(parts)-1], "_"))
}

func deriveDisplayName(parts []string) string {
	if len(parts) == 0 {
		return "Unnamed Permission"
	}
	if len(parts) == 1 {
		return humanize(parts[0])
	}

	action := humanize(parts[len(parts)-1])
	targetParts := parts[:len(parts)-1]
	if len(targetParts) > 1 {
		targetParts = targetParts[1:]
	}

	target := "General"
	if len(targetParts) > 0 {
		target = humanize(strings.Join(targetParts, " "))
	}

	return action + " " + target
}

func normalizeKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	replacer := strings.NewReplacer(".", "_", "-", "_", " ", "_")
	value = replacer.Replace(value)
	for strings.Contains(value, "__") {
		value = strings.ReplaceAll(value, "__", "_")
	}
	return strings.Trim(value, "_")
}

func humanize(value string) string {
	if value == "" {
		return ""
	}

	value = strings.ReplaceAll(value, ".", " ")
	value = strings.ReplaceAll(value, "_", " ")
	words := strings.Fields(strings.ToLower(value))
	for i, word := range words {
		words[i] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, " ")
}
