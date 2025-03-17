package validator

import (
	"errors"
	"regexp"
	"strconv"
)

// ValidateBody checks request data against validation rules
func ValidateBody(bodyData map[string]interface{}, props map[string]interface{}) error {
	for field, rules := range props {
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
			case "format":
				if isString {
					switch param.(string) {
					case "email":
						if !isValidEmail(strVal) {
							return errors.New(field + " must be a valid email")
						}
					case "phone":
						if !isValidPhone(strVal) {
							return errors.New(field + " must be a valid phone number")
						}
					}
				}
			}
		}
	}
	return nil
}

// isValidEmail checks if the given string is a valid email
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// isValidPhone checks if the given string is a valid phone number
func isValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	return re.MatchString(phone)
}
