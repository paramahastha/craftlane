package codex

import "github.com/craftlane/craftlane/internal/providers"

func Preview(root string) (providers.SetupPreview, error) { return New(root).Preview() }
func Apply(root string) (providers.SetupResult, error)    { return New(root).Apply() }
func Remove(root string) (providers.SetupResult, error)   { return New(root).Remove() }
