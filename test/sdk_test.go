package test

import (
	apper "../../apper-go"
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	app, err := apper.GetApper()
	if err != nil {
		fmt.Print("err")
	}
	err = app.Connect("47.99.72.199:4222")
	if err != nil {
		fmt.Print("err")
	}
	txID, err := app.Start("./task_test.yaml")

	if app.Ready(txID) {
		inter, err := app.GetVal("abcd", txID)
		t := inter.([]string)
		fmt.Println(t, err)
	}
}
