package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chatgp/gpt3"
)

const ApiFileName = "chat-token.txt"

func main() {
	fmt.Println("Ensure you start a new console if you want a new conversations. Tokens are expensive!!")
	reader := bufio.NewReader(os.Stdin)
	apiKey := getKey(reader)

	cli, _ := gpt3.NewClient(&gpt3.Options{
		ApiKey:  apiKey,
		Timeout: 600 * time.Second,
		Debug:   false,
	})
	uri := "/v1/chat/completions"
	var m []Message
	fmt.Println("Ask me something")

	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		// request api
		addMessage(&m, "user", text)

		params := map[string]interface{}{
			"model":    "gpt-3.5-turbo",
			"messages": m,
		}

		res, err := cli.Post(uri, params)
		if err != nil {
			log.Fatalf("request api failed: %v", err)
		}

		message := res.Get("choices.0.message.content").String()
		if message == "" {
			fmt.Println("Something went wrong at Chat GPT's end")
			fmt.Println(res.Get("message"))
		}
		fmt.Printf(message)
		fmt.Println("")
		fmt.Println("--------")
		addMessage(&m, "assistant", message)
	}
}

func addMessage(m *[]Message, role string, message string) {
	dict := make(map[string]interface{})
	dict["role"] = role
	dict["content"] = message

	*m = append(*m, dict)
}

func getKey(reader *bufio.Reader) string {
	homeDir, err := os.UserHomeDir()

	_, err = os.Stat(homeDir + "/" + ApiFileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Create(homeDir + homeDir + "/" + ApiFileName)
	}

	readFile, err := os.Open(homeDir + "/" + ApiFileName)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	fileScanner.Scan()
	line := fileScanner.Text()
	var apiKey = line
	readFile.Close()
	if apiKey == "" {
		fmt.Println("Please provide your API Key")
		key, _ := reader.ReadString('\n')
		key = strings.Replace(key, "\n", "", -1)
		os.WriteFile(homeDir+"/"+ApiFileName, []byte(key), 777)
		apiKey = key
	}
	return apiKey
}

type Message map[string]interface{}
