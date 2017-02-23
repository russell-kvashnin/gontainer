package main

import (
	"github.com/russell-kvashnin/gontainer/example/constructor"
	"github.com/russell-kvashnin/gontainer/example/property"
	"github.com/russell-kvashnin/gontainer/example/setter"
)

func main() {
	constructor.ConstructorInjectionExamples()
	setter.SetterInjectionExamples()
	property.PropertyInjectionExamples()
}
