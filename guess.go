package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


var toPN string = "/tmp/toP"     // "toP + player number" when we (I) understand how we get that
var fromPN string = "/tmp/fromP" // same here "fromP + player number"
// Writes to "fromPN" named pipe


// Reads from "toPN" named pipe
func readFromPipe(fd *os.File, err error) string {
	if err != nil {
		fmt.Errorf(err.Error())
	}
	// var buff bytes.Buffer
	buff := make([]byte, 1024)
	n, err := fd.Read(buff)
	for n == 0 {
		n, err = fd.Read(buff)
	}
	if err != nil {
		fmt.Errorf(err.Error())
	}
	if len(buff) > 0 {
		return string(buff)
	}
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return ""
}

func writeToPipe(fd1 *os.File, err1 error, args string) {
	if err1 != nil {
		fmt.Errorf(err1.Error())
	}
	fmt.Println("string-"+args)
	fd1.Write([]byte(args))
}

func main() {
	fmt.Println("What player number are you: ")
	reader := bufio.NewReader(os.Stdin)
	playerNumber, _ := reader.ReadString('\n')
	playerNumber = strings.Replace(playerNumber, "\n", "", -1)
	fmt.Println()
	toPN += playerNumber
	fromPN += playerNumber
	fmt.Println("Do you want to guess? Answer with Y or N: ")
	
	var Response string
	fmt.Scanf("%s", &Response)

	if(Response == "Y"){
		fd, err := os.OpenFile(toPN, os.O_RDONLY, os.ModeNamedPipe) // opens toPN named pipe

		fd1, err1 := os.OpenFile(fromPN, os.O_RDWR, 0600)           // opens fromPN named pipe

		readFromPipe(fd, err)

		var first_token string
		var second_token string

		fmt.Println("Choose first token: ")
		fmt.Scanf("%s", &first_token)

		fmt.Println("Choose second token: ")
		fmt.Scanf("%s", &second_token)

		writeToPipe(fd1, err1, "07:P"+playerNumber+","+first_token+","+second_token) // writes 1 time


	}



}
