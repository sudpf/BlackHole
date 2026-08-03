package apperror

import (
	"fmt"
	"sort"
)

type Definition struct {
	Code       Code
	HTTPStatus int
}

type Catalog struct {
	definitions map[Code]Definition
}

func NewCatalog(definitions ...Definition) (*Catalog, error) {
	catalog := &Catalog{
		definitions: make(map[Code]Definition, len(definitions)),
	}

	for _, definition := range definitions {
		if definition.Code < 0 {
			return nil, fmt.Errorf("invalid error code %d", definition.Code)
		}
		if definition.HTTPStatus < 100 || definition.HTTPStatus > 599 {
			return nil, fmt.Errorf("invalid HTTP status %d for error code %d", definition.HTTPStatus, definition.Code)
		}
		if definition.Code == Success && (definition.HTTPStatus < 200 || definition.HTTPStatus > 299) {
			return nil, fmt.Errorf("success code %d requires a 2xx HTTP status", Success)
		}
		if _, exists := catalog.definitions[definition.Code]; exists {
			return nil, fmt.Errorf("duplicate error code %d", definition.Code)
		}
		catalog.definitions[definition.Code] = definition
	}

	if _, exists := catalog.definitions[Success]; !exists {
		return nil, fmt.Errorf("success code %d is not defined", Success)
	}

	return catalog, nil
}

func (c *Catalog) Lookup(code Code) (Definition, bool) {
	if c == nil {
		return Definition{}, false
	}
	definition, exists := c.definitions[code]
	return definition, exists
}

func (c *Catalog) MessageIDs() []string {
	if c == nil {
		return nil
	}

	codes := make([]int, 0, len(c.definitions))
	for code := range c.definitions {
		codes = append(codes, int(code))
	}
	sort.Ints(codes)

	messageIDs := make([]string, 0, len(codes))
	for _, code := range codes {
		messageIDs = append(messageIDs, MessageID(Code(code)))
	}
	return messageIDs
}

func MessageID(code Code) string {
	return fmt.Sprintf("error_%d", code)
}
