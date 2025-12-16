package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestTTRSSHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		title          string
		content        string
		wantTitle      string
		wantContent    string
		wantStatusCode int
		wantError      string
	}{
		{
			name:           "s2t success",
			path:           "/s2t",
			title:          "简体中文",
			content:        "这是一个测试",
			wantTitle:      "簡體中文",
			wantContent:    "這是一個測試",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "t2s success",
			path:           "/t2s",
			title:          "簡體中文",
			content:        "這是一個測試",
			wantTitle:      "简体中文",
			wantContent:    "这是一个测试",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "empty path defaults to t2s",
			path:           "/",
			title:          "簡體中文",
			content:        "這是一個測試",
			wantTitle:      "简体中文",
			wantContent:    "这是一个测试",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "t2s html content",
			path:           "/t2s",
			title:          "HTML測試",
			content:        "<div><h1>標題</h1><p>這是一個測試</p></div>",
			wantTitle:      "HTML测试",
			wantContent:    "<div><h1>标题</h1><p>这是一个测试</p></div>",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "t2s xml content",
			path:           "/t2s",
			title:          "XML測試",
			content:        "<note><to>用戶</to><from>管理員</from><heading>提醒</heading><body>這是一個測試</body></note>",
			wantTitle:      "XML测试",
			wantContent:    "<note><to>用户</to><from>管理员</from><heading>提醒</heading><body>这是一个测试</body></note>",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "t2s markdown content",
			path:           "/t2s",
			title:          "Markdown測試",
			content:        "# 標題\n\n這是一個測試",
			wantTitle:      "Markdown测试",
			wantContent:    "# 标题\n\n这是一个测试",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "invalid scheme",
			path:           "/invalid_scheme",
			title:          "text",
			content:        "content",
			wantStatusCode: http.StatusOK, // The handler currently returns 200 for invalid scheme
			wantError:      "Invalid convert scheme.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("title", tt.title)
			formData.Set("content", tt.content)

			req, err := http.NewRequest("POST", tt.path, strings.NewReader(formData.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(handler)
			handlerFunc.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatusCode)
			}

			if tt.wantError != "" {
				if rr.Body.String() != tt.wantError {
					t.Errorf("handler returned wrong error message: got %q want %q",
						rr.Body.String(), tt.wantError)
				}
				return
			}

			var ret Ret
			if err := json.Unmarshal(rr.Body.Bytes(), &ret); err != nil {
				t.Errorf("failed to unmarshal response: %v, body: %s", err, rr.Body.String())
				return
			}

			if !strings.Contains(ret.Title, tt.wantTitle) {
				t.Errorf("Title not converted correctly, got: %s, want part: %s", ret.Title, tt.wantTitle)
			}
			if !strings.Contains(ret.Content, tt.wantContent) {
				t.Errorf("Content not converted correctly, got: %s, want part: %s", ret.Content, tt.wantContent)
			}
		})
	}
}

func TestCommonJsonApiHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		wantStatusCode int
		wantSubstring  string
	}{
		{
			name:           "s2t success",
			method:         "POST",
			path:           "/api/s2t",
			body:           `{"text": "简体中文"}`,
			wantStatusCode: http.StatusOK,
			wantSubstring:  "簡體中文",
		},
		{
			name:           "t2s success",
			method:         "POST",
			path:           "/api/t2s",
			body:           `{"text": "簡體中文"}`,
			wantStatusCode: http.StatusOK,
			wantSubstring:  "简体中文",
		},
		{
			name:           "invalid method",
			method:         "GET",
			path:           "/api/s2t",
			body:           "",
			wantStatusCode: http.StatusMethodNotAllowed,
			wantSubstring:  "Method not allowed",
		},
		{
			name:           "invalid scheme",
			method:         "POST",
			path:           "/api/invalid_scheme",
			body:           `{"text": "abc"}`,
			wantStatusCode: http.StatusBadRequest,
			wantSubstring:  "Invalid convert scheme",
		},
		{
			name:           "invalid json",
			method:         "POST",
			path:           "/api/s2t",
			body:           `{invalid json}`,
			wantStatusCode: http.StatusBadRequest,
			wantSubstring:  "Invalid JSON payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(apiHandler)
			handlerFunc.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatusCode)
			}

			body := rr.Body.String()
			if tt.wantStatusCode == http.StatusOK {
				var resp map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}
				if got, ok := resp["converted"]; !ok {
					t.Errorf("response missing 'converted' field")
				} else {
					if !strings.Contains(got, tt.wantSubstring) {
						t.Errorf("Text not converted correctly, got: %s, want part: %s", got, tt.wantSubstring)
					}
				}
			} else {
				// Error cases often return plain text or simple body
				if !strings.Contains(strings.TrimSpace(body), tt.wantSubstring) {
					t.Errorf("Error message mismatch: got %q, want part %q", body, tt.wantSubstring)
				}
			}
		})
	}
}
