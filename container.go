package gontainer

import (
	"errors"
)

// Container - implementation of ServiceContainer interface
type Container struct {
	services     map[string]interface{}
	builder      containerBuilder
	injector     serviceInjector
	instantiator objectInstantiator
	compiled     bool
}

// ServiceContainer describes a DI container behavior
// Pass ServiceDefinitions to Compile method, after successful compilation
// services will be available via Get(name) method
type ServiceContainer interface {
	Compile(defs ServiceDefinitions) []error
	Get(name string) (interface{}, error)
	set(name string, svc interface{})
	getDefinition(name Injection) (interface{}, error)
}

// NewContainer - Container object constructor
func NewContainer() *Container {
	cnt := new(Container)
	cnt.services = make(map[string]interface{})
	cnt.builder = newContainerBuilder()

	r := newReflector()
	cnt.instantiator = newInstantiator(cnt, r)
	cnt.injector = newInjector(cnt, r)

	return cnt
}

// Compile method - entry point for Container configuration
// accept array of ServiceDefinition (ServiceDefinitions)
// return array of error's if something goes wrong
func (c *Container) Compile(defs ServiceDefinitions) []error {
	var errs []error

	if c.compiled {
		msg := "Container already compiled"

		return []error{errors.New(msg)}
	}
	errs = c.doCompile(defs)
	if len(errs) > 0 {
		return errs
	}

	c.compiled = true

	return nil
}

func (c *Container) doCompile(defs ServiceDefinitions) []error {
	var errs []error

	errs = c.builder.build(defs)
	if len(errs) > 0 {
		return errs
	}

	for i := range c.builder.getDefinitions() {
		def := defs[i]

		_, err := c.Get(def.Name)
		if err == nil {
			continue
		}

		inst, err := c.instantiator.instantiate(def)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		if len(def.injections) > 0 {
			err := c.injector.inject(inst, def.injections)
			if err != nil {
				errs = append(errs, err)

				continue
			}
		}

		c.services[def.Name] = inst
	}

	return errs
}

func (c *Container) has(name string) bool {
	_, ok := c.services[name]

	return ok
}

func (c *Container) set(name string, svc interface{}) {
	if !c.has(name) {
		c.services[name] = svc
	}
}

// Get provide access or configured services
// accept service name as argument and return service if it exists,
// or error if not
func (c *Container) Get(name string) (interface{}, error) {
	if c.has(name) {
		return c.services[name], nil
	}

	return nil, errors.New("Undefined service")
}

// Method shortcut
func (c *Container) getDefinition(name Injection) (interface{}, error) {
	return c.builder.getDefinition(name)
}
