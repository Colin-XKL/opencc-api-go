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
