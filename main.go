package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Law     string   `json:"law"`
	Options []Option `json:"options"`
}

type Option struct {
	Text  string `json:"text"`
	Score int    `json:"score"`
}

type QuizSubmission struct {
	Answers map[string]int `json:"answers"`
}

type ScoreResult struct {
	TotalScore    int     `json:"totalScore"`
	MaxScore      int     `json:"maxScore"`
	Percentage    float64 `json:"percentage"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	LawBreakdown  []LawScore `json:"lawBreakdown"`
}

type LawScore struct {
	Law       string `json:"law"`
	Score     int    `json:"score"`
	MaxScore  int    `json:"maxScore"`
}

var questions = []Question{
	{
		ID:   1,
		Text: "You're about to give a presentation and realize you forgot to prepare one slide. You:",
		Law:  "L1: Confidence",
		Options: []Option{
			{Text: "Panic and apologize profusely to the audience", Score: 1},
			{Text: "Nervously skip over it hoping no one notices", Score: 2},
			{Text: "Acknowledge it briefly and move on smoothly", Score: 3},
			{Text: "Improvise confidently as if it was always the plan", Score: 4},
		},
	},
	{
		ID:   2,
		Text: "Someone questions your decision in front of others. You:",
		Law:  "L1: Confidence",
		Options: []Option{
			{Text: "Immediately second-guess yourself and backtrack", Score: 1},
			{Text: "Get defensive and argue your point aggressively", Score: 2},
			{Text: "Calmly explain your reasoning once", Score: 3},
			{Text: "Stand firm with quiet certainty, no explanation needed", Score: 4},
		},
	},
	{
		ID:   3,
		Text: "You spent 6 hours perfecting your outfit for an event. When someone asks about it, you:",
		Law:  "L2: Effortless",
		Options: []Option{
			{Text: "Tell them exactly how long it took and every detail", Score: 1},
			{Text: "Humble brag about finding it last minute", Score: 2},
			{Text: "Casually say 'oh, this old thing?'", Score: 3},
			{Text: "Simply smile and accept the compliment gracefully", Score: 4},
		},
	},
	{
		ID:   4,
		Text: "You achieved something impressive after months of hard work. On social media, you:",
		Law:  "L2: Effortless",
		Options: []Option{
			{Text: "Post a detailed journey of every struggle", Score: 1},
			{Text: "Share it with a long caption about the grind", Score: 2},
			{Text: "Post a simple photo with minimal caption", Score: 3},
			{Text: "Let others discover it naturally, or don't post at all", Score: 4},
		},
	},
	{
		ID:   5,
		Text: "You catch yourself nerding out about your favorite hobby. You:",
		Law:  "L3: Balance",
		Options: []Option{
			{Text: "Feel embarrassed and quickly change the subject", Score: 1},
			{Text: "Apologize for being 'such a nerd'", Score: 2},
			{Text: "Own it but read the room for interest", Score: 3},
			{Text: "Embrace it fully - your passion is magnetic", Score: 4},
		},
	},
	{
		ID:   6,
		Text: "At a party, someone asks what you do for fun. You:",
		Law:  "L3: Balance",
		Options: []Option{
			{Text: "List only 'cool' hobbies and hide your real interests", Score: 1},
			{Text: "Mention safe, generic activities", Score: 2},
			{Text: "Share your genuine interests with confidence", Score: 3},
			{Text: "Make even your quirkiest hobby sound intriguing", Score: 4},
		},
	},
	{
		ID:   7,
		Text: "There's an awkward silence in a conversation. You:",
		Law:  "L4: Mystery",
		Options: []Option{
			{Text: "Frantically fill it with nervous chatter", Score: 1},
			{Text: "Make a self-deprecating joke to ease tension", Score: 2},
			{Text: "Wait comfortably for someone else to speak", Score: 3},
			{Text: "Let the silence work for you - it builds intrigue", Score: 4},
		},
	},
	{
		ID:   8,
		Text: "Someone asks about your weekend plans. You:",
		Law:  "L4: Mystery",
		Options: []Option{
			{Text: "Give a detailed hour-by-hour breakdown", Score: 1},
			{Text: "Over-explain to fill the conversation", Score: 2},
			{Text: "Give a brief, honest answer", Score: 3},
			{Text: "Keep it vague and intriguing - 'I have a few things lined up'", Score: 4},
		},
	},
	{
		ID:   9,
		Text: "You realize your current image no longer represents who you want to be. You:",
		Law:  "L5: Reinvention",
		Options: []Option{
			{Text: "Feel stuck - 'this is just who I am'", Score: 1},
			{Text: "Make small changes but worry what others will think", Score: 2},
			{Text: "Gradually evolve while staying true to yourself", Score: 3},
			{Text: "Boldly reinvent yourself without seeking permission", Score: 4},
		},
	},
	{
		ID:   10,
		Text: "How do you set standards for how others treat you?",
		Law:  "L6: Royalty",
		Options: []Option{
			{Text: "Accept whatever treatment comes your way", Score: 1},
			{Text: "Hint at boundaries but rarely enforce them", Score: 2},
			{Text: "Communicate boundaries when crossed", Score: 3},
			{Text: "Your energy naturally commands respect - no words needed", Score: 4},
		},
	},
}

func getScoreTitle(percentage float64) (string, string) {
	switch {
	case percentage >= 90:
		return "Supreme Baddie", "You've mastered all the laws. Your presence commands attention and respect effortlessly. You embody confidence, mystery, and royalty in everything you do."
	case percentage >= 75:
		return "Certified Baddie", "You're well on your way to baddie mastery. You understand the power of confidence and mystery, with just a few areas to refine."
	case percentage >= 60:
		return "Rising Baddie", "You have strong baddie instincts but sometimes let doubt creep in. Keep trusting yourself and saying less - your glow up is in progress."
	case percentage >= 45:
		return "Baddie in Training", "You're learning the laws but haven't fully internalized them yet. Focus on building unshakeable confidence and embracing your authentic self."
	case percentage >= 30:
		return "Aspiring Baddie", "You have potential but often prioritize others' perceptions over your own power. Time to start treating yourself like the prize you are."
	default:
		return "Baddie Novice", "Your baddie journey is just beginning. Study the laws, practice confidence, and remember: you deserve to take up space."
	}
}

func calculateScore(answers map[string]int) ScoreResult {
	totalScore := 0
	maxScore := len(questions) * 4

	lawScores := make(map[string]int)
	lawMaxScores := make(map[string]int)

	for _, q := range questions {
		lawMaxScores[q.Law] += 4
		// Answer keys are question IDs as strings (e.g., "1", "2", ..., "10")
		key := strconv.Itoa(q.ID)
		if score, exists := answers[key]; exists {
			totalScore += score
			lawScores[q.Law] += score
		}
	}

	percentage := float64(totalScore) / float64(maxScore) * 100
	title, description := getScoreTitle(percentage)

	var breakdown []LawScore
	lawOrder := []string{"L1: Confidence", "L2: Effortless", "L3: Balance", "L4: Mystery", "L5: Reinvention", "L6: Royalty"}
	for _, law := range lawOrder {
		if maxS, exists := lawMaxScores[law]; exists {
			breakdown = append(breakdown, LawScore{
				Law:      law,
				Score:    lawScores[law],
				MaxScore: maxS,
			})
		}
	}

	return ScoreResult{
		TotalScore:   totalScore,
		MaxScore:     maxScore,
		Percentage:   percentage,
		Title:        title,
		Description:  description,
		LawBreakdown: breakdown,
	}
}

func main() {
	router := gin.Default()

	// Register custom template functions
	router.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/quiz", func(c *gin.Context) {
		c.HTML(200, "questionnaire.html", gin.H{
			"questions": questions,
		})
	})

	router.POST("/api/score", func(c *gin.Context) {
		var submission QuizSubmission
		if err := c.ShouldBindJSON(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := calculateScore(submission.Answers)
		c.JSON(http.StatusOK, result)
	})

	router.GET("/results", func(c *gin.Context) {
		c.HTML(200, "results.html", nil)
	})

	router.Run() // listens on 0.0.0.0:8080 by default
}
