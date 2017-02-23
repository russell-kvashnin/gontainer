package gontainer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceStub struct {
	Public    string
	injection *InjectionStub
}
type InjectionStub struct {
}

func (stub *ServiceStub) SetInjection(injection InjectionStub) {
	stub.injection = &injection
}

func (stub *ServiceStub) Test(arg string) string {
	return arg
}

type ContainerMock struct {
	mock.Mock
}

func (cm *ContainerMock) Get(name string) (interface{}, error) {
	args := cm.Called(name)

	return args.Get(0), args.Error(1)
}
func (cm *ContainerMock) Compile(defs ServiceDefinitions) []error {
	args := cm.Called()

	return args.Get(0).([]error)
}
func (cm *ContainerMock) getDefinition(name Injection) (interface{}, error) {
	args := cm.Called(name)

	return args.Get(0), args.Error(1)
}
func (cm *ContainerMock) set(name string, svc interface{}) {
}

type ContainerTestSuite struct {
	suite.Suite
	cbm            *ContainerBuilderMock
	inst           *InstantiatorMock
	inj            *InjectorMock
	cnt            *Container
	svc            *ServiceStub
	defs           ServiceDefinitions
	defsInjFailed  ServiceDefinitions
	defsInstFailed ServiceDefinitions
}

func (suite *ContainerTestSuite) SetupSuite() {
	def := ServiceDefinition{
		Name: "test.service_stub",
	}
	suite.defs = ServiceDefinitions{
		def,
	}

	defInstFailed := ServiceDefinition{
		Name: "test.service_stub_inst_failed",
	}
	suite.defsInstFailed = ServiceDefinitions{
		defInstFailed,
	}

	defInjFailed := ServiceDefinition{
		Name: "test.service_stub_inj_failed",
		injections: injections{
			Injection("test.service_stub_inj_failed"),
		},
	}
	suite.defsInjFailed = ServiceDefinitions{
		defInjFailed,
	}

	suite.cbm = new(ContainerBuilderMock)
	suite.inst = new(InstantiatorMock)
	suite.inj = new(InjectorMock)

	suite.cbm.On("build", suite.defs).
		Return([]error{})

	suite.cbm.On("build", suite.defsInstFailed).
		Return([]error{})

	suite.cbm.On("build", suite.defsInjFailed).
		Return([]error{})

	suite.cbm.On("getDefinitions").
		Return(suite.defs)

	suite.inst.On("instantiate", suite.defs[0]).
		Return(new(ServiceStub), nil)

	suite.inst.On("instantiate", suite.defsInstFailed[0]).
		Return(nil, errors.New("Arguments count mismatch"))

	suite.inst.On("instantiate", suite.defsInjFailed[0]).
		Return(suite.svc, nil)

	suite.inj.On("inject", suite.svc, suite.defsInjFailed[0].injections).
		Return(errors.New("Cannot inject test.service_stub_inj_failed, service not found"))
}

func (suite *ContainerTestSuite) SetupTest() {
	suite.svc = new(ServiceStub)

	cnt := new(Container)
	cnt.builder = suite.cbm
	cnt.injector = suite.inj
	cnt.instantiator = suite.inst
	cnt.services = make(map[string]interface{})

	suite.cnt = cnt
}

func (suite *ContainerTestSuite) TestConstructor() {
	cnt := NewContainer()

	assert.Implements(suite.T(), (*ServiceContainer)(nil), cnt)
}

func (suite *ContainerTestSuite) TestCompile() {
	suite.cnt.Compile(suite.defs)

	assert.True(suite.T(), suite.cnt.compiled)
}

func (suite *ContainerTestSuite) TestCompilationFailed() {
	suite.cnt.compiled = true
	suite.cnt.Compile(suite.defs)

	assert.True(suite.T(), suite.cnt.compiled)
}

func (suite *ContainerTestSuite) TestCompilationFailedOnInstantiation() {
	errs := suite.cnt.Compile(suite.defsInstFailed)

	assert.EqualError(suite.T(), errs[0], "Arguments count mismatch")
}

func (suite *ContainerTestSuite) TestCompilationFailedOnInjection() {
	errs := suite.cnt.Compile(suite.defsInjFailed)

	assert.EqualError(suite.T(), errs[0], "Cannot inject test.service_stub_inj_failed, service not found")
}

func (suite *ContainerTestSuite) TestGet() {
	suite.cnt.Compile(suite.defs)

	svc, err := suite.cnt.Get("test.service_stub")

	assert.IsType(suite.T(), &ServiceStub{}, svc)
	assert.Nil(suite.T(), err)
}

func (suite *ContainerTestSuite) TestGetUndefined() {
	suite.cnt.Compile(suite.defs)

	svc, err := suite.cnt.Get("test.service_undefined")

	assert.IsType(suite.T(), nil, svc)
	assert.Error(suite.T(), err)
}

func (suite *ContainerTestSuite) TestSet() {
	suite.cnt.set("test.service_stub", suite.svc)

	svc := suite.cnt.services["test.service_stub"]

	assert.Equal(suite.T(), suite.svc, svc)
}

func TestContainerTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}
