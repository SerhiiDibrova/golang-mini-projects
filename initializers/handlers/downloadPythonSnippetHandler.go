package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"text/template"

	"initializers/services"
	"initializers/utils"
)

func downloadPythonSnippetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	downloadToken, err := utils.GenerateDownloadToken()
	if err != nil {
		http.Error(w, "Failed to generate download token", http.StatusInternalServerError)
		return
	}

	err = _addDownloadTokenCookie(w, downloadToken)
	if err != nil {
		http.Error(w, "Failed to set download token cookie", http.StatusInternalServerError)
		return
	}

	pythonSnippet, err := services.DownloadPythonSnippetService()
	if err != nil {
		http.Error(w, "Failed to download Python snippet", http.StatusInternalServerError)
		return
	}

	if pythonSnippet == nil {
		http.Error(w, "Python snippet is empty", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("python_snippet.tmpl")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pythonSnippet)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func _addDownloadTokenCookie(w http.ResponseWriter, downloadToken string) error {
	cookie := &http.Cookie{
		Name:  "download_token",
		Value: downloadToken,
	}
	http.SetCookie(w, cookie)
	return nil
}