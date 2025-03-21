package ollama

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
)

func getOllamaRest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

type ollamaTags struct {
	Models []ollamaTag `json:"models"`
}

type ollamaTag struct {
	Model string `json:"model"`
	Name  string `json:"name"`
}

func OllamaConnectionFromContext(c *cli.Context) error {
	// get the connection string
	connectionString := c.String("ollama-url")
	if connectionString == "" {
		return fmt.Errorf("ollama-url is required")
	}

	// make sure the response body is 'Ollama is running'
	body, err := getOllamaRest(connectionString)
	if err != nil {
		return fmt.Errorf("failed to get the response from %s: %w", connectionString, err)
	}
	if string(body) != "Ollama is running" {
		return fmt.Errorf("something was at %s, but it seems like it was not Ollama: %v", connectionString, body)
	}

	// check that a model of type nomic-embed-text is available
	body, err = getOllamaRest(connectionString + "/api/tags")
	if err != nil {
		return fmt.Errorf("failed to get tags from %s: %w", connectionString, err)
	}

	var tags ollamaTags
	if err := json.Unmarshal(body, &tags); err != nil {
		return fmt.Errorf("failed to read the existing models from ollama: %w", err)
	}

	// for now we hardcode nomic-embed-text
	var foundTag ollamaTag
	for _, tag := range tags.Models {
		if strings.HasPrefix(tag.Model, "nomic-embed-text") {
			foundTag = tag
			break
		}
	}
	if foundTag == (ollamaTag{}) {
		return fmt.Errorf("datailama needs the nomic-embed-text model, which was not found Run \n ollama pull nomic-embed-text:latest\n to get it")
	} else {
		fmt.Printf("found %v\n", foundTag.Name)
	}
	return nil
}
