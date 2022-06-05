package PasswordGeneratorTest

import (
	"findPassword/pkg/password"
	"sync"
	"testing"
)


func TestSetTableLen(t *testing.T){
	passwordFinder := new(passwordGenerator.PasswordFinder)
	passwordFinder.Init().SetTableClear().AddTableAllLowercase().AddTableAllNumber().AddTableAllUppercase().AddTableAllUppercase().AddTableAllSpecial()
	passwordFinder.AddTable(string([]byte{0, 1, 1, 1, 2, 127, 127}))
	t.Log(passwordFinder.GetTable())

	if len(passwordFinder.GetTable()) != 95 {
		t.Error(len(passwordFinder.GetTable()), "!=", 95)
	}
}

func TestSetReg(t *testing.T){
	passwordFinder := new(passwordGenerator.PasswordFinder)
	passwordFinder.Init().AddTableAllUppercase().SetReg("[0-9][a-z][A-Z]")
	passwordGenerator := passwordFinder.PasswordGenerator(3)
	counter := 0
	for password := range passwordGenerator {
		go func(s string){
		}(password)

		counter++
	}

	if counter != 10*26*26 {
		t.Error(counter, "!=", 10*26*26)
	}
}

func TestStop(t *testing.T) {
	passwordFinder := new(passwordGenerator.PasswordFinder)
	passwordFinder.Init()
	passwordGenerator := passwordFinder.PasswordGenerator(3)
	wg := &sync.WaitGroup{}

	counter := 0
	for password := range passwordGenerator {
		wg.Add(1)
		go func(s string, wg *sync.WaitGroup){
			defer wg.Done()
			t.Log(s)
		}(password, wg)

		counter++
		if counter == 10 {
			passwordFinder.Stop()
			break
		}
	}

	if passwordFinder.IsRunning() {
		t.Error("cannot stop!")
	}
}

func TestGeneratorOutput(t *testing.T){
	passwordFinder := new(passwordGenerator.PasswordFinder)
	passwordFinder.Init()
	wg := &sync.WaitGroup{}
	mutex := &sync.RWMutex{}
	counter := 0

	passwordGenerator := passwordFinder.PasswordGenerator(3)

	for password := range passwordGenerator {
		wg.Add(1)
		go func(s string, wg *sync.WaitGroup){
			defer wg.Done()
			t.Log(s)
			mutex.Lock()
			counter++
			mutex.Unlock()
		}(password, wg)

	}
	wg.Wait()

	if counter != 46656 {
		t.Error(counter, "!=", 46656)
	}
}