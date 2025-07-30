package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-worker/internal/commons"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type DetectionResponse struct {
	Result []struct {
		Box struct {
			Probability float64 `json:"probability"`
			XMin        int     `json:"x_min"`
			XMax        int     `json:"x_max"`
			YMin        int     `json:"y_min"`
			YMax        int     `json:"y_max"`
		} `json:"box"`
		Embedding []float64 `json:"embedding"`
	} `json:"result"`
}

func downloadImageFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %v", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %v", err)
	}

	return imgData, nil
}

func GenerateEmbedding(ctx context.Context, photoPath string) (*[]byte, error) {
	imageData, err := downloadImageFromURL(photoPath)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, commons.NewPermanentError(fmt.Sprintf("image not found: %v", err))
		}
		return nil, fmt.Errorf("download failed: %w", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to copy image data: %v", err)
	}
	writer.Close()

	url := os.Getenv("COMPRE_API_BASE") + "/v1/detection/detect?face_plugins=calculator"
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("x-api-key", os.Getenv("COMPREFACE_DETECTION_KEY"))

	resp, err := DefaultHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		bodyErr, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status: %s, body: %s", resp.Status, string(bodyErr))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	minProbability, err := strconv.ParseFloat(os.Getenv("MIN_PROBABILITY"), 64)
	if err != nil {
		minProbability = 0.9
	}

	minHeight, err := strconv.Atoi(os.Getenv("MIN_HEIGHT"))
	if err != nil {
		minHeight = 100
	}

	minWidth, err := strconv.Atoi(os.Getenv("MIN_WIDTH"))
	if err != nil {
		minWidth = 100
	}

	var result DetectionResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, err
	}

	var finalResult DetectionResponse

	for _, item := range result.Result {
		width := item.Box.XMax - item.Box.XMin
		height := item.Box.YMax - item.Box.YMin
		if item.Box.Probability > minProbability || (width >= minWidth && height >= minHeight) {
			finalResult.Result = append(finalResult.Result, item)
		}
	}
	data, err := json.Marshal(finalResult.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal final result: %v", err)
	}

	return &data, nil

}

// func GenerateEmbedding(imageUrl string) *[]byte {
// 	url := os.Getenv("COMPRE_API_BASE")
// 	verifyKey := os.Getenv("COMPREFACE_DETECTION_KEY")

// 	imageData, err := downloadImageFromURL(imageUrl)
// 	if err != nil {
// 		fmt.Printf("Error downloading image: %v\n", err)
// 	}

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	part, err := writer.CreateFormFile("file", "image.jpg")
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("failed to create form file: %v", err))
// 	}

// 	_, err = io.Copy(part, bytes.NewReader(imageData))
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("failed to copy image data: %v", err))
// 	}

// 	err = writer.Close()
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("failed to close writer: %v", err))
// 	}

// 	request, err := http.NewRequest("POST", url+"/v1/detection/detect?face_plugins=calculator", body)
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("failed to create new request: %v", err))
// 	}

// 	request.Header.Set("x-api-key", verifyKey)
// 	request.Header.Set("Content-Type", writer.FormDataContentType())

// 	client := &http.Client{}
// 	response, err := client.Do(request)

// 	if err != nil {
// 		fmt.Println("Error request")
// 		return nil
// 	}
// 	defer response.Body.Close()

// 	responseBody, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		fmt.Println("Cannot read the response body")
// 	}

// 	minProbability, err := strconv.ParseFloat(os.Getenv("MIN_PROBABILITY"), 64)
// 	if err != nil {
// 		minProbability = 0.9
// 	}

// 	minHeight, err := strconv.Atoi(os.Getenv("MIN_HEIGHT"))
// 	if err != nil {
// 		minHeight = 100
// 	}

// 	minWidth, err := strconv.Atoi(os.Getenv("MIN_WIDTH"))
// 	if err != nil {
// 		minWidth = 100
// 	}

// 	var result DetectionResponse

// 	json.Unmarshal(responseBody, &result)

// 	var finalResult DetectionResponse

// 	for _, item := range result.Result {
// 		if item.Box.Probability > minProbability || (((item.Box.XMax - item.Box.XMin) >= minWidth) && ((item.Box.YMax - item.Box.YMin) >= minHeight)) {
// 			finalResult.Result = append(finalResult.Result, item)
// 		}
// 	}

// 	data, _ := json.Marshal(finalResult.Result)

// 	return &data
// 	// return nil
// }
