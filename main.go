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

type fn func(string) string

func selectedFunction(f fn, val string) { // selectedFunction provides functionality to call specific function by its id [:2] of args string
	f(val)
}

var functions = map[string]  fn {
	"01": g.playerNO,
	"02": p.readMyTerrain,
	"03": g.leftoverTokens,
	"04": playerTurn,
	"05": chooseDice,
	"06": SendInterrogation,
	//"07": tobeImplemented,
//	"08": tobeImplemented,
//"09": tobeImplemented,
	 "10": tokenInfoSwap,
	 "11": remainingWinner,
}

func (gm *Game) playerNO(args string) string{
	gm.totalPlayers, _ = strconv.Atoi(args[len(args)-1:])
	gm.activePlayers = gm.totalPlayers
	return "";
}

func (pl *Player) readMyTerrain(args string) string{
	pl.terrains = strings.Split(args[3:], ",")
	return ""; 
}

func (gm *Game) leftoverTokens(args string) string{
	gm.leftTokens = strings.Split(args[3:], ",")
	return "";
}

func tokenInfoSwap(args string) string{
	message := strings.Split(args[3:], ",")

	if string(message[0][1]) == strconv.Itoa(p.no) {
		fmt.Printf("You let %s know you got a token %s\n", message[1], message[2])
	} else {
		fmt.Printf("You acknowledge %s got a token %s\n", message[0], message[2])
	}
	return ""
}

func remainingWinner(args string) string{
	fmt.Println("string-"+args)
	message := strings.Split(args[3:], ",")

	fmt.Printf("%s wins as the only remaining player. All others have guessed incorrectly and been disqualified. The treasures are located at %s and %s\n",
		message[0], message[1], message[2])
		return ""
}



func playerTurn(args string) string{
	stringSlice := strings.Split(args, ":")
		stringSlice2 := strings.Split(stringSlice[1], ",")

		fmt.Println("Choose any two dice options from the following or choose A")
			for j :=1;j<len(stringSlice2);j++{
             fmt.Println(stringSlice2[j])
			}
	return ""
}

func chooseDice(args string) string{


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

		var temp string = "05:"+Dice1+","+Dice2+","+Terrain+",P"+Player
		return temp
}
	
	
func SendInterrogation(args string) string{
	
	stringSlice := strings.Split(args, ":")
		stringSlice2 := strings.Split(stringSlice[1], ",")

		fmt.Printf("Player %s asks %s how many locations they've searched between %s and %s in %s terrain.\n",
		stringSlice2[5],stringSlice2[4],stringSlice2[0],stringSlice2[1],stringSlice2[2])
		
		fmt.Printf("Player %s responds %s.\n",
		stringSlice2[4],stringSlice2[3])
		return ""
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
