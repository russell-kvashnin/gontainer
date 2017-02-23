package gontainer

import (
	"errors"
	"fmt"
	"reflect"
)

type constructorInjection struct {
	SvcName string
}
type setterInjection struct {
	SvcName    string
	MethodName string
}
type propertyInjection struct {
	SvcName      string
	PropertyName string
}

// Injection used in ServiceDefinition
// to configure constructor injection
type Injection string
type injections []interface{}

type injector struct {
	cnt ServiceContainer
	r   reflector
}
type serviceInjector interface {
	inject(svc interface{}, injections injections) error
}

func newInjector(cnt ServiceContainer, r reflector) *injector {
	inst := new(injector)
	inst.cnt = cnt
	inst.r = r

	return inst
}

func (inj *injector) inject(svc interface{}, injections injections) error {
	for _, injection := range injections {
		switch reflect.TypeOf(injection) {
		case typ.SetterInjection:
			err := inj.setterInjection(svc, injection.(setterInjection))
			if err != nil {
				return err
			}
			break
		case typ.PropertyInjection:
			err := inj.propertyInjection(svc, injection.(propertyInjection))
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func (inj *injector) setterInjection(svc interface{}, injection setterInjection) error {
	injSvc, err := inj.cnt.Get(injection.SvcName)
	if err != nil {
		msg := fmt.Sprintf("Cannot inject %s , service not found.",
			injection.SvcName)

		return errors.New(msg)
	}

	_, err = inj.r.runMethod(svc, injection.MethodName, []interface{}{injSvc})

	return err
}

func (inj *injector) propertyInjection(svc interface{}, injection propertyInjection) error {
	injSvc, err := inj.cnt.Get(injection.SvcName)
	if err != nil {
		msg := fmt.Sprintf("Cannot inject %s , service not found.",
			injection.SvcName)

		return errors.New(msg)
	}

	err = inj.r.setFieldValue(svc, injection.PropertyName, injSvc)

	return err
}
