package passwordGenerator

import (
	"log"
	"os/exec"
	"sync"
)

// install path -------------------------------------------------------------------------
var path7z string

// set install path ----------------------------------------------------------------------
func Set7zPath(path string) {
	path7z = path
}

// Public mathos -------------------------------------------------------------------------
func (obj *PasswordFinder)Find7ZPassword(startLen uint8, endLen uint8, filename string, runSize uint8) string{
	log.Println("filename:" ,filename)
	f := func(password string) bool {
		err := exec.Command(path7z, "t", filename, "-p" + password).Run()
		return err == nil
	}
	return obj.findPassword(startLen, endLen, filename, f, runSize)
}

// Private method -------------------------------------------------------------------------
type funcExe func(string) bool

func (obj *PasswordFinder)findPassword(startLen uint8, endLen uint8, filename string, exe funcExe, runSize uint8) string{
	wg := sync.WaitGroup{}
	runBuffer := make(chan struct{}, runSize)
	result := ""
	for i := startLen; i <= endLen; i++ {
		passwordGenerator := obj.PasswordGenerator(i)
		for password := range passwordGenerator {

			wg.Add(1)
			runBuffer<-struct{}{}
			go func(password string, wg *sync.WaitGroup, res *string, runBuffer <-chan struct{}){
				defer func () {
					<-runBuffer
					wg.Done()
				}()

				if exe(password) {
					obj.Stop()

					result = password
				}
			}(password, &wg, &result, runBuffer)
		}
		wg.Wait()
		
		if result != "" {
			return result
		}
	}

	return result
}

