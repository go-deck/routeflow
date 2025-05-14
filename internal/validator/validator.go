package validator

import (
	"errors"
	"strconv"
)

// ValidateBody checks request data against validation rules
func ValidateBody(bodyData map[string]interface{}, validation map[string]interface{}) error {
	for field, rules := range validation {
		ruleMap, ok := rules.(map[string]interface{})
		if !ok {
			continue
		}

		value := bodyData[field]
		strVal, isString := value.(string)

		// Required field validation
		if req, exists := ruleMap["required"]; exists && req.(bool) && !exists {
			return errors.New(field + " is required")
		}

		// Check for specific validation rules
		for key, param := range ruleMap {
			switch key {
			case "min_length":
				if isString {
					min := int(param.(int))
					if len(strVal) < min {
						return errors.New(field + " must be at least " + strconv.Itoa(min) + " characters long")
					}
				}
			case "max_length":
				if isString {
					max := int(param.(int))
					if len(strVal) > max {
						return errors.New(field + " must be at most " + strconv.Itoa(max) + " characters long")
					}
				}
			}
		}
	}
	return nil
}
