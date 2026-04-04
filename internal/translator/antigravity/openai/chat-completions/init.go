package chat_completions

import (
	. "github.com/Rovemark/LogicaProxy/v6/internal/constant"
	"github.com/Rovemark/LogicaProxy/v6/internal/interfaces"
	"github.com/Rovemark/LogicaProxy/v6/internal/translator/translator"
)

func init() {
	translator.Register(
		OpenAI,
		Antigravity,
		ConvertOpenAIRequestToAntigravity,
		interfaces.TranslateResponse{
			Stream:    ConvertAntigravityResponseToOpenAI,
			NonStream: ConvertAntigravityResponseToOpenAINonStream,
		},
	)
}
