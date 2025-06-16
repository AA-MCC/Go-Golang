package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

}
