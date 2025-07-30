package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-worker/internal/entity"
	"io"
	"log"
	"net/http"
	"os"
)

type ResponseSimiliarity struct {
	IsSimiliar  bool
	Similiarity float64
}

type VerifyRequest struct {
	Source  []float64   `json:"source"`
	Targets [][]float64 `json:"targets"`
}

type VerifyResponse struct {
	Result []struct {
		Embedding  []float64 `json:"embedding"`
		Similarity float64   `json:"similarity"`
	} `json:"result"`
}

// verify to compreface service
func Verify(ctx context.Context, source entity.Embed, targets []entity.Embed, threshold float64) ResponseSimiliarity {
	url := os.Getenv("COMPRE_API_BASE")
	verifyKey := os.Getenv("COMPREFACE_VERIFICATION_KEY")

	resSimiliarity := ResponseSimiliarity{
		IsSimiliar:  false,
		Similiarity: 0.0,
	}

	var newTargets [][]float64

	for _, target := range targets {
		newTargets = append(newTargets, target.Embedding)
	}

	body := VerifyRequest{
		Source:  source.Embedding,
		Targets: newTargets,
	}

	bodyTransform, err := json.Marshal(body)

	if err != nil {
		fmt.Println("Cannot transform the body")
		return resSimiliarity
	}

	request, err := http.NewRequestWithContext(ctx, "POST", url+"/v1/verification/embeddings/verify", bytes.NewBuffer(bodyTransform))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return resSimiliarity
	}
	request.Header.Set("x-api-key", verifyKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := DefaultHTTPClient.Do(request)

	if err != nil {
		return resSimiliarity
	}
	defer func() {
		io.Copy(io.Discard, response.Body)
		response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return resSimiliarity
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Cannot read the response body")
		return resSimiliarity
	}

	var result VerifyResponse

	if err := json.Unmarshal(responseBody, &result); err != nil {
		log.Printf("Failed to unmarshal verify response: %v", err)
		return resSimiliarity
	}

	isOk := false
	similiarityThreshold := 0.0

	for _, item := range result.Result {
		if item.Similarity > threshold {
			isOk = true
			similiarityThreshold = item.Similarity
			break
		}
	}

	resSimiliarity.IsSimiliar = isOk
	resSimiliarity.Similiarity = similiarityThreshold

	return resSimiliarity
}
