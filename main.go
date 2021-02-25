package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Player struct
type Player struct {
	no       int
	terrains []string
}

//Game struct
type Game struct {
	totalPlayers  int
	activePlayers int
	leftTokens    []string
}

var p Player
var g Game

var toPN string = "/tmp/toP"     // "toP + player number" when we (I) understand how we get that
var fromPN string = "/tmp/fromP" // same here "fromP + player number"

type fn func(string)

func selectedFunction(f fn, val string) { // selectedFunction provides functionality to call specific function by its id [:2] of args string
	f(val)
}

var functions = map[string]fn{
	"01": g.playerNO,
	"02": p.readMyTerrain,
	"03": g.leftoverTokens,
	// "04": toBeImplemented,
	// "05": toBeImplemented,
	// "06": toBeImplemented,
	// "07": toBeImplemented,
	// "08": toBeImplemented,
	// "09": toBeImplemented,
	// "10": toBeImplemented,
	// "11": toBeImplemented,
}

func (gm *Game) playerNO(args string) {
	gm.totalPlayers, _ = strconv.Atoi(args[len(args)-1:])
	gm.activePlayers = gm.totalPlayers
}

func (pl *Player) readMyTerrain(args string) {
	pl.terrains = strings.Split(args[3:], ",")
}

func (gm *Game) leftoverTokens(args string) {
	gm.leftTokens = strings.Split(args[3:], ",")
}

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

// Writes to "fromPN" named pipe
func writeToPipe(fd1 *os.File, err1 error, args string) {
	if err1 != nil {
		fmt.Errorf(err1.Error())
	}
	fd1.Write([]byte(args))
}

func main() {
	fmt.Println("What player number are you: ")
	reader := bufio.NewReader(os.Stdin)
	playerNumber, _ := reader.ReadString('\n')
	playerNumber = strings.Replace(playerNumber, "\n", "", -1)
	toPN += playerNumber
	fromPN += playerNumber
	fmt.Println(toPN, fromPN)

	fd, err := os.OpenFile(toPN, os.O_RDONLY, os.ModeNamedPipe) // opens toPN named pipe
	fd1, err1 := os.OpenFile(fromPN, os.O_RDWR, 0600)           // opens fromPN named pipe
	for i := 0; i < 3; i++ {                                    // reads 3 times
		fmt.Println(readFromPipe(fd, err))
	}
	writeToPipe(fd1, err1, "Hello") // writes 1 time
	fd.Close()
	fd1.Close()
}
