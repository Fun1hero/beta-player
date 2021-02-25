package main

import (
	"bufio"
	"fmt"
	"os"
	"log"
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
	fmt.Println(toPN, fromPN)

	fd, err := os.OpenFile(toPN, os.O_RDONLY, os.ModeNamedPipe) // opens toPN named pipe
	fd1, err1 := os.OpenFile(fromPN, os.O_RDWR, 0600)           // opens fromPN named pipe
	//for i := 0; i < 3; i++ {                           // reads 3 times
		input :=readFromPipe(fd, err)
		
	//}
	if err1 != nil {
		log.Fatal(err)
	}
	input2 := strings.SplitAfter(input, "\n")
	fmt.Print(input2)
	for i := 0; i < len(input2)-1; i++ {
		fmt.Println(i)
		stringSlice := strings.Split(input2[i], ":")
		stringSlice2 := strings.Split(stringSlice[1], ",")
		if stringSlice[0]=="04" {
			fmt.Println("Choose any two dice options from the following or choose A")
			for j :=0;j<len(stringSlice2);j++{
             fmt.Println(stringSlice2[j])
			}
		}
	}
	var Dice1 string
	fmt.Println("Choose first dice option")
	fmt.Scanf("%s", &Dice1)

		var Dice2 string	
         fmt.Println("Choose second dice option")
		fmt.Scanf("%s", &Dice2)
		var Terrain string	
		var Player string	

		fmt.Println("Choose Terrian")
		fmt.Scanf("%s", &Terrain)
		fmt.Println("Choose Player that you want to interrogate")
		fmt.Scanf("%s", &Player)

		writeToPipe(fd1, err1, "05:"+Dice1+","+Dice2+","+Terrain+","+Player) // writes 1 time

	fd.Close()
	fd1.Close()
}
