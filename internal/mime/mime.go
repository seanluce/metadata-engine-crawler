package mime

import (
	"mime"
	"strings"
)

// fallback map for common extensions not always registered by the OS
var fallback = map[string]string{
	".md":    "text/markdown",
	".ts":    "application/typescript",
	".tsx":   "application/typescript",
	".jsx":   "text/jsx",
	".go":    "text/x-go",
	".rs":    "text/x-rust",
	".yaml":  "application/yaml",
	".yml":   "application/yaml",
	".toml":  "application/toml",
	".sh":    "application/x-sh",
	".bash":  "application/x-sh",
	".zsh":   "application/x-sh",
	".csv":   "text/csv",
	".svg":   "image/svg+xml",
	".webp":  "image/webp",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".avif":  "image/avif",
}

// TypeByExtension returns a MIME type for the given file extension.
// Extension must include the leading dot (e.g. ".pdf").
func TypeByExtension(ext string) string {
	if ext == "" {
		return "application/octet-stream"
	}
	lower := strings.ToLower(ext)
	if t, ok := fallback[lower]; ok {
		return t
	}
	if t := mime.TypeByExtension(lower); t != "" {
		// strip parameters like "; charset=utf-8"
		if idx := strings.Index(t, ";"); idx != -1 {
			return strings.TrimSpace(t[:idx])
		}
		return t
	}
	return "application/octet-stream"
}
