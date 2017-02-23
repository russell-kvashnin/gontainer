package main

import (
	"github.com/russell-kvashnin/gontainer/example/constructor"
	"github.com/russell-kvashnin/gontainer/example/setter"
	"github.com/russell-kvashnin/gontainer/example/property"
)

func main() {
	constructor.ConstructorInjectionExamples()
	setter.SetterInjectionExamples()
	property.PropertyInjectionExamples()
}