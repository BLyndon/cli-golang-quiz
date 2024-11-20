package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

var questions = []Question{
	{ID: 1, Question: "Where are the offices?", Options: []string{"Spain, Malta, Sweden", "Spain, Malta, France", "Spain, Malta, Germany"}},
	{ID: 2, Question: "What is the name of the CTO?", Options: []string{"Alice Smith", "Patrik Potocki", "Simon Lidzén"}},
	{ID: 3, Question: "What backend language transition was made?", Options: []string{"Java to Node.js", "Node.js to Java", "PHP to Go"}},
	{ID: 4, Question: "Where is the headquarter?", Options: []string{"Malta", "Barcelona", "Sweden"}},
	{ID: 5, Question: "What is the main product?", Options: []string{"CRM", "Risk Management Software", "A Database"}},
	{ID: 6, Question: "What happens on Wednesdays in Barcelona at the office?", Options: []string{"Yoga", "Spanish only day", "Everyone brings their dog to work"}},
}

var correctAnswers = map[int]int{
	1: 0,
	2: 1,
	3: 2,
	4: 0,
	5: 0,
	6: 1,
}

var scoreDistribution = map[int]int{
	0: 1,
	1: 2,
	2: 10,
	3: 18,
	4: 32,
	5: 20,
	6: 4,
}

func questionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func submissionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var answers []Answer
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	correctCount := 0
	var results []Result
	for _, answer := range answers {
		correct, exists := correctAnswers[answer.QuestionID]
		if !exists {
			http.Error(w, fmt.Sprintf("Question ID %d does not exist", answer.QuestionID), http.StatusBadRequest)
			return
		}
		isCorrect := answer.Answer == correct
		if isCorrect {
			correctCount++
		}
		results = append(results, Result{
			QuestionID: answer.QuestionID,
			Correct:    isCorrect,
		})
	}

	totalPlayers := 0
	playersWorse := 0
	for score, count := range scoreDistribution {
		totalPlayers += count
		if score < correctCount {
			playersWorse += count
		}
	}

	betterThanPercentage := 0.0
	if totalPlayers > 0 {
		betterThanPercentage = (float64(playersWorse) / float64(totalPlayers)) * 100
	}

	scoreDistribution[correctCount]++

	for i := range results {
		results[i].BetterThanPercentage = betterThanPercentage
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/questions", questionsHandler)
	http.HandleFunc("/submission", submissionHandler)
	fmt.Println("Backend is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
