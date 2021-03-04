package main

import (
	"bufio"
	"fmt"
	"io"
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
var toPN string = "/tmp/betatoP"     // "toP + player number" when we (I) understand how we get that
var fromPN string = "/tmp/betafromP" // same here "fromP + player number"

type fn func(string) string

func selectedFunction(f fn, val string) string { // selectedFunction provides functionality to call specific function by its id [:2] of args string
	return f(val)
}

var functions = map[string]fn{
	"01": g.playerNO,
	"02": p.readMyTerrain,
	"03": g.leftoverTokens,
	"04": playerTurn,
	"05": chooseDice,
	"06": SendInterrogation,
	"07": guessTokens,
	"08": guessCorrect,
	"09": guessIncorrect,
	"10": tokenInfoSwap,
	"11": remainingWinner,
}

func (gm *Game) playerNO(args string) string {
	gm.totalPlayers, _ = strconv.Atoi(args[len(args)-1:])
	gm.activePlayers = gm.totalPlayers
	return ""
}

func (pl *Player) readMyTerrain(args string) string {
	pl.terrains = strings.Split(args[3:], ",")
	return ""
}

func (gm *Game) leftoverTokens(args string) string {
	gm.leftTokens = strings.Split(args[3:], ",")
	return ""
}

func tokenInfoSwap(args string) string {
	message := strings.Split(args[3:], ",")

	if string(message[0][1]) == strconv.Itoa(p.no) {
		fmt.Printf("You let %s know you got a token %s\n", message[1], message[2])
	} else {
		fmt.Printf("You acknowledge %s got a token %s\n", message[0], message[2])
	}
	return ""
}

func remainingWinner(args string) string {
	message := strings.Split(args[3:], ",")

	fmt.Printf("%s wins as the only remaining player. All others have guessed incorrectly and been disqualified. The treasures are located at %s and %s\n",
		message[0], message[1], message[2])
	return ""
}

func playerTurn(args string) string {
	stringSlice := strings.Split(args, ":")
	stringSlice2 := strings.Split(stringSlice[1], ",")
	fmt.Println("\n")
	fmt.Println("Player "+stringSlice2[0][1:]+" has rolled "+stringSlice2[1]+","+stringSlice2[2]+","+stringSlice2[3])
	if("P"+strconv.Itoa(p.no)!=stringSlice2[0]){
	
		return "";
	}
	var response string
	fmt.Println("Choose any two dice options from the following or choose A")
	for j := 1; j < len(stringSlice2); j++ {
		fmt.Println(stringSlice2[j])
	}
	fmt.Println("Would you like to guess? Y/N")
	fmt.Scanln(&response)
	if response == "Y" || response == "y" {
		return guessTokens("")
	}
	return chooseDice(args)

}

func chooseDice(args string) string {

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

	var temp string = "05:" + Dice1 + "," + Dice2 + "," + Terrain + ",P" + Player
	return temp
}

func SendInterrogation(args string) string {

	stringSlice := strings.Split(args, ":")
	stringSlice2 := strings.Split(stringSlice[1], ",")

	fmt.Printf("Player %s asks %s how many locations they've searched between %s and %s in %s terrain.\n",
		stringSlice2[5], stringSlice2[4], stringSlice2[0], stringSlice2[1], stringSlice2[2])

	fmt.Printf("Player %s responds %s.\n",
		stringSlice2[4], stringSlice2[3])
	return ""
}

func guessTokens(playerNumber string) string {
	fmt.Println("fromPn-" + fromPN)
	var first_token string
	var second_token string
	fmt.Println("Choose first token: ")
	fmt.Scanf("%s", &first_token)
	fmt.Println("Choose second token: ")
	fmt.Scanf("%s", &second_token)
	var temp string = "07:P" + strconv.Itoa(p.no) + "," + first_token + "," + second_token
	return temp
}

func guessCorrect(args string) string {
	stringSlice := strings.Split(args, ":")
	stringSlice2 := strings.Split(stringSlice[1], ",")
	fmt.Printf("Player %s is correct! They have won the game.\n",
		stringSlice2[0])
	fmt.Printf("The treasures were located at %s and %s.\n",
		stringSlice2[1], stringSlice2[2])
	return ""
}

func guessIncorrect(args string) string {
	message := strings.Split(args, ":")
	fmt.Printf("Player %s is submitting a guess at the treasure locations!. Player %s was wrong. They are now disqualified from winning.\n",
		message[1], message[1])
	return ""
}

// Reads from "toPN" named pipe
func readFromPipe(fd *os.File, rd *bufio.Reader) string {
	buff, err := rd.ReadString('\n')
	if err == io.EOF {
		return "exit"
	}
	if err != nil {
		fmt.Errorf(err.Error())
	}
	if len(buff) > 0 {
		return buff
	}
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return ""
}

// Writes to "fromPN" named pipe
func writeToPipe(fd1 *os.File, args string) {
	fd1.Write([]byte(args))
}

func main() {
	fmt.Println("What player number are you: ")
	reader := bufio.NewReader(os.Stdin)
	playerNumber, _ := reader.ReadString('\n')
	playerNumber = strings.Replace(playerNumber, "\n", "", -1)
	p.no, _ = strconv.Atoi(playerNumber)
	toPN += playerNumber
	fromPN += playerNumber
	fmt.Println(toPN, fromPN)

	fd, err := os.OpenFile(toPN, os.O_RDONLY, os.ModeNamedPipe) // opens toPN named pipe
	if err != nil {
		fmt.Errorf(err.Error())
	}
	rd := bufio.NewReader(fd)

	fd1, err1 := os.OpenFile(fromPN, os.O_RDWR, 0600) // opens fromPN named pipe
	if err1 != nil {
		fmt.Errorf(err1.Error())
	}

	for {
		serverSaid := readFromPipe(fd, rd)
		if serverSaid == "exit" {
			break
		}
		playerReply := selectedFunction(functions[serverSaid[:2]], serverSaid)
		if playerReply != "" {
			writeToPipe(fd1, playerReply)
		}
	}
}
