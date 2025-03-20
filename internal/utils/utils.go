package utils

import (
	"log"
	"reflect"

	"github.com/go-deck/routeflow/internal/ctx"
	"github.com/go-deck/routeflow/internal/loader"
)

// DiscoverHandlers extracts handler functions dynamically
func DiscoverHandlers(handlerStructs ...interface{}) map[string]ctx.HandlerFunc {
	handlerMap := make(map[string]ctx.HandlerFunc)

	for _, structInstance := range handlerStructs {
		v := reflect.ValueOf(structInstance)

		// Check if we have a pointer to struct
		if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
			t := v.Type()
			log.Printf("🔍 Inspecting handler: %s", t.String())

			for i := 0; i < v.NumMethod(); i++ {
				method := v.Method(i)
				methodName := t.Method(i).Name

				handlerFunc, ok := method.Interface().(func(*ctx.Context) (interface{}, int))
				if !ok {
					log.Printf("❌ Skipping method %s: Invalid signature", methodName)
					continue
				}

				handlerMap[methodName] = handlerFunc
				log.Printf("✅ Registered handler: %s", methodName)
			}
		} else {
			log.Printf("⚠️ Invalid handler type: %T (expected pointer to struct)", structInstance)
		}
	}

	return handlerMap
}

// Map YAML handlers to discovered handlers
func MapHandlersFromYAML(cfg *loader.Config, discovered map[string]ctx.HandlerFunc) map[string]func(*ctx.Context) (interface{}, int) {
	mappedHandlers := make(map[string]func(*ctx.Context) (interface{}, int))

	for _, group := range cfg.Routes.Groups {
		for _, route := range group.Routes {
			handlerName := route.Handler // Handler name from YAML

			if fn, exists := discovered[handlerName]; exists {
				// Explicit type conversion
				mappedHandlers[handlerName] = func(c *ctx.Context) (interface{}, int) {
					return fn(c) // Invoke the discovered function
				}
				log.Printf("✔ Mapped handler: %s", handlerName)
			} else {
				log.Fatalf("❌ Handler not found: %s", handlerName)
			}
		}
	}

	return mappedHandlers
}
