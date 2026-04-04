package responses

import (
	. "github.com/Rovemark/LogicaProxy/v6/internal/constant"
	"github.com/Rovemark/LogicaProxy/v6/internal/interfaces"
	"github.com/Rovemark/LogicaProxy/v6/internal/translator/translator"
)

func init() {
	translator.Register(
		OpenaiResponse,
		OpenAI,
		ConvertOpenAIResponsesRequestToOpenAIChatCompletions,
		interfaces.TranslateResponse{
			Stream:    ConvertOpenAIChatCompletionsResponseToOpenAIResponses,
			NonStream: ConvertOpenAIChatCompletionsResponseToOpenAIResponsesNonStream,
		},
	)
}
