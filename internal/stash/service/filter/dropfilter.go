package filter

import (
	"BlackHole/pkg/config"
	"strings"
)

func DropFilter(conds []config.ConditionConf) FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		var qualify bool
		for _, cond := range conds {
			var qualifyOnce bool
			switch cond.Type {
			case typeMatch:
				qualifyOnce = cond.Value == m[cond.Key]
			case typeContains:
				if val, ok := m[cond.Key].(string); ok {
					qualifyOnce = strings.Contains(val, cond.Value)
				}
			}

			switch cond.Op {
			case opAnd:
				if !qualifyOnce {
					return m
				} else {
					qualify = true
				}
			case opOr:
				if qualifyOnce {
					qualify = true
				}
			}
		}

		if qualify {
			return nil
		} else {
			return m
		}
	}
}
