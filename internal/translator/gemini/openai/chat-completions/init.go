package chat_completions

import (
	. "github.com/Rovemark/LogicaProxy/v6/internal/constant"
	"github.com/Rovemark/LogicaProxy/v6/internal/interfaces"
	"github.com/Rovemark/LogicaProxy/v6/internal/translator/translator"
)

func init() {
	translator.Register(
		OpenAI,
		Gemini,
		ConvertOpenAIRequestToGemini,
		interfaces.TranslateResponse{
			Stream:    ConvertGeminiResponseToOpenAI,
			NonStream: ConvertGeminiResponseToOpenAINonStream,
		},
	)
}
