package gontainer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type InjectorMock struct {
	mock.Mock
}

func (injMock *InjectorMock) inject(svc interface{}, injections injections) error {
	args := injMock.Called(svc, injections)

	return args.Error(0)
}

type InjectionTestSuite struct {
	suite.Suite
	svcStub               *ServiceStub
	injStub               *InjectionStub
	cnt                   *ContainerMock
	refl                  *ReflectorMock
	injector              serviceInjector
	setterInjections      injections
	propertyInjections    injections
	wrongSetterInjections injections
	wrongPropInjections   injections
}

func (suite *InjectionTestSuite) SetupSuite() {
	suite.svcStub = new(ServiceStub)
	suite.injStub = new(InjectionStub)

	suite.setterInjections = injections{
		setterInjection{
			MethodName: "SetInjection",
			SvcName:    "test.service_stub",
		},
	}
	suite.propertyInjections = injections{
		propertyInjection{
			PropertyName: "Public",
			SvcName:      "test.service_stub",
		},
	}
	suite.wrongSetterInjections = injections{
		setterInjection{
			MethodName: "SetInjection",
			SvcName:    "test.wrong",
		},
	}
	suite.wrongPropInjections = injections{
		propertyInjection{
			SvcName:      "test.wrong",
			PropertyName: "WRONG",
		},
	}

	suite.cnt = new(ContainerMock)
	suite.cnt.On("Get", "test.service_stub").Return(suite.injStub, nil)
	suite.cnt.On("Get", "test.wrong").Return(nil, errors.New("Undefined service"))

	suite.refl = new(ReflectorMock)
	suite.refl.On("runMethod", suite.svcStub, "SetInjection", suite.injStub).Return([]interface{}{}, nil)
	suite.refl.On("setFieldValue", suite.svcStub, "Public", suite.injStub).Return(nil)
}

func (suite *InjectionTestSuite) SetupTest() {
	inj := new(injector)

	inj.cnt = suite.cnt
	inj.r = suite.refl

	suite.injector = inj
}

func (suite *InjectionTestSuite) TestSetterInjection() {
	s := new(ServiceStub)
	err := suite.injector.inject(s, suite.setterInjections)

	assert.Nil(suite.T(), err)
}

func (suite *InjectionTestSuite) TestPropertyInjection() {
	s := new(ServiceStub)
	err := suite.injector.inject(s, suite.propertyInjections)

	assert.Nil(suite.T(), err)
}

func (suite *InjectionTestSuite) TestInjectSetterFailed() {
	s := new(ServiceStub)
	err := suite.injector.inject(s, suite.wrongSetterInjections)

	assert.EqualError(suite.T(), err, "Cannot inject test.wrong , service not found.")
}

func (suite *InjectionTestSuite) TestInjectPropertyFailed() {
	s := new(ServiceStub)
	err := suite.injector.inject(s, suite.wrongPropInjections)

	assert.EqualError(suite.T(), err, "Cannot inject test.wrong , service not found.")
}

func TestInjectionTestSuite(t *testing.T) {
	suite.Run(t, new(InjectionTestSuite))
}
