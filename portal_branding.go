package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

type PortalBranding struct {
	CssContent string `json:"css_content"`
}

func (p *PortalBranding) GetCssContent() (string, error) {
	if p == nil {
		return "", logError("PortalBranding struct is not properly initialized")
	}
	if p.CssContent == "" {
		return "", logError("CssContent field is empty")
	}
	sanitizedCssContent := p.sanitizeCssContent(p.CssContent)
	if !p.isValidCssContent(sanitizedCssContent) {
		return "", logError("CssContent field is not properly sanitized")
	}
	return sanitizedCssContent, nil
}

func (p *PortalBranding) sanitizeCssContent(cssContent string) string {
	// Remove any script tags
	cssContent = regexp.MustCompile(`<script>.*?</script>`).ReplaceAllString(cssContent, "")
	// Remove any HTML tags
	cssContent = regexp.MustCompile(`<.*?>`).ReplaceAllString(cssContent, "")
	// Remove any JavaScript code
	cssContent = regexp.MustCompile(`javascript:.*?;`).ReplaceAllString(cssContent, "")
	// Remove any URL that starts with javascript:
	cssContent = regexp.MustCompile(`url\(javascript:.*?\)`).ReplaceAllString(cssContent, "")
	// Remove any CSS rules that contain JavaScript code
	cssContent = regexp.MustCompile(`.*?{.*?javascript:.*?;.*?}`).ReplaceAllString(cssContent, "")
	return strings.TrimSpace(cssContent)
}

func (p *PortalBranding) isValidCssContent(cssContent string) bool {
	// Check if the CSS content contains any malicious code
	if regexp.MustCompile(`<script>.*?</script>`).MatchString(cssContent) {
		return false
	}
	if regexp.MustCompile(`<.*?>`).MatchString(cssContent) {
		return false
	}
	if regexp.MustCompile(`javascript:.*?;`).MatchString(cssContent) {
		return false
	}
	if regexp.MustCompile(`url\(javascript:.*?\)`).MatchString(cssContent) {
		return false
	}
	if regexp.MustCompile(`.*?{.*?javascript:.*?;.*?}`).MatchString(cssContent) {
		return false
	}
	return true
}

func logError(message string) error {
	log.Println(message)
	return nil
}

func main() {
	portalBranding := &PortalBranding{
		CssContent: "body { background-color: #f2f2f2; } <script>alert('XSS')</script>",
	}
	cssContent, err := portalBranding.GetCssContent()
	if err != nil {
		log.Println(err)
	}
	log.Println(cssContent)

	// Test case with empty CssContent field
	portalBrandingEmpty := &PortalBranding{
		CssContent: "",
	}
	cssContentEmpty, errEmpty := portalBrandingEmpty.GetCssContent()
	if errEmpty != nil {
		log.Println(errEmpty)
	}
	log.Println(cssContentEmpty)

	// Test case with nil PortalBranding struct
	var portalBrandingNil *PortalBranding
	cssContentNil, errNil := portalBrandingNil.GetCssContent()
	if errNil != nil {
		log.Println(errNil)
	}
	log.Println(cssContentNil)
}