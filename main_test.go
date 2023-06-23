package main

import (
	"testing"
)

func Testwritting_into_file(t *testing.T) {
	var r ResponseDetails
	r.Name = "Tech"
	var r1 []ResponseDetails
	r1 = append(r1, r)
	err := writting_into_file(&r1)
	if err != nil {
		t.Errorf("Error returned is: %s.", err)
	}
}
func TestloadConfig(t *testing.T) {

	_, err := loadConfig(`D:\GoWorkspace\src\Golang_code_techmahindra`)
	if err != nil {
		t.Errorf("Error returned is: %s.", err)
	}
}
