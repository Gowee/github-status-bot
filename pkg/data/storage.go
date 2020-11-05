package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gowee/github-status-bot/pkg/utils"
)

// A Storage backed by JSON, simplistic
type Storage struct {
	FilePath string
}

func (s *Storage) EnsureInitialized(empty interface{}) error {
	if utils.IsPathNotExisting(s.FilePath) {
		if err := s.Store(empty); err != nil {
			return fmt.Errorf("Failed to initialize data file: %w", err)
		}
	}
	return nil
}

func (s *Storage) Load(data interface{}) error {
	raw, err := ioutil.ReadFile(s.FilePath)
	if err = json.Unmarshal(raw, data); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Store(data interface{}) error {
	serd, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.FilePath, serd, 0644)
}

// WTF: why lint prevent self and this?
// 	ref: https://stackoverflow.com/questions/23482068/in-go-is-naming-the-receiver-variable-self-misleading-or-good-practice

// WTF: why no generic struct so that to be typesafe? e.g. HashMap, data.Database/Storage
