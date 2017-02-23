package gontainer

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

var typ struct {
	InjectionType        reflect.Type
	ConstructorInjection reflect.Type
	SetterInjection      reflect.Type
	PropertyInjection    reflect.Type
}

func init() {
	typ.InjectionType = reflect.TypeOf(Injection(""))

	typ.ConstructorInjection = reflect.TypeOf(constructorInjection{})
	typ.SetterInjection = reflect.TypeOf(setterInjection{})
	typ.PropertyInjection = reflect.TypeOf(propertyInjection{})
}

type reflectionUtils struct {
}
type reflector interface {
	runMethod(obj interface{}, methodName string, args []interface{}) ([]interface{}, error)
	runFunction(fn interface{}, args []interface{}) ([]interface{}, error)
	setFieldValue(obj interface{}, fieldName string, arg interface{}) error
}

func newReflector() *reflectionUtils {
	inst := new(reflectionUtils)

	return inst
}

func (r *reflectionUtils) runMethod(obj interface{}, methodName string, args []interface{}) ([]interface{}, error) {
	objVal := reflect.ValueOf(obj)
	method := objVal.MethodByName(methodName)

	if method.Kind() != reflect.Func {
		errMsg := fmt.Sprintf("Method with name '%s' not found", methodName)

		return nil, errors.New(errMsg)
	}

	err := r.checkFnArgCount(method, methodName, args)
	if err != nil {
		return nil, err
	}

	argVals := r.arrToValsArr(args)
	retVals := method.Call(argVals)

	ret := make([]interface{}, len(retVals))

	for i, retVal := range retVals {
		ret[i] = retVal.Interface()
	}

	return ret, nil
}

func (r *reflectionUtils) runFunction(fn interface{}, args []interface{}) ([]interface{}, error) {
	var ret []interface{}
	fnValue := reflect.ValueOf(fn)
	fnBody := reflect.TypeOf(fn).String()

	if fnValue.Kind() != reflect.Func {
		errMsg := fmt.Sprintf("Given argument '%s' is not a function", fnBody)

		return nil, errors.New(errMsg)
	}

	err := r.checkFnArgCount(fnValue, fnBody, args)
	if err != nil {
		return nil, err
	}

	argVals := r.arrToValsArr(args)
	resVal := fnValue.Call(argVals)

	ret = r.valsArrToArr(resVal)

	return ret, nil
}

func (r *reflectionUtils) setFieldValue(obj interface{}, fieldName string, arg interface{}) error {
	var err error
	objVal := reflect.ValueOf(obj)
	argValue := reflect.ValueOf(arg)

	field := objVal.Elem().FieldByName(fieldName)
	if field.Kind() == reflect.Invalid {
		errMsg := fmt.Sprintf("Field with name '%s' not found", fieldName)

		return errors.New(errMsg)
	}

	field.Set(argValue)

	return err
}

func (r *reflectionUtils) checkFnArgCount(fn reflect.Value, name string, args []interface{}) error {
	argCount := fn.Type().NumIn()
	if len(args) != argCount {
		errMsg := fmt.Sprintf("Invalid arguments count for '%s' method, expected '%d' but got '%d'",
			name, argCount, len(args))
		return errors.New(errMsg)
	}

	return nil
}

func (r *reflectionUtils) arrToValsArr(arr []interface{}) []reflect.Value {
	elemVals := make([]reflect.Value, len(arr))
	for i, elem := range arr {
		elemVal := reflect.ValueOf(elem)
		elemVals[i] = elemVal
	}

	return elemVals
}

func (r *reflectionUtils) valsArrToArr(valsArr []reflect.Value) []interface{} {
	arr := make([]interface{}, len(valsArr))
	for i, elemVal := range valsArr {
		elem := elemVal.Interface()
		arr[i] = elem
	}

	return arr
}
