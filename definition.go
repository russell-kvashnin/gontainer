package gontainer

// ServiceDefinition - struct for service configuration.
// Name represents service name
// injection - not exported, and filled automatically by containerBuilder
type ServiceDefinition struct {
	Name       string
	Factory    Factory
	injections injections
}

// Factory for service instantiation
// Constructor - function, that return target object
type Factory struct {
	Constructor interface{}
	Args        ConstructorArguments
}

// ConstructorArguments - array of object constructor arguments
type ConstructorArguments []interface{}

// ServiceDefinitions - array of ServiceDefinition objects
type ServiceDefinitions []ServiceDefinition

// Has - provide "is definition exists in ServiceDefinitions"
func (defs ServiceDefinitions) Has(name string) bool {
	for i := range defs {
		if defs[i].Name == name {
			return true
		}
	}

	return false
}
