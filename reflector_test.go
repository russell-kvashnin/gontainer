package gontainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ReflectorMock struct {
	mock.Mock
}

func (rm *ReflectorMock) runMethod(obj interface{}, methodName string, args []interface{}) ([]interface{}, error) {
	mArgs := rm.Called(obj, methodName, args[0])

	return mArgs.Get(0).([]interface{}), mArgs.Error(1)
}
func (rm *ReflectorMock) runFunction(fn interface{}, args []interface{}) ([]interface{}, error) {
	mArgs := rm.Called(fn, args)

	return mArgs.Get(0).([]interface{}), mArgs.Error(1)
}
func (rm *ReflectorMock) setFieldValue(obj interface{}, fieldName string, arg interface{}) error {
	args := rm.Called(obj, fieldName, arg)

	return args.Error(0)
}

type ReflectorSuite struct {
	suite.Suite
	r reflector
}

func (suite *ReflectorSuite) SetupTest() {
	r := new(reflectionUtils)

	suite.r = r
}

func (suite *ReflectorSuite) TestRunFunction() {
	fn := func(arg string) string {
		return arg
	}

	args := []interface{}{
		"arg",
	}

	res, err := suite.r.runFunction(fn, args)

	assert.Equal(suite.T(), "arg", res[0])
	assert.Nil(suite.T(), err)
}

func (suite *ReflectorSuite) TestRunFunctionNotCallable() {
	fn := "func"

	args := []interface{}{
		"arg",
	}

	res, err := suite.r.runFunction(fn, args)

	assert.Nil(suite.T(), res)
	assert.EqualError(suite.T(), err, "Given argument 'string' is not a function")
}

func (suite *ReflectorSuite) TestRunFunctionWrongArgCount() {
	fn := func(arg string) {}

	args := []interface{}{
		1, 2, 3,
	}

	res, err := suite.r.runFunction(fn, args)

	assert.Nil(suite.T(), res)
	assert.EqualError(suite.T(), err, "Invalid arguments count for 'func(string)' method, expected '1' but got '3'")
}

func (suite *ReflectorSuite) TestRunMethod() {
	svcStub := new(ServiceStub)
	arg := "test"

	res, err := suite.r.runMethod(svcStub, "Test", []interface{}{arg})

	assert.Equal(suite.T(), arg, res[0])
	assert.Nil(suite.T(), err)
}

func (suite *ReflectorSuite) TestRunMethodUndefined() {
	svcStub := new(ServiceStub)
	arg := "test"

	res, err := suite.r.runMethod(svcStub, "Undefined", []interface{}{arg})

	assert.Nil(suite.T(), res)
	assert.EqualError(suite.T(), err, "Method with name 'Undefined' not found")
}

func (suite *ReflectorSuite) TestRunMethodWrongArgCount() {
	svcStub := new(ServiceStub)
	args := []interface{}{
		1, 2, 3,
	}

	res, err := suite.r.runMethod(svcStub, "Test", args)

	assert.Nil(suite.T(), res)
	assert.EqualError(suite.T(), err, "Invalid arguments count for 'Test' method, expected '1' but got '3'")
}

func (suite *ReflectorSuite) TestSetFieldValue() {
	svcStub := new(ServiceStub)
	arg := "set_field_value_test"

	err := suite.r.setFieldValue(svcStub, "Public", arg)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), arg, svcStub.Public)
}

func (suite *ReflectorSuite) TestSetFieldValueFieldNotFound() {
	svcStub := new(ServiceStub)
	arg := "set_field_value_test"

	err := suite.r.setFieldValue(svcStub, "Undefined", arg)

	assert.EqualError(suite.T(), err, "Field with name 'Undefined' not found")
}

func TestReflectorSuite(t *testing.T) {
	suite.Run(t, new(ReflectorSuite))
}
