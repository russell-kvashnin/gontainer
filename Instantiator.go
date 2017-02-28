package gontainer

import (
	"errors"
	"fmt"
	"reflect"
)

type instantiator struct {
	cnt   ServiceContainer
	r     reflector
	chain []string
}
type objectInstantiator interface {
	instantiate(def ServiceDefinition) (interface{}, error)
}

func newInstantiator(cnt ServiceContainer, r reflector) *instantiator {
	inst := new(instantiator)
	inst.cnt = cnt
	inst.r = r

	return inst
}

func (ins *instantiator) instantiate(def ServiceDefinition) (interface{}, error) {
	var inst interface{}
	ins.chain = append(ins.chain, def.Name)

	args, err := ins.prepareConstructorArgs(def.Factory.Args, def.Factory.Constructor)
	if err != nil {
		ins.chain = []string{}

		return nil, err
	}

	res, err := ins.r.runFunction(def.Factory.Constructor, args)
	if err != nil {
		return nil, err
	}

	inst = res[0]
	ins.chain = []string{}

	return inst, nil
}

func (ins *instantiator) prepareConstructorArgs(args ConstructorArguments, fn interface{}) ([]interface{}, error) {
	cArgsCount := reflect.TypeOf(fn).NumIn()
	if len(args) != cArgsCount {
		msg := "Arguments count mismatch"

		return nil, errors.New(msg)
	}

	for i, v := range args {
		fnArgType := reflect.TypeOf(fn).In(i)
		vArgType := reflect.TypeOf(v)
		argTypeAll := reflect.TypeOf([]interface {}{})

		// Declared function argument and given has different types
		// Except case when arg type is []interface{} (any type)
		if vArgType != fnArgType && vArgType != typ.InjectionType && fnArgType != argTypeAll {
			err := fmt.Sprintf(
				"Argument type mismatch, expected %s, got %s",
				fnArgType,
				reflect.TypeOf(v))

			return nil, errors.New(err)
		}

		// If argument type is constructor injection
		if vArgType == typ.InjectionType {
			svcVal, err := ins.prepareConstructorInjArgs(v.(Injection))
			if err != nil {
				return nil, err
			}
			args[i] = svcVal

			continue
		}
	}

	return args, nil
}

func (ins *instantiator) prepareConstructorInjArgs(arg Injection) (interface{}, error) {
	var ret interface{}

	svc, err := ins.cnt.Get(string(arg))
	if err != nil {
		for i := range ins.chain {
			if ins.chain[i] == string(arg) {
				last := ins.chain[len(ins.chain)-1]
				msg := fmt.Sprintf(
					"Cannot compile container, circular reference found, %s <-> %s",
					ins.chain[i], last)

				return ret, errors.New(msg)
			}
		}

		def, err := ins.cnt.getDefinition(arg)
		if err != nil {
			return ret, err
		}

		inj, err := ins.instantiate(def.(ServiceDefinition))
		if err != nil {
			return ret, err
		}

		ins.cnt.set(string(arg), inj)

		return inj, nil
	}

	return svc, nil
}
