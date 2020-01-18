package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var quitLoop = false // this bool flips to true if the user requests to quit early
var questionFiles = []string{"ports.json", "webservers.json", "acronyms.json"}

// DataJSON is used to parse the parent JSON object file
type DataJSON struct {
	Questions []QuestionJSON `json:"Questions"`
}

// QuestionJSON is used to parse the question data
type QuestionJSON struct {
	Question string       `json:"Question"`
	Answers  []AnswerJSON `json:"Answers"`
}

// AnswerJSON is used to parse the answer data
type AnswerJSON struct {
	Answer    string `json:"Answer"`
	IsCorrect string `json:"isCorrect"`
}

var allQuestions []*QuestionData // main question bank used to randomly obtain our next question
var userAnswers []*UserAnswer    // slice of all the questions that have already been answered

// QuestionData is the object that the json is unmarshal'd into
type QuestionData struct {
	Question string
	Answers  []*AnswerData
}

// AnswerData is the object that the json is unmarshal'd into
type AnswerData struct {
	Answer    string
	IsCorrect bool
}

// UserAnswer tracks the result of the answered question by the user
type UserAnswer struct {
	Question          *QuestionData
	AnsweredCorrectly bool
}

func populateQuestionPool() {
	for _, v := range questionFiles {
		// read file in and unmarshal into objects
		fmt.Println("[DEBUG] Parsing current file:", v)
		file, _ := ioutil.ReadFile(fmt.Sprintf("data/questions/%s", v))
		data := DataJSON{}

		err := json.Unmarshal([]byte(file), &data)
		if err != nil {
			fmt.Println("[!] JSON failed to load!")
			fmt.Println(err)
			os.Exit(1)
		}

		for _, d := range data.Questions {
			newQuestion := &QuestionData{}
			newQuestion.Question = d.Question

			for _, a := range d.Answers {
				newAnswer := &AnswerData{}
				newAnswer.Answer = a.Answer

				b, err := strconv.ParseBool(a.IsCorrect)
				if err != nil {
					fmt.Println("Error parsing answers for question ->", d.Question)
					continue
				}
				newAnswer.IsCorrect = b
				newQuestion.Answers = append(newQuestion.Answers, newAnswer)
			}

			// todo: shuffle the order the answers are in the memory before doing anything with them
			rand.Shuffle(len(newQuestion.Answers), func(i, j int) {
				newQuestion.Answers[i], newQuestion.Answers[j] = newQuestion.Answers[j], newQuestion.Answers[i]
			})

			// the question should be fully serialized into an object now
			allQuestions = append(allQuestions, newQuestion)
		}
	}

	fmt.Println("Total questions loaded:", len(allQuestions))
}

func getRandomQuestion() *QuestionData {
	// get random question from the pool
	randomElement := rand.Intn(len(allQuestions))
	randomQuestion := allQuestions[randomElement]
	allQuestionsAfterRemoval := append(allQuestions[:randomElement], allQuestions[randomElement+1:]...)
	allQuestions = allQuestionsAfterRemoval
	return randomQuestion
}

func produceQuestion(q *QuestionData) {
	toPrint := fmt.Sprintf(`[Q] %s
	
1) %s
2) %s
3) %s
4) %s
`, q.Question, q.Answers[0].Answer, q.Answers[1].Answer, q.Answers[2].Answer, q.Answers[3].Answer)

	fmt.Println("")
	fmt.Printf(toPrint)
	fmt.Println("")
}

func handleAnswer(q *QuestionData) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Answer (1 - 4): ")
		userInput, _ := reader.ReadString('\n')
		// cause windows and macs don't like each other
		userInput = strings.Replace(userInput, "\n", "", -1)
		userInput = strings.Replace(userInput, "\r", "", -1)

		if userInputInt, err := strconv.Atoi(userInput); err == nil {
			if userInputInt < 1 || userInputInt > 4 {
				fmt.Println("Input does not look like a number between 1 and 4")
			} else {
				// everything should be okay if we are here, clean and break the loop
				ua := &UserAnswer{}
				ua.Question = q
				ua.AnsweredCorrectly = q.Answers[userInputInt-1].IsCorrect
				userAnswers = append(userAnswers, ua)

				if !ua.AnsweredCorrectly {
					showCorrectAnswer(q)
				} else {
					fmt.Println("[+] Correct!")
				}

				break
			}
		} else if strings.ToLower(userInput) == "quit" || strings.ToLower(userInput) == "q" {
			quitLoop = true
			break // quit out the logic loop
		} else {
			fmt.Println("[!] User's answer does not appear to be a number between 1 and 4")
		}
	}
}

func showCorrectAnswer(q *QuestionData) {
	for _, v := range q.Answers {
		if v.IsCorrect {
			fmt.Println("[TIP] The correct answer is:", v.Answer)
			break
		}
	}
}

func showFinalResults() {
	totalQuestionsAsked := len(userAnswers)
	totalCorrect := 0

	for _, v := range userAnswers {
		if v.AnsweredCorrectly {
			totalCorrect++
		}
	}

	percentage := float32(0)
	if totalQuestionsAsked > 0 {
		percentage = float32(totalCorrect) / float32(totalQuestionsAsked) * 100.00
	}

	fmt.Println("")
	fmt.Println("[*] Results:")
	if totalQuestionsAsked > 0 {
		fmt.Printf("\t%d / %d (%.1f%%)\n", totalCorrect, totalQuestionsAsked, percentage)
	} else {
		fmt.Println("No questions were answered")
	}
}

func logicLoop() {
	// while
	for {
		if !quitLoop {
			if len(allQuestions) > 0 {
				question := getRandomQuestion()
				produceQuestion(question)
				handleAnswer(question)
			} else {
				fmt.Println("[DEBUG] Question bank is empty")
				break
			}
		} else {
			fmt.Println("[DEBUG] Player requested to quit")
			break
		}
	}

	showFinalResults()
}

func main() {
	fmt.Println("Let's do the learning thing")
	fmt.Println("Made by Gary @monobehaviour, questions were obtained from general internet surfing")

	fmt.Println("")
	fmt.Println("[*] Loading question bank...")

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	populateQuestionPool()

	fmt.Println("")
	fmt.Println("[*] Starting the questions... (answer with 'quit' at any time to leave)")
	fmt.Println("")

	logicLoop()
}
