package main

import (
	"math/rand"
	"strings"
)

var (
	computerName = []string{"laptop", "workstation", "server"}
	firstName    = []string{"alice", "bob", "john"}
	lastName     = []string{"smith", "johnson", "williams"}
	jobName      = []string{"dev", "qa", "ops", "hr"}
)

func getRandomNumber(length int) int {
	return rand.Intn(length)
}

func getComputerName() string {
	dstSlice := make([]string, 4)
	srcSlice := make([]string, 4)
	srcSlice[0] = computerName[getRandomNumber(len(computerName))]
	srcSlice[1] = firstName[getRandomNumber(len(firstName))]
	srcSlice[2] = lastName[getRandomNumber(len(lastName))]
	srcSlice[3] = jobName[getRandomNumber(len(jobName))]
	perm := rand.Perm(4)
	for i, v := range perm {
		dstSlice[v] = srcSlice[i]
	}
	return strings.Join(dstSlice[:], "")
}
