package gontainer

import (
	"errors"
	"fmt"
	"reflect"
)

type definitionBuilder struct {
	defs  ServiceDefinitions
	built bool
}
type containerBuilder interface {
	build(defs ServiceDefinitions) []error
	getDefinition(name Injection) (interface{}, error)
	getDefinitions() ServiceDefinitions
}

func newContainerBuilder() containerBuilder {
	inst := new(definitionBuilder)
	inst.defs = ServiceDefinitions{}

	return inst
}

func (db *definitionBuilder) build(defs ServiceDefinitions) []error {
	var errs []error

	for _, def := range defs {
		if db.defs.Has(def.Name) {
			msg := fmt.Sprintf("Service definition with name '%s' already defined", def.Name)
			errs = append(errs, errors.New(msg))

			continue
		}

		def, err := db.buildDefinition(def)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		db.defs = append(db.defs, def)
	}

	if len(errs) == 0 {
		db.built = true
	}

	return errs
}

func (db *definitionBuilder) buildDefinition(def ServiceDefinition) (ServiceDefinition, error) {
	sType := reflect.TypeOf(def.Factory.Constructor).Out(0).Elem()
	fCnt := sType.NumField()

	fields := make([]reflect.StructField, fCnt)
	for i := 0; i < fCnt; i++ {
		fields[i] = sType.Field(i)
	}

	injections, err := db.buildInjections(fields)
	if err != nil {
		return def, err
	}

	def.injections = injections

	return def, nil
}

func (db *definitionBuilder) buildInjections(fields []reflect.StructField) ([]interface{}, error) {
	injections := []interface{}{}

	for _, field := range fields {
		tag := field.Tag.Get("inject")

		if len(tag) > 0 {
			injType := field.Tag.Get("inject_type")

			switch injType {
			case "setter":
				mName := field.Tag.Get("inject_method")
				if len(mName) == 0 {
					return nil, errors.New("Must provide setter metod name for setter injection")
				}

				inj := setterInjection{
					SvcName:    tag,
					MethodName: mName,
				}
				injections = append(injections, inj)
			case "property":
				inj := propertyInjection{
					SvcName:      tag,
					PropertyName: field.Name,
				}
				injections = append(injections, inj)
			default:
				return nil, fmt.Errorf("Unknown injection type '%s'", injType)
			}
		}
	}

	return injections, nil
}

func (db *definitionBuilder) getDefinition(name Injection) (interface{}, error) {
	if !db.built {
		return nil, errors.New("Not built yet")
	}

	for i := range db.defs {
		if db.defs[i].Name == string(name) {
			return db.defs[i], nil
		}
	}

	return nil, errors.New("Definition not exists")
}

func (db *definitionBuilder) getDefinitions() ServiceDefinitions {
	if !db.built {
		return ServiceDefinitions{}
	}

	return db.defs
}
