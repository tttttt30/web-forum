package system

import (
	"fmt"
	"log"
)

func ErrLog(funcName string, err error) error {
	fmtPrint := fmt.Sprintf("%s: %s", funcName, err)
	log.Println(fmtPrint)
	return fmt.Errorf(fmtPrint)
}

func FatalLog(funcName string, err error) {
	fmtPrint := fmt.Sprintf("%s: %s", funcName, err)
	log.Fatalln(fmtPrint)
}
