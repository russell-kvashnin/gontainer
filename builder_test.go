package gontainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ContainerBuilderMock struct {
	mock.Mock
}

func (cbm *ContainerBuilderMock) build(defs ServiceDefinitions) []error {
	args := cbm.Called(defs)

	return args.Get(0).([]error)
}
func (cbm *ContainerBuilderMock) getDefinition(name Injection) (interface{}, error) {
	args := cbm.Called(name)

	return args.Get(0), args.Error(1)
}
func (cbm *ContainerBuilderMock) getDefinitions() ServiceDefinitions {
	args := cbm.Called()

	return args.Get(0).(ServiceDefinitions)
}

type TaggedServiceStub struct {
	Setter   string `inject:"nothing" inject_type:"setter" inject_method:"SetNothing"`
	Property string `inject:"nothing" inject_type:"property"`
}

type WrongTaggedServiceStub struct {
	Wrong string `inject:"nothing" inject_type:"wrong"`
}

type MethodNotProvidedServiceStub struct {
	Wrong string `inject:"nothing" inject_type:"setter"`
}

type ContainerBuilderTestSuite struct {
	suite.Suite
	builder         *definitionBuilder
	defs            ServiceDefinitions
	wrongDefs       ServiceDefinitions
	wrongSetterDefs ServiceDefinitions
	repeatableDefs  ServiceDefinitions
}

func (suite *ContainerBuilderTestSuite) SetupSuite() {
	svcStub := ServiceDefinition{
		Name: "test.service_stub",
		Factory: Factory{
			Constructor: func() *TaggedServiceStub { return new(TaggedServiceStub) },
		},
	}
	suite.defs = ServiceDefinitions{
		svcStub,
	}

	wrongStub := ServiceDefinition{
		Name: "test.service_stub",
		Factory: Factory{
			Constructor: func() *WrongTaggedServiceStub { return new(WrongTaggedServiceStub) },
		},
	}
	suite.wrongDefs = ServiceDefinitions{
		wrongStub,
	}

	wrongSetterStub := ServiceDefinition{
		Name: "test.service_stub",
		Factory: Factory{
			Constructor: func() *MethodNotProvidedServiceStub { return new(MethodNotProvidedServiceStub) },
		},
	}
	suite.wrongSetterDefs = ServiceDefinitions{
		wrongSetterStub,
	}

	suite.repeatableDefs = ServiceDefinitions{
		svcStub,
		svcStub,
	}
}

func (suite *ContainerBuilderTestSuite) SetupTest() {
	inst := new(definitionBuilder)

	suite.builder = inst
}

func (suite *ContainerBuilderTestSuite) TestConstructor() {
	inst := newContainerBuilder()

	assert.Implements(suite.T(), (*containerBuilder)(nil), inst)
}

func (suite *ContainerBuilderTestSuite) TestBuild() {
	errs := suite.builder.build(suite.defs)
	defs := suite.builder.getDefinitions()

	assert.IsType(suite.T(), ServiceDefinitions{}, defs)
	assert.Empty(suite.T(), errs)
}

func (suite *ContainerBuilderTestSuite) TestBuildFailedOnInjections() {
	errs := suite.builder.build(suite.wrongDefs)
	defs := suite.builder.getDefinitions()

	assert.IsType(suite.T(), ServiceDefinitions{}, defs)
	if assert.NotEmpty(suite.T(), errs) {
		assert.EqualError(suite.T(), errs[0], "Unknown injection type 'wrong'")
	}
}

func (suite *ContainerBuilderTestSuite) TestBuildFailedOnSetterInjectionMethodNotExists() {
	errs := suite.builder.build(suite.wrongSetterDefs)
	defs := suite.builder.getDefinitions()

	assert.IsType(suite.T(), ServiceDefinitions{}, defs)
	if assert.NotEmpty(suite.T(), errs) {
		assert.EqualError(suite.T(), errs[0], "Must provide setter metod name for setter injection")
	}
}

func (suite *ContainerBuilderTestSuite) TestBuildAlreadyDefined() {
	errs := suite.builder.build(suite.repeatableDefs)

	assert.EqualError(suite.T(), errs[0], "Service definition with name 'test.service_stub' already defined")
}

func (suite *ContainerBuilderTestSuite) TestGetDefinition() {
	suite.builder.build(suite.defs)

	def, err := suite.builder.getDefinition("test.service_stub")

	assert.IsType(suite.T(), ServiceDefinition{}, def)
	assert.Nil(suite.T(), err)
}

func (suite *ContainerBuilderTestSuite) TestGetDefinitionFailedNotBuilt() {
	def, err := suite.builder.getDefinition("test.service_stub")

	assert.Nil(suite.T(), def)
	assert.EqualError(suite.T(), err, "Not built yet")
}

func (suite *ContainerBuilderTestSuite) TestGetDefinitionNotExists() {
	suite.builder.build(suite.defs)

	def, err := suite.builder.getDefinition("test.not_exists")

	assert.Nil(suite.T(), def)
	assert.EqualError(suite.T(), err, "Definition not exists")
}

func TestContainerBuilderSuite(t *testing.T) {
	suite.Run(t, new(ContainerBuilderTestSuite))
}
