package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type DownloadPythonSnippetService struct {
	pythonSnippetTemplate *template.Template
}

func NewDownloadPythonSnippetService(pythonSnippetTemplate *template.Template) *DownloadPythonSnippetService {
	return &DownloadPythonSnippetService{pythonSnippetTemplate: pythonSnippetTemplate}
}

func (s *DownloadPythonSnippetService) downloadPythonSnippetService(w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	if w == nil {
		return nil, errors.New("response writer is nil")
	}
	if r == nil {
		return nil, errors.New("request is nil")
	}
	if r.Method != http.MethodGet {
		return nil, errors.New("invalid request method")
	}
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{}}
	err := s.addDownloadTokenCookie(w)
	if err != nil {
		return nil, err
	}
	s.setContentTypeAndDisposition(w)
	err = s.renderPythonSnippetTemplate(w)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *DownloadPythonSnippetService) addDownloadTokenCookie(w http.ResponseWriter) error {
	downloadToken := "download-token"
	downloadTokenValue := "token-value"
	expirationTime := time.Now().Add(1 * time.Hour)
	cookie := &http.Cookie{
		Name:    downloadToken,
		Value:   downloadTokenValue,
		Expires: expirationTime,
	}
	err := http.SetCookie(w, cookie)
	if err != nil {
		return err
	}
	return nil
}

func (s *DownloadPythonSnippetService) setContentTypeAndDisposition(w http.ResponseWriter) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	if w.Header().Get("Content-Disposition") == "" {
		w.Header().Set("Content-Disposition", "attachment; filename=python_snippet.py")
	}
}

func (s *DownloadPythonSnippetService) renderPythonSnippetTemplate(w http.ResponseWriter) error {
	if s.pythonSnippetTemplate == nil {
		return errors.New("python snippet template is nil")
	}
	err := s.pythonSnippetTemplate.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}