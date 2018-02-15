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
	alph         = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	num          = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
)

func getRandomNumber(length int) int {
	return rand.Intn(length)
}

func getRandomNumberBetween(min, max int) int {
	return rand.Intn(max-min) + min
}

func getComputerName() string {
	dstSlice := make([]string, 5)
	srcSlice := make([]string, 5)
	srcSlice[0] = computerName[getRandomNumber(len(computerName))]
	srcSlice[1] = firstName[getRandomNumber(len(firstName))]
	srcSlice[2] = lastName[getRandomNumber(len(lastName))]
	srcSlice[3] = jobName[getRandomNumber(len(jobName))]
	srcSlice[4] = num[getRandomNumber(len(num))]
	perm := rand.Perm(5)
	for i, v := range perm {
		dstSlice[v] = srcSlice[i]
	}
	return strings.Join(dstSlice[:], "")
}

func randomString() string {
	strLen := getRandomNumberBetween(16, 32)
	compName := make([]string, strLen)
	for x := 0; x < strLen; x++ {
		numOrStr := getRandomNumber(2)
		upperOrLower := getRandomNumber(2)
		switch numOrStr {
		case 0:
			compName[x] = num[getRandomNumber(len(num))]
		case 1:
			switch upperOrLower {
			case 0:
				compName[x] = alph[getRandomNumber(len(alph))]
			case 1:
				compName[x] = strings.ToUpper(alph[getRandomNumber(len(alph))])
			}
		}
	}
	return strings.Join(compName, "")
}
