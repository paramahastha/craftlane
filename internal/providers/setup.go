package providers

import (
	"fmt"
	"os"
	"path/filepath"
)

func MarkedHook(provider string) map[string]any {
	return map[string]any{"_craftlane": "observability-v1", "provider": provider, "command": "craftlane observe bridge --provider " + provider}
}

func PreviewFile(path, provider string, warning string, trust bool) (SetupPreview, error) {
	object, err := readObject(path)
	if err != nil {
		return SetupPreview{}, err
	}
	preview := SetupPreview{Provider: provider, Warnings: []string{warning}, RequiresTrust: trust}
	if hasMarkedHook(object, provider) {
		preview.AlreadyApplied = true
		preview.Changes = []SetupChange{{Path: path, Action: "retain", Detail: "Craftlane observation hook already present"}}
	} else {
		preview.Changes = []SetupChange{{Path: path, Action: "merge", Detail: "Add one marked observation hook; preserve unrelated settings"}}
	}
	return preview, nil
}

func ApplyFile(path, provider string) (bool, error) {
	object, err := readObject(path)
	if err != nil {
		return false, err
	}
	if hasMarkedHook(object, provider) {
		return false, nil
	}
	hooks, _ := object["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
	}
	items, _ := hooks["SessionStart"].([]any)
	hooks["SessionStart"] = append(items, MarkedHook(provider))
	object["hooks"] = hooks
	return true, writeObject(path, object)
}

func RemoveFile(path, provider string) (bool, error) {
	object, err := readObject(path)
	if err != nil {
		return false, err
	}
	changed := false
	if hooks, ok := object["hooks"].(map[string]any); ok {
		for event, raw := range hooks {
			items, ok := raw.([]any)
			if !ok {
				continue
			}
			kept := items[:0]
			for _, item := range items {
				marked, _ := item.(map[string]any)
				if marked["_craftlane"] == "observability-v1" && marked["provider"] == provider {
					changed = true
					continue
				}
				kept = append(kept, item)
			}
			if len(kept) == 0 {
				delete(hooks, event)
			} else {
				hooks[event] = kept
			}
		}
		if len(hooks) == 0 {
			delete(object, "hooks")
		}
	}
	if !changed {
		return false, nil
	}
	if err := deleteIfEmpty(path, object); err != nil {
		return false, err
	}
	if len(object) > 0 {
		return true, writeObject(path, object)
	}
	return true, nil
}

func hasMarkedHook(object map[string]any, provider string) bool {
	hooks, _ := object["hooks"].(map[string]any)
	for _, raw := range hooks {
		items, _ := raw.([]any)
		for _, item := range items {
			marked, _ := item.(map[string]any)
			if marked["_craftlane"] == "observability-v1" && marked["provider"] == provider {
				return true
			}
		}
	}
	return false
}

func ensurePath(root, relative string) string { return filepath.Join(root, relative) }

func filePreview(path, provider, detail string, warning string) SetupPreview {
	return SetupPreview{Provider: provider, Changes: []SetupChange{{Path: path, Action: "create", Detail: detail}}, Warnings: []string{warning}}
}

func removeOwnedFile(path, provider string) (SetupResult, error) {
	preview := SetupPreview{Provider: provider, Changes: []SetupChange{{Path: path, Action: "remove", Detail: "Remove only the exact Craftlane-owned plugin"}}}
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return SetupResult{Preview: preview}, nil
	}
	return SetupResult{Preview: preview, Changed: err == nil}, err
}

func pluginPath(root, provider string) string {
	return filepath.Join(root, ".opencode", "plugins", fmt.Sprintf("craftlane-observability.%s.ts", provider))
}
