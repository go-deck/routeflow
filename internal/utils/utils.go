package utils

import (
	"log"
	"reflect"

	"github.com/go-deck/routeflow/internal/ctx"
	"github.com/go-deck/routeflow/internal/loader"
)

// Enhanced discovery with clear logging
func DiscoverHandlers(handlerStructs ...interface{}) map[string]ctx.HandlerFunc {
	handlers := make(map[string]ctx.HandlerFunc)

	for _, structInstance := range handlerStructs {
		v := reflect.ValueOf(structInstance)
		if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
			log.Printf("⚠️ Invalid handler type: %T", structInstance)
			continue
		}

		t := v.Type()
		for i := 0; i < v.NumMethod(); i++ {
			method := v.Method(i)
			name := t.Method(i).Name

			if fn, ok := method.Interface().(func(*ctx.Context) (interface{}, int)); ok {
				handlers[name] = fn
				log.Printf("✅ Registered handler: %s", name)
			}
		}
	}
	return handlers
}

// Map YAML handlers to discovered handlers
func MapHandlersFromYAML(cfg *loader.Config, discovered map[string]ctx.HandlerFunc) map[string]func(*ctx.Context) (interface{}, int) {
	mappedHandlers := make(map[string]func(*ctx.Context) (interface{}, int))

	for _, group := range cfg.Routes.Groups {
		for _, route := range group.Routes {
			handlerName := route.Handler

			if fn, exists := discovered[handlerName]; exists {
				// Direct assignment avoids closure issues
				mappedHandlers[handlerName] = fn
				log.Printf("✔ Mapped handler: %s", handlerName)
			} else {
				log.Fatalf("❌ Handler not found: %s", handlerName)
			}
		}
	}

	return mappedHandlers
}

func MapMiddlewareHandlersFromYAML(cfg *loader.Config, discovered map[string]ctx.HandlerFunc) map[string]func(*ctx.Context) (interface{}, int) {
	mappedHandlers := make(map[string]func(*ctx.Context) (interface{}, int))

	// Collect all middleware references from all levels
	var allMiddlewareNames []string

	// Global middlewares
	allMiddlewareNames = append(allMiddlewareNames, cfg.Middlewares.BuiltIn...)
	allMiddlewareNames = append(allMiddlewareNames, cfg.Middlewares.Custom...)

	// Traverse route groups and routes
	for _, group := range cfg.Routes.Groups {
		// Group-level middlewares
		allMiddlewareNames = append(allMiddlewareNames, group.Middlewares.BuiltIn...)
		allMiddlewareNames = append(allMiddlewareNames, group.Middlewares.Custom...)

		// Route-level middlewares
		for _, route := range group.Routes {
			allMiddlewareNames = append(allMiddlewareNames, route.Middlewares.BuiltIn...)
			allMiddlewareNames = append(allMiddlewareNames, route.Middlewares.Custom...)
		}
	}

	// Process all collected middleware names
	for _, name := range unique(allMiddlewareNames) {
		if fn, exists := discovered[name]; exists {
			mappedHandlers[name] = fn
			log.Printf("✔ Mapped middleware: %s", name)
		} else {
			log.Printf("⚠️ Middleware not found: %s", name)
		}
	}

	return mappedHandlers
}

// Helper function to get unique values
func unique(input []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(input))
	for _, val := range input {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}
