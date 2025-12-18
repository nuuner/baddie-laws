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
		Text: "You're presenting and realize you missed a slide. What feels most natural?",
		Law:  "L1: Confidence",
		Options: []Option{
			{Text: "Acknowledge it and offer to send the missing info later", Score: 2},
			{Text: "Adapt on the fly - the audience won't know the difference", Score: 4},
			{Text: "Take a moment to apologize and explain what happened", Score: 1},
			{Text: "Keep your composure and weave the content into what's next", Score: 3},
		},
	},
	{
		ID:   2,
		Text: "Someone challenges your decision publicly. How do you typically respond?",
		Law:  "L1: Confidence",
		Options: []Option{
			{Text: "Hold your position - you've thought this through", Score: 4},
			{Text: "Walk them through your reasoning to get them on board", Score: 3},
			{Text: "Push back firmly so they understand you're serious", Score: 2},
			{Text: "Consider their perspective - they might have a point", Score: 1},
		},
	},
	{
		ID:   3,
		Text: "You put real effort into tonight's look. Someone compliments you. You say:",
		Law:  "L2: Effortless",
		Options: []Option{
			{Text: "Thanks! I actually found this piece at the most amazing place...", Score: 2},
			{Text: "You're sweet - just threw something together", Score: 3},
			{Text: "That's kind of you", Score: 4},
			{Text: "Oh this? Let me tell you about the hunt for these shoes", Score: 1},
		},
	},
	{
		ID:   4,
		Text: "You landed something big after months of work. Time to share it?",
		Law:  "L2: Effortless",
		Options: []Option{
			{Text: "A tasteful post - let the achievement speak for itself", Score: 3},
			{Text: "Share the story - the journey is what makes it meaningful", Score: 1},
			{Text: "Maybe mention it if it comes up naturally", Score: 4},
			{Text: "Celebrate the win with a caption about what it took", Score: 2},
		},
	},
	{
		ID:   5,
		Text: "You're deep in conversation about something you love. You notice others glazing over. You:",
		Law:  "L3: Balance",
		Options: []Option{
			{Text: "Wrap it up smoothly and ask about their interests", Score: 3},
			{Text: "Keep going - genuine enthusiasm is refreshing", Score: 4},
			{Text: "Laugh it off and admit you got carried away", Score: 2},
			{Text: "Pivot to something more universally relatable", Score: 1},
		},
	},
	{
		ID:   6,
		Text: "Someone at a social event asks what you're into. You:",
		Law:  "L3: Balance",
		Options: []Option{
			{Text: "Share what you're genuinely excited about lately", Score: 3},
			{Text: "Mention a few things that usually land well", Score: 2},
			{Text: "Turn it into something they'll want to hear more about", Score: 4},
			{Text: "Keep it general - no need to overshare with strangers", Score: 1},
		},
	},
	{
		ID:   7,
		Text: "The conversation hits a lull. What's your instinct?",
		Law:  "L4: Mystery",
		Options: []Option{
			{Text: "Sit with it - sometimes silence is comfortable", Score: 3},
			{Text: "Use it as a natural pause before moving on", Score: 4},
			{Text: "Bridge the gap with a light comment or joke", Score: 2},
			{Text: "Jump in with a new topic to keep things flowing", Score: 1},
		},
	},
	{
		ID:   8,
		Text: "Someone asks what you're up to this weekend. You say:",
		Law:  "L4: Mystery",
		Options: []Option{
			{Text: "A few things here and there - should be good", Score: 4},
			{Text: "Oh, just the usual stuff - nothing too exciting", Score: 2},
			{Text: "Actually, I've got [detailed Saturday and Sunday plans]...", Score: 1},
			{Text: "Some plans with friends, maybe check out that new spot", Score: 3},
		},
	},
	{
		ID:   9,
		Text: "You've outgrown your current vibe. What's your approach?",
		Law:  "L5: Reinvention",
		Options: []Option{
			{Text: "Start fresh - who you were doesn't define who you'll be", Score: 4},
			{Text: "Evolve naturally over time without making it a thing", Score: 3},
			{Text: "Test the waters with small updates first", Score: 2},
			{Text: "Accept that people have certain expectations of you", Score: 1},
		},
	},
	{
		ID:   10,
		Text: "When it comes to how people treat you, which sounds most like you?",
		Law:  "L6: Royalty",
		Options: []Option{
			{Text: "I speak up when something crosses a line", Score: 3},
			{Text: "People generally know where I stand without me saying much", Score: 4},
			{Text: "I try to be flexible - relationships require compromise", Score: 1},
			{Text: "I drop hints when something bothers me", Score: 2},
		},
	},
	{
		ID:   11,
		Text: "Your style used to work, but lately it feels off. You:",
		Law:  "L5: Reinvention",
		Options: []Option{
			{Text: "Experiment with something completely different", Score: 4},
			{Text: "Tweak a few things and see how it feels", Score: 2},
			{Text: "Stick with what's familiar - it's worked before", Score: 1},
			{Text: "Refresh gradually as inspiration strikes", Score: 3},
		},
	},
	{
		ID:   12,
		Text: "People remember you a certain way. Does that shape your choices?",
		Law:  "L5: Reinvention",
		Options: []Option{
			{Text: "Somewhat - consistency builds trust", Score: 2},
			{Text: "Not really - I do what feels right now", Score: 4},
			{Text: "I try to meet their expectations mostly", Score: 1},
			{Text: "I balance who I was with who I'm becoming", Score: 3},
		},
	},
	{
		ID:   13,
		Text: "Someone keeps canceling plans with you last minute. You:",
		Law:  "L6: Royalty",
		Options: []Option{
			{Text: "Stop initiating - they'll reach out if they want to", Score: 4},
			{Text: "Give them another chance, things happen", Score: 1},
			{Text: "Mention it casually next time you talk", Score: 2},
			{Text: "Have a direct conversation about it", Score: 3},
		},
	},
	{
		ID:   14,
		Text: "You're at a gathering and someone talks over you repeatedly. You:",
		Law:  "L6: Royalty",
		Options: []Option{
			{Text: "Let it go - not worth the awkwardness", Score: 1},
			{Text: "Finish your thought firmly when there's a gap", Score: 3},
			{Text: "Continue as if they hadn't interrupted", Score: 4},
			{Text: "Make a lighthearted comment about being interrupted", Score: 2},
		},
	},
}

func getScoreTitle(percentage float64) (string, string) {
	switch {
	case percentage >= 96:
		return "Supreme Baddie", "You've mastered all the laws. Your presence commands attention and respect effortlessly. You embody confidence, mystery, and royalty in everything you do."
	case percentage >= 85:
		return "Certified Baddie", "You're well on your way to baddie mastery. You understand the power of confidence and mystery, with just a few areas to refine."
	case percentage >= 72:
		return "Rising Baddie", "You have strong baddie instincts but sometimes let doubt creep in. Keep trusting yourself and saying less - your glow up is in progress."
	case percentage >= 58:
		return "Baddie in Training", "You're learning the laws but haven't fully internalized them yet. Focus on building unshakeable confidence and embracing your authentic self."
	case percentage >= 42:
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

	router.GET("/roadmap", func(c *gin.Context) {
		c.HTML(200, "roadmap.html", nil)
	})

	router.Run() // listens on 0.0.0.0:8080 by default
}
