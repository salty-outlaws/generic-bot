package util

import "encoding/json"

func mapToString(v map[string]any) string {
	s, _ := json.Marshal(v)
	return string(s)
}

func Embed(title, summary string, field map[string]string) string {

	return mapToString(map[string]any{
		"type":    "embed",
		"title":   title,
		"summary": summary,
		"fields":  field,
	})
}

func Image(link string) string {
	return mapToString(map[string]any{
		"type":  "image",
		"image": link,
	})
}

func Text(message string) string {
	return mapToString(map[string]any{
		"type":    "text",
		"message": message,
	})
}
