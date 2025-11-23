package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	// Test case: s2t
	formData := url.Values{}
	formData.Set("title", "简体中文")
	formData.Set("content", "这是一个测试")

	req, err := http.NewRequest("POST", "/s2t", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var ret Ret
	if err := json.Unmarshal(rr.Body.Bytes(), &ret); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	t.Logf("Title: %s", ret.Title)
	t.Logf("Content: %s", ret.Content)

	// "简" -> "簡"
	if !strings.Contains(ret.Title, "簡") {
		t.Errorf("Title not converted correctly, got: %s", ret.Title)
	}
	// "这" -> "這", "测试" -> "測試"
	if !strings.Contains(ret.Content, "測試") {
		t.Errorf("Content not converted correctly, got: %s", ret.Content)
	}
}

func TestApiHandler(t *testing.T) {
	// Test case: s2t via JSON API
	reqBody := `{"text": "简体中文"}`
	req, err := http.NewRequest("POST", "/api/s2t", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	// We need to register the handler same as main, but since main registers it globally and we are in test,
	// we might need to simulate the router or just test the specific handler if I export it.
	// For now, let's assume I'll name the new handler `apiHandler`.
	// However, `http.HandleFunc` in main is not easily accessible here unless I expose the mux or register it.
	// The existing test creates a HandlerFunc from `handler`.
	// I'll assume I'll create a function `apiHandler` and test it directly here.

	handler := http.HandlerFunc(apiHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if got, ok := resp["converted"]; !ok {
		t.Errorf("response missing 'converted' field")
	} else {
		t.Logf("Converted: %s", got)
		if !strings.Contains(got, "簡體中文") {
			t.Errorf("Text not converted correctly, got: %s", got)
		}
	}
}
