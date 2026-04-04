package middleware

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// AIRouterConfig defines model routing rules per task category.
type AIRouterConfig struct {
	Enabled      bool              `yaml:"enabled" json:"enabled"`
	DefaultModel string            `yaml:"default_model" json:"default_model"`
	Routes       map[string]string `yaml:"routes" json:"routes"`
}

// DefaultAIRouterConfig returns sensible defaults for AI routing.
func DefaultAIRouterConfig() AIRouterConfig {
	return AIRouterConfig{
		Enabled:      false,
		DefaultModel: "",
		Routes: map[string]string{
			"code":      "claude-sonnet-4-6",
			"reasoning": "claude-opus-4-6",
			"fast":      "claude-haiku-4-5-20251001",
			"creative":  "claude-sonnet-4-6",
			"analysis":  "claude-opus-4-6",
		},
	}
}

type taskCategory struct {
	name     string
	keywords []string
}

var categories = []taskCategory{
	{
		name: "code",
		keywords: []string{
			"write code", "implement", "function", "class", "debug", "fix bug",
			"refactor", "test", "compile", "build", "deploy", "api",
			"endpoint", "database", "query", "sql", "html", "css",
			"javascript", "python", "golang", "rust", "typescript",
			"react", "vue", "angular", "node", "docker", "kubernetes",
			"git", "merge", "commit", "pull request", "code review",
			"def ", "func ", "import ", "package ", "class ", "struct ",
		},
	},
	{
		name: "reasoning",
		keywords: []string{
			"analyze", "explain why", "compare", "evaluate", "assess",
			"strategy", "architecture", "design system", "trade-off",
			"pros and cons", "decision", "complex", "multi-step",
			"reasoning", "think through", "plan", "roadmap",
		},
	},
	{
		name: "fast",
		keywords: []string{
			"translate", "summarize", "list", "format", "convert",
			"extract", "simple", "quick", "short", "brief",
			"yes or no", "true or false", "one word",
		},
	},
	{
		name: "creative",
		keywords: []string{
			"write", "story", "poem", "creative", "imagine",
			"blog post", "article", "essay", "copy", "slogan",
			"marketing", "content", "narrative", "description",
		},
	},
	{
		name: "analysis",
		keywords: []string{
			"data", "chart", "graph", "statistics", "report",
			"financial", "metrics", "benchmark", "performance",
			"trend", "forecast", "predict", "insight",
		},
	},
}

// classifyTask determines the task category from the prompt content.
func classifyTask(prompt string) string {
	lower := strings.ToLower(prompt)
	scores := make(map[string]int)

	for _, cat := range categories {
		for _, kw := range cat.keywords {
			if strings.Contains(lower, kw) {
				scores[cat.name]++
			}
		}
	}

	best := ""
	bestScore := 0
	for name, score := range scores {
		if score > bestScore {
			best = name
			bestScore = score
		}
	}

	if bestScore == 0 {
		return ""
	}
	return best
}

// extractPrompt pulls the user message from an OpenAI or Claude request body.
func extractPrompt(body []byte) string {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return ""
	}

	// OpenAI format: messages[].content
	if msgs, ok := data["messages"].([]interface{}); ok {
		var parts []string
		for _, m := range msgs {
			if msg, ok := m.(map[string]interface{}); ok {
				if role, _ := msg["role"].(string); role == "user" {
					if content, ok := msg["content"].(string); ok {
						parts = append(parts, content)
					}
				}
			}
		}
		return strings.Join(parts, " ")
	}

	// Single prompt field
	if prompt, ok := data["prompt"].(string); ok {
		return prompt
	}

	return ""
}

// AIRouterMiddleware returns middleware that sets X-AI-Route and X-AI-Category headers.
// When config.Enabled is true and no model is explicitly set (or model is "auto"),
// it rewrites the model field in the request body.
func AIRouterMiddleware(cfg AIRouterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Enabled {
			c.Next()
			return
		}

		// Only process POST requests with JSON body
		if c.Request.Method != "POST" || c.Request.Body == nil {
			c.Next()
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

		prompt := extractPrompt(bodyBytes)
		if prompt == "" {
			c.Next()
			return
		}

		category := classifyTask(prompt)
		if category == "" {
			c.Set("ai_category", "general")
			c.Header("X-AI-Category", "general")
			c.Next()
			return
		}

		c.Set("ai_category", category)
		c.Header("X-AI-Category", category)

		// Check if we should route to a specific model
		if routeModel, ok := cfg.Routes[category]; ok && routeModel != "" {
			c.Header("X-AI-Route", routeModel)
			c.Set("ai_routed_model", routeModel)

			// Rewrite model in body if current model is "auto" or empty
			var data map[string]interface{}
			if json.Unmarshal(bodyBytes, &data) == nil {
				currentModel, _ := data["model"].(string)
				if currentModel == "" || currentModel == "auto" {
					data["model"] = routeModel
					if newBody, err := json.Marshal(data); err == nil {
						c.Request.Body = io.NopCloser(strings.NewReader(string(newBody)))
						c.Request.ContentLength = int64(len(newBody))
					}
				}
			}
		}

		c.Next()
	}
}
