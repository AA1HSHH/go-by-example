package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	maxNum := 100
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Intn(maxNum)
	// fmt.Println("The secret number is ", secretNumber)

	fmt.Println("Please input your guess")
	//reader := bufio.NewReader(os.Stdin)
	for {
		//input, err := reader.ReadString('\n')
		//if err != nil {
		//	fmt.Println("An error occured while reading input. Please try again", err)
		//	continue
		//}
		//input = strings.Trim(input, "\r\n")
		input := ""
		fmt.Scanf("%s", &input)
		guess, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter an integer value")
			continue
		}
		fmt.Println("You guess is", guess)
		if guess > secretNumber {
			fmt.Println("Your guess is bigger than the secret number. Please try again")
		} else if guess < secretNumber {
			fmt.Println("Your guess is smaller than the secret number. Please try again")
		} else {
			fmt.Println("Correct, you Legend!")
			break
		}
	}
}