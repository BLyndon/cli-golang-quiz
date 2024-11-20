package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type Answer struct {
	QuestionID int `json:"question_id"`
	Answer     int `json:"answer"`
}

type Result struct {
	QuestionID           int     `json:"question_id"`
	Correct              bool    `json:"correct"`
	BetterThanPercentage float64 `json:"better_than_percentage"`
}

var questions []Question
var answers map[int]int

var startQuizCmd = &cobra.Command{
	Use:   "start-quiz",
	Short: "Start the quiz and answer questions",
	Run: func(cmd *cobra.Command, args []string) {
		questionsEndpoint := "http://localhost:8080/questions"

		resp, err := http.Get(questionsEndpoint)
		if err != nil {
			fmt.Printf("Error fetching questions: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to fetch questions: HTTP %d\n", resp.StatusCode)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			return
		}

		err = json.Unmarshal(body, &questions)
		if err != nil {
			fmt.Printf("Error parsing questions: %v\n", err)
			return
		}

		answers = make(map[int]int)
		for _, q := range questions {
			fmt.Printf("\nQuestion %d: %s\n", q.ID, q.Question)
			for i, option := range q.Options {
				fmt.Printf("  %d. %s\n", i+1, option)
			}

			var userAnswer int
			for {
				fmt.Printf("Your answer (1-%d): ", len(q.Options))
				_, err := fmt.Scan(&userAnswer)
				if err == nil && userAnswer >= 1 && userAnswer <= len(q.Options) {
					answers[q.ID] = userAnswer - 1
					break
				} else {
					fmt.Println("Invalid input. Please try again.")
				}
			}
		}

		submitAnswers(answers)
	},
}

func submitAnswers(answers map[int]int) {
	answersEndpoint := "http://localhost:8080/submission"

	var answerList []Answer
	for id, ans := range answers {
		answerList = append(answerList, Answer{QuestionID: id, Answer: ans})
	}
	jsonData, err := json.Marshal(answerList)
	if err != nil {
		fmt.Printf("Error serializing answers: %v\n", err)
		return
	}

	resp, err := http.Post(answersEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error submitting answers: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var results []Result
		err := json.Unmarshal(body, &results)
		if err != nil {
			fmt.Printf("Error parsing results: %v\n", err)
			return
		}

		fmt.Printf("\nQuiz Results:\n")
		for _, result := range results {
			correctStr := "Incorrect"
			if result.Correct {
				correctStr = "Correct"
			}
			fmt.Printf("Question %d: %s (Better than %.2f%% of players)\n",
				result.QuestionID, correctStr, result.BetterThanPercentage)
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to submit answers: HTTP %d\n%s\n", resp.StatusCode, string(body))
	}
}

func init() {
	rootCmd.AddCommand(startQuizCmd)
}
