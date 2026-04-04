package gemini

import (
	. "github.com/Rovemark/LogicaProxy/v6/internal/constant"
	"github.com/Rovemark/LogicaProxy/v6/internal/interfaces"
	"github.com/Rovemark/LogicaProxy/v6/internal/translator/translator"
)

func init() {
	translator.Register(
		Gemini,
		Claude,
		ConvertGeminiRequestToClaude,
		interfaces.TranslateResponse{
			Stream:     ConvertClaudeResponseToGemini,
			NonStream:  ConvertClaudeResponseToGeminiNonStream,
			TokenCount: GeminiTokenCount,
		},
	)
}
