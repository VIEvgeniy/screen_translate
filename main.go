package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	langFrom = "en"
	langTo   = "ru"
)

func translateText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	params := url.Values{}
	params.Set("q", text)
	params.Set("langpair", langFrom+"|"+langTo)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.mymemory.translated.net/get?" + params.Encode())
	if err != nil {
		return fmt.Sprintf("[Ошибка сети: %v]", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		ResponseStatus json.RawMessage `json:"responseStatus"`
		ResponseData   struct {
			TranslatedText string `json:"translatedText"`
		} `json:"responseData"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Sprintf("[Ошибка парсинга: %v]", err)
	}
	// responseStatus может быть числом 200 или строкой "200"
	status := strings.Trim(string(data.ResponseStatus), `"`)
	if status == "200" {
		return data.ResponseData.TranslatedText
	}
	return fmt.Sprintf("[Ошибка перевода: %s]", status)
}

func takeScreenshot() (string, error) {
	tmp, err := os.CreateTemp("", "scrtrans_*.png")
	if err != nil {
		return "", err
	}
	tmp.Close()
	path := tmp.Name()

	// gnome-screenshot -a: интерактивное выделение области
	cmd := exec.Command("gnome-screenshot", "-a", "-f", path)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("gnome-screenshot: %v", err)
	}
	// При отмене (ESC) gnome-screenshot создаёт пустой файл
	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 {
		os.Remove(path)
		return "", nil
	}
	return path, nil
}
func ocrImage(imgPath string) (string, error) {
	cmd := exec.Command("tesseract", imgPath, "stdout", "-l", "eng+rus", "--psm", "6")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("tesseract: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func showResult(original, translated string, short bool) {
	if translated == "" {
		translated = "(текст не распознан)"
	}
	var text string
	if short {
		text = translated
	} else {
		text = fmt.Sprintf("ПЕРЕВОД:\n%s\n\nОРИГИНАЛ:\n%s", translated, original)
	}
	exec.Command("zenity",
		"--info",
		"--title=Экранный переводчик",
		"--text="+text,
		"--width=460",
		"--height=280",
		"--no-markup",
	).Run()
}

func showError(msg string) {
	exec.Command("zenity", "--error", "--title=Переводчик", "--text="+msg, "--width=300").Run()
}

func main() {
	short := len(os.Args) > 1 && os.Args[1] == "-short"

	imgPath, err := takeScreenshot()
	if err != nil {
		showError(fmt.Sprintf("Ошибка скриншота:\n%v", err))
		os.Exit(1)
	}
	if imgPath == "" {
		os.Exit(0) // пользователь отменил
	}
	defer os.Remove(imgPath)

	original, err := ocrImage(imgPath)
	if err != nil {
		showError(fmt.Sprintf("Ошибка OCR:\n%v", err))
		os.Exit(1)
	}

	translated := translateText(original)
	showResult(original, translated, short)
}
