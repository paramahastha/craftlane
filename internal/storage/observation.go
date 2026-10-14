package storage

import (
	"encoding/json"
	"errors"
	"os"
)

const observationSecretFile = "observation.token"
const observationFingerprintsFile = "observation-fingerprints.json"

func (s Store) SaveObservationSecret(secret string) error {
	if secret == "" {
		return errors.New("empty observation secret")
	}
	if err := os.MkdirAll(s.Dir, 0700); err != nil {
		return err
	}
	return os.WriteFile(s.path(observationSecretFile), []byte(secret), 0600)
}

func (s Store) LoadObservationSecret() (string, error) {
	b, err := os.ReadFile(s.path(observationSecretFile))
	return string(b), err
}

func (s Store) ClearObservationSecret() error {
	err := os.Remove(s.path(observationSecretFile))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s Store) LoadObservationFingerprints() ([]string, error) {
	b, err := os.ReadFile(s.path(observationFingerprintsFile))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var fingerprints []string
	if err := json.Unmarshal(b, &fingerprints); err != nil {
		return nil, err
	}
	return fingerprints, nil
}

func (s Store) SaveObservationFingerprints(fingerprints []string) error {
	return s.write(observationFingerprintsFile, fingerprints)
}

func (s Store) AddObservationFingerprint(fingerprint string, limit int) (bool, error) {
	if fingerprint == "" {
		return false, errors.New("empty observation fingerprint")
	}
	items, err := s.LoadObservationFingerprints()
	if err != nil {
		return false, err
	}
	for _, item := range items {
		if item == fingerprint {
			return false, nil
		}
	}
	items = append(items, fingerprint)
	if limit < 1 {
		limit = 1000
	}
	if len(items) > limit {
		items = items[len(items)-limit:]
	}
	return true, s.SaveObservationFingerprints(items)
}
