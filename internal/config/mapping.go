package config

import (
	_ "embed"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type Mapping struct {
	Version int               `yaml:"version"`
	Images  []ImagesChangelog `yaml:"images"`
}

type ImagesChangelog struct {
	Names []string      `yaml:"names"`
	Url   string        `yaml:"url"`
	Type  ChangelogType `yaml:"type"`
}

//go:generate go run github.com/dmarkham/enumer -type=ChangelogType -yaml -json
type ChangelogType int

const (
	Plaintext ChangelogType = iota
	Markdown
	Asciidoc
	Html
	Github
)

//go:embed mapping.yml
var mapping []byte

func loadInternalMappings() ([]ImagesChangelog, error) {
	var cfg Mapping
	err := yaml.Unmarshal(mapping, &cfg)
	if err != nil {
		return nil, err
	}

	log.Tracef("loaded default mappings: %d", len(cfg.Images))
	return cfg.Images, nil
}

func LoadMappings(name string) ([]ImagesChangelog, error) {
	internal, err := loadInternalMappings()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		return internal, nil
	}

	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var cfg Mapping
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	res := make([]ImagesChangelog, len(cfg.Images))
	copy(res, cfg.Images)

	for _, v := range internal {
		exists := false
		for _, name := range v.Names {
			cmp := func(i ImagesChangelog) bool {
				return slices.Contains(i.Names, name)
			}
			idx := slices.IndexFunc(res, cmp)
			if idx != -1 {
				exists = true
				break
			}
		}

		if !exists {
			res = append(res, v)
		}
	}

	return res, nil
}
