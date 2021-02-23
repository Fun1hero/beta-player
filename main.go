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

func main() {
	fmt.Println("What player number are you: ")
	reader := bufio.NewReader(os.Stdin)
	playerNumber, _ := reader.ReadString('\n')
	playerNumber = strings.Replace(playerNumber, "\n", "", -1)
	toPN += playerNumber
	fromPN += playerNumber
	fmt.Println(toPN, fromPN)
	selectedFunction(functions["01"], "01:03")
	fmt.Println(g.totalPlayers)
}
