package utils

import (
	"path/filepath"
)

func DetermineFileType(filePath string) string {
	fileExtension := filepath.Ext(filePath)
	switch fileExtension {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".java":
		return "java"
	case ".c":
		return "c"
	case ".cpp", ".cc", ".cxx":
		return "cpp"
	case ".cs":
		return "csharp"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".rs":
		return "rust"
	case ".swift":
		return "swift"
	case ".kt", ".kts":
		return "kotlin"
	case ".sh":
		return "shell"
	case ".html":
		return "html"
	case ".css":
		return "css"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	case ".yaml", ".yml":
		return "yaml"
	case ".sql":
		return "sql"
	default:
		return "other"
	}
}
