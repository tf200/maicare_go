package util

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Keep filenames predictable and safe for use in object keys.
	// Allowed: letters, digits, dash, underscore. Everything else becomes underscore.
	filenameBaseSanitizer = regexp.MustCompile(`[^a-zA-Z0-9-_]+`)
	underscoreCollapse    = regexp.MustCompile(`_+`)

	// Extension is sanitized separately: allow dot plus alnum (".pdf", ".tar.gz" isn't preserved by filepath.Ext anyway).
	filenameExtSanitizer = regexp.MustCompile(`[^a-zA-Z0-9.]+`)
)

// SanitizeFilename normalizes a user-provided filename into a safe ASCII-ish form.
// If maxLen <= 0, no length limit is applied.
// The returned string is never a path: separators are removed/replaced.
func SanitizeFilename(filename string, maxLen int) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return ""
	}

	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	base = filenameBaseSanitizer.ReplaceAllString(base, "_")
	base = underscoreCollapse.ReplaceAllString(base, "_")
	base = strings.Trim(base, "._- _")
	if base == "" {
		base = "file"
	}

	ext = filenameExtSanitizer.ReplaceAllString(ext, "")
	if ext == "." {
		ext = ""
	}

	out := base + ext
	if maxLen > 0 && len(out) > maxLen {
		// Prefer keeping the extension intact.
		keepExt := ext
		if len(keepExt) >= maxLen {
			keepExt = ""
		}
		maxBase := maxLen - len(keepExt)
		if maxBase < 1 {
			maxBase = maxLen
			keepExt = ""
		}
		baseTrunc := base
		if len(baseTrunc) > maxBase {
			baseTrunc = baseTrunc[:maxBase]
		}
		out = baseTrunc + keepExt
	}

	return out
}
