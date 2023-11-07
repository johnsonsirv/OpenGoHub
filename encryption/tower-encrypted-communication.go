package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var encOffset int32 = 3

func main()  {
	stringChan := make(chan string)
	cellTower1Chan := make(chan string)
	cellTower2Chan := make(chan string)

	go cellTower1(stringChan, cellTower1Chan, encOffset)
	go cellTower2(stringChan, cellTower2Chan, encOffset)

	for i := 0; i < 2; i++ {
		select {
		case msg := <- cellTower1Chan:
			fmt.Printf("\nControlTower::Message from T1:: %v", msg)
		case msg := <- cellTower2Chan:
			fmt.Printf("\nControl Tower::Message from T2:: %v", msg)
		}
	}
}


func cellTower1(s chan string, t1 chan string, offset int32)  {
	stream := bufio.NewReader(os.Stdin)
	fmt.Println("T1::Enter message for cell T2: ")
	
	streamInput, _ := stream.ReadString('\n')
	streamInput = strings.Replace(streamInput, "\r\n", "", -1)

	encrypted := encryptSync(streamInput)

	fmt.Printf("\nT1::Encrypted message: %s", encrypted)
	s <- encrypted
	t1 <- "Message sent to T2"
}

func cellTower2(s chan string, t2 chan string, offset int32)  {
	encMsg := <- s

	decrypted := decryptSync(encMsg)

	fmt.Printf("\nT2::Decrypted Message: %s", decrypted)
	t2 <- "Message received from T1"
}

func encryptSync(s string) (string) {
	var enc string
	for _, c := range s {
		enc += string(c + encOffset)
	}

	return enc
}

func decryptSync(s string) (string)  {
	var dec string
	
	for _, c := range s {
		dec += string(c - encOffset)
	}

	return dec
}