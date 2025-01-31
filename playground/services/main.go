package main

import (
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/service"
)

func main() {
	env := &environment.EnvironmentVariables{}
	serviceInstance := service.New(env)

	serviceInstance.JobService.DoSomething()
}
