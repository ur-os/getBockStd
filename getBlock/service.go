package getBlock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Service struct {
	getBlock *GetBlock
}

type ResponseTopFive struct {
	Status  string   `json:"status"`
	Code    int      `json:"code"`
	Payload []string `json:"payload"`
}

func NewService() *Service {
	endpoint := os.Getenv("GET_BLOCK_ENDPOINT")
	apiKey := os.Getenv("GET_BLOCK_API_KEY")
	timeout := os.Getenv("GET_BLOCK_TIMEOUT")
	depth := os.Getenv("GET_BLOCK_DEPTH")
	pullingStep := os.Getenv("GET_BLOCK_PULLING_STEP")
	rpsLimit := os.Getenv("GET_BLOCK_RPS_LIMIT")

	return &Service{
		New(endpoint, apiKey, timeout, depth, pullingStep, rpsLimit),
	}
}

const InternalError = "InternalError"

func (s *Service) GetTopFiveUserActivity(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getTopFiveUserActivity" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	start := time.Now()
	topUsers := s.getBlock.GetTop5Users()
	end := time.Now()

	fmt.Printf("\nbench: %s", end.Sub(start).String())

	response := ResponseTopFive{Status: "ok", Code: http.StatusOK}

	for _, user := range topUsers {
		response.Payload = append(response.Payload, fmt.Sprintf("%s = %d", user.Address, user.Count))
	}

	resultByte, err := json.Marshal(response)
	if err != nil {
		errorResponse, _ := json.Marshal(ResponseTopFive{Status: InternalError, Code: http.StatusInternalServerError})
		fmt.Fprintf(w, string(errorResponse))
	}

	fmt.Fprintf(w, string(resultByte))
}
