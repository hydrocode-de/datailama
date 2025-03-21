package ollama

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pgvector/pgvector-go"
	"github.com/urfave/cli/v2"
)

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func EmbedText(c *cli.Context, text string) (pgvector.Vector, error) {
	embed_model, err := CheckOllamaConnection(c)
	if err != nil {
		return pgvector.Vector{}, err
	}

	payload := map[string]string{
		"model":  embed_model,
		"prompt": text,
	}

	bytePayload, err := json.Marshal(payload)
	if err != nil {
		return pgvector.Vector{}, err
	}

	ollama_url := c.String("ollama-url")

	response, err := http.Post(ollama_url+"/api/embeddings", "application/json", bytes.NewBuffer(bytePayload))
	if err != nil {
		return pgvector.Vector{}, err
	}
	defer response.Body.Close()

	var embeddingResponse EmbeddingResponse
	err = json.NewDecoder(response.Body).Decode(&embeddingResponse)
	if err != nil {
		return pgvector.Vector{}, err
	}

	vector := pgvector.NewVector(embeddingResponse.Embedding)
	return vector, nil
}
