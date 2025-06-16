package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var pl = fmt.Println

func main() {
	pl("Hello Go \n")
	pl("What is your name?")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err == nil {
		pl("Hello", name)
	} else {
		log.Fatal(err)
	}

	//var name type   -- Note that  if a variable function or type starts with a CAPITAL letter it is considered EXPORTED, and can be accessed outside the package
	// camelCase is used

	var vName string = "Aaron"
	var v1, v2 = 1.2, 3.4
	var v3 = "hello" //figures out the type for us
	v4 := 2.4        // shortcut when initially assigning value
	v4 = 5.4
	pl(vName, v1, v2, v3, v4)

	// int, float64, bool, string, rune
	// Default type 0,0.0, false, ""
	pl(reflect.TypeOf(25))
	pl(reflect.TypeOf(3.14))
	pl(reflect.TypeOf(true))
	pl(reflect.TypeOf("Hello"))

	cV1 := 1.5
	cV2 := int(cV1)
	pl(cV2)

	cV3 := "50000000"
	cV4, err := strconv.Atoi(cV3) // to Ascii to integer
	pl(cV4, err, reflect.TypeOf(cV4))
	cV5 := cV4
	cV6 := strconv.Itoa(cV5) // integer to ascii
	pl(cV6)

	cV7 := "3.14"
	if cV8, err := strconv.ParseFloat(cV7, 64); err == nil {
		pl(cV8)
	}
	cV9 := fmt.Sprintf("%f", 3.14)
	pl(cV9)

	//Conditional operators - standard < > >= <= == !=     Logical operators - && || !

	sV1 := "A word"
	replacer := strings.NewReplacer("A", "Another")
	sV2 := replacer.Replace(sV1)
	pl(sV2)
	pl("Length :", len(sV2))
	pl("Contains Another :", strings.Contains(sV2, "Another"))
	pl("o index :", strings.Index(sV2, "o"))
	pl("Replace :", strings.Replace(sV2, "o", "0", -1))
	sV3 := "\nSome Words\n"
	sV3 = strings.TrimSpace(sV3)
	pl("Split :", strings.Split("a-b-c-d", "-"))
	pl("Lower :", strings.ToLower(sV2))
	pl("Upper :", strings.ToUpper(sV2))
	pl("Prefix :", strings.HasPrefix("tacocat", "taco"))
	pl("Suffix :", strings.HasSuffix("tacocat", "taco"))

}
