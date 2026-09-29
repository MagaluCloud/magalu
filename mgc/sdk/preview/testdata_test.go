package preview

import (
	"errors"
	"io/fs"
	"sort"
)

const fooSpec = `openapi: 3.0.3
info:
  title: Foo
  description: Foo preview API.
  version: '1.0'
servers:
- url: https://foo.example.com
tags:
- name: things
  description: Things of foo.
paths:
  /v1/things:
    get:
      tags: [things]
      operationId: list-things
      summary: List things
      description: List all things.
      responses:
        '200':
          description: OK
`

const brokenSpec = `openapi: 3.0.3
info: [this is not
`

var extensionPrefix = "x-mgc"

// memStore is a Store kept in memory.
type memStore struct {
	specs map[string]string
	err   error
}

func newMemStore(specs map[string]string) *memStore {
	return &memStore{specs: specs}
}

func (s *memStore) Names() ([]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	names := make([]string, 0, len(s.specs))
	for name := range s.specs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func (s *memStore) Read(name string) (Package, error) {
	spec, ok := s.specs[name]
	if !ok {
		return Package{}, fs.ErrNotExist
	}
	return Package{Name: name, Spec: []byte(spec)}, nil
}

func (s *memStore) String() string {
	return "mem-store"
}

// fakeValidator refuses the packages listed in refused and records what it saw.
type fakeValidator struct {
	refused map[string]string
	seen    []string
}

func (v *fakeValidator) Validate(level Level, pkg Package) error {
	v.seen = append(v.seen, level.Name+"/"+pkg.Name)
	if msg, ok := v.refused[pkg.Name]; ok {
		return errors.New(msg)
	}
	return nil
}
