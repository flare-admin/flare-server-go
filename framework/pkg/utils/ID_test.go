package utils

import (
	"log"
	"testing"
)

const Id = "332623197205242271"

func Test_IsValidIDCard(t *testing.T) {
	if !IsValidIDCard(Id) {
		t.Error("id should be valid")
	}
}

func Test_GetIDCard(t *testing.T) {
	birthday, a, err := ExtractBirthdayAndAge(Id)
	if err != nil {
		t.Error(err)
	}
	log.Println(birthday, a)
}
