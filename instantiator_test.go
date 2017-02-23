package gontainer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type InstantiatorMock struct {
	mock.Mock
}

func (instMock *InstantiatorMock) instantiate(def ServiceDefinition) (interface{}, error) {
	args := instMock.Called(def)

	return args.Get(0), args.Error(1)
}

type InstantiatorSuite struct {
	suite.Suite
	cnt                   *ContainerMock
	refl                  *ReflectorMock
	inst                  objectInstantiator
	def                   ServiceDefinition
	defInj                ServiceDefinition
	defWrongArgCount      ServiceDefinition
	defWrongArgType       ServiceDefinition
	defUndefinedInjection ServiceDefinition
	defCircularReference  ServiceDefinition
}

func (suite *InstantiatorSuite) SetupSuite() {
	suite.def = ServiceDefinition{
		Name: "test.service_stub",
		Factory: Factory{
			Constructor: func(name string, injection constructorInjection) *ServiceStub { return new(ServiceStub) },
			Args: ConstructorArguments{
				"fake_name",
				Injection("test.injection"),
			},
		},
	}

	suite.defInj = ServiceDefinition{
		Name: "test.injection",
		Factory: Factory{
			Constructor: func() *InjectionStub { return new(InjectionStub) },
		},
	}

	suite.defWrongArgCount = ServiceDefinition{
		Name: "test.service_stub",
		Factory: Factory{
			Constructor: func(name string, injection constructorInjection) *ServiceStub { return new(ServiceStub) },
			Args: ConstructorArguments{
				"fake_name",
			},
		},
	}

	suite.defWrongArgType = ServiceDefinition{
		Name: "test.service_stub",
		Factory: Factory{
			Constructor: func(name string, injection constructorInjection) *ServiceStub { return new(ServiceStub) },
			Args: ConstructorArguments{
				0,
				Injection("test.wrong"),
			},
		},
	}

	suite.defUndefinedInjection = ServiceDefinition{
		Name: "test.service_stub",
		Factory: Factory{
			Constructor: func(name string, injection constructorInjection) *ServiceStub { return new(ServiceStub) },
			Args: ConstructorArguments{
				"string",
				Injection("test.wrong"),
			},
		},
	}

	suite.defCircularReference = ServiceDefinition{
		Name: "test.service_stub",
		Factory: Factory{
			Constructor: func(name string, injection constructorInjection) *ServiceStub { return new(ServiceStub) },
			Args: ConstructorArguments{
				"string",
				Injection("test.service_stub"),
			},
		},
	}

	suite.refl = new(ReflectorMock)
	suite.refl.On("runFunction", mock.Anything, mock.Anything).
		Return([]interface{}{new(ServiceStub)}, nil)

	suite.cnt = new(ContainerMock)
	suite.cnt.On("Get", "test.service_stub").
		Return(nil, errors.New("Undefined service"))
	suite.cnt.On("Get", "test.injection").
		Return(new(InjectionStub), nil)

	suite.cnt.On("Get", "test.wrong").
		Return(nil, errors.New("Undefined service"))

	suite.cnt.On("getDefinition", Injection("test.wrong")).
		Return(suite.defInj, nil)
}

func (suite *InstantiatorSuite) SetupTest() {

	inst := new(instantiator)
	inst.r = suite.refl
	inst.cnt = suite.cnt

	suite.inst = inst
}

func (suite *InstantiatorSuite) TestConstructor() {
	cntMock := new(ContainerMock)
	reflMock := new(ReflectorMock)

	inst := newInstantiator(cntMock, reflMock)

	assert.Implements(suite.T(), (*objectInstantiator)(nil), inst)
}

func (suite *InstantiatorSuite) TestInstantiate() {
	svc, err := suite.inst.instantiate(suite.def)

	assert.IsType(suite.T(), &ServiceStub{}, svc)
	assert.Nil(suite.T(), err)
}

func (suite *InstantiatorSuite) TestInstantiateWrongArgCount() {
	svc, err := suite.inst.instantiate(suite.defWrongArgCount)

	assert.EqualError(suite.T(), err, "Arguments count mismatch")
	assert.Nil(suite.T(), svc)
}

func (suite *InstantiatorSuite) TestInstantiateWrongArgType() {
	svc, err := suite.inst.instantiate(suite.defWrongArgType)

	assert.EqualError(suite.T(), err, "Argument type mismatch, expected string, got int")
	assert.Nil(suite.T(), svc)
}

func (suite *InstantiatorSuite) TestInstantiateUndefinedInjection() {
	svc, err := suite.inst.instantiate(suite.defUndefinedInjection)

	assert.IsType(suite.T(), &ServiceStub{}, svc)
	assert.IsType(suite.T(), &InjectionStub{}, svc.(*ServiceStub).injection)
	assert.Nil(suite.T(), err)
}

func (suite *InstantiatorSuite) TestInstantiateCircularReference() {
	svc, err := suite.inst.instantiate(suite.defCircularReference)

	errStr := err.Error()

	assert.Nil(suite.T(), svc)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), errStr, "Cannot compile container, circular reference found")
}

func TestInstantiatorSuite(t *testing.T) {
	suite.Run(t, new(InstantiatorSuite))
}
