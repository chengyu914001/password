package passwordGenerator

import (
	"regexp"
	"sort"
)  

type PasswordFinder struct {
	table []byte
	reg string
	isRunning bool
}

func (obj *PasswordFinder) Init() (*PasswordFinder) {
	if obj.isRunning {
		return obj
	}

	obj.AddTableAllLowercase().AddTableAllNumber().SetReg(".*")

	return obj
}

// Get var --------------------------------------------------------------------------
func (obj *PasswordFinder) IsRunning() bool {
	return obj.isRunning
}

func (obj *PasswordFinder) GetTable() string {
	return string(obj.table)
}

func (obj *PasswordFinder) GetReg() string {
	return obj.reg
}

// Set var --------------------------------------------------------------------------
//	-- Table --
func (obj *PasswordFinder) SetTableClear() (*PasswordFinder) {
	if obj.isRunning{
		return obj
	}

	obj.table = make([]byte, 0)

	return obj
}

func (obj *PasswordFinder) AddTable(s string) (*PasswordFinder) {
	if obj.isRunning{
		return obj
	}

	byte_s := []byte(s)
	sort.Slice(byte_s, func(i, j int) bool {
		return byte_s[i] < byte_s[j]
	})

	i, j := 0, 0
	res := make([]byte, 0, len(obj.table) + len(s))
	for i < len(byte_s) && j < len(obj.table) {
		if byte_s[i] < obj.table[j] {
			res = append(res, byte_s[i])
			i++
	
		}else if byte_s[i] > obj.table[j] {
			res = append(res, obj.table[j])
			j++

		}else {
			res = append(res, byte_s[i])
			i++
			j++
		}
	}
	res = append(res, byte_s[i:]...)
	res = append(res, obj.table[j:]...)

	for i = 0; res[i] < ' '; i++ {}
	res = res[i:]

	for i = 0; i < len(res); i++ {
		for j = i + 1; j < len(res) && res[i] == res[j]; j++ {}
		res = append(res[:i+1], res[j:]...)

		if res[i] > '~' {
			res = res[:i]
			break
		}
	}
	
	func (a, b *[]byte) {
		tmp := *a
		*a = *b
		*b = tmp
	}(&res, &obj.table)

	return obj
} 

func (obj *PasswordFinder) AddTableAllNumber() (*PasswordFinder) {
	obj.AddTable("0123456789")

	return obj
}

func (obj *PasswordFinder) AddTableAllLowercase() (*PasswordFinder) {
	obj.AddTable("abcdefghijklmnopqrstuvwxyz")

	return obj
}

func (obj *PasswordFinder) AddTableAllUppercase() (*PasswordFinder) {
	obj.AddTable("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	return obj
}

func (obj *PasswordFinder) AddTableAllSpecial() (*PasswordFinder) {
	obj.AddTable(" !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")

	return obj
}

//	-- others --

func (obj *PasswordFinder) SetReg(reg string) (*PasswordFinder) {
	if obj.isRunning {
		return obj
	}

	obj.reg = reg

	return obj
}

// Run ---------------------------------------------------------------------------------
func (obj *PasswordFinder) Stop() {
	obj.isRunning = false
}

func (obj *PasswordFinder) PasswordGenerator(passwordLen uint8) (<-chan string) {
	if obj.isRunning {
		return nil
	}
	obj.isRunning = true

	res := make(chan string)
	go func(ch chan<- string){
		var passwordGenerator func([]byte, []byte, int, *bool, *string, chan<- string)
		passwordGenerator = func (password []byte, table []byte, passwordIdx int, isRunning *bool, reg *string, res chan<- string){
			if !*isRunning || passwordIdx >= len(password) { 
				return 
			}
		
			if passwordIdx == len(password) - 1 {
				for i := 0; i < len(table); i++ {
					if !*isRunning {
						return
					}
		
					password[passwordIdx] = table[i]
					s := string(password)
					match, _ := regexp.MatchString(*reg, s)
					if match {
						res <- s
					}
				}
		
			} else {
				for i := 0; i < len(table); i++ {
					if !*isRunning {
						return
					}
		
					password[passwordIdx] = table[i]
					passwordGenerator(password, table, passwordIdx + 1, isRunning, reg, res)
				}
			}
		}
		passwordGenerator(make([]byte, passwordLen), obj.table, 0, &obj.isRunning, &obj.reg, ch)
		obj.isRunning = false
		close(res)
		
	}(res)
	
	return res
}

// Private method -------------------------------------------------------------------------

