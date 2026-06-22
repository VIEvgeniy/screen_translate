# screen_translate

Экранный переводчик для Linux (GNOME/Wayland). Выделяешь область экрана мышью — получаешь перевод во всплывающем окне.

## Как работает

1. Нажимаешь хоткей (настраивается через GNOME)
2. Курсор превращается в прицел — выделяешь область с текстом
3. `gnome-screenshot` делает скриншот области
4. `tesseract` распознаёт текст (OCR, локально)
5. `translate-shell` переводит текст через Google Translate
6. Результат появляется в окне `zenity`

## Требования

- Linux с GNOME на Wayland
- Go 1.18+

## Установка зависимостей

### ALT Linux
```bash
bash install.sh
```

### Debian / Ubuntu
```bash
sudo apt install tesseract-ocr tesseract-ocr-eng tesseract-ocr-rus \
    gnome-screenshot zenity translate-shell
```

### Fedora
```bash
sudo dnf install tesseract tesseract-langpack-eng tesseract-langpack-rus \
    gnome-screenshot zenity translate-shell
```

## Сборка

```bash
git clone https://github.com/ВАШ_НИК/screen_translate
cd screen_translate
go build -o screen_translate .
```

## Настройка хоткея в GNOME

Добавляем хоткей **Alt+T**:

```bash
# Читаем текущий список хоткеев
gsettings get org.gnome.settings-daemon.plugins.media-keys custom-keybindings

# Добавляем новый слот (не удаляя существующие — вставь свой список)
gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings \
  "['/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/translate-keybind/']"

# Настраиваем хоткей
BASE="org.gnome.settings-daemon.plugins.media-keys.custom-keybinding"
PATH="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/translate-keybind/"

gsettings set ${BASE}:${PATH} name    'Экранный переводчик'
gsettings set ${BASE}:${PATH} command '/ПОЛНЫЙ/ПУТЬ/ДО/screen_translate'
gsettings set ${BASE}:${PATH} binding '<Alt>t'
```

Замени `/ПОЛНЫЙ/ПУТЬ/ДО/screen_translate` на вывод команды:
```bash
realpath ./screen_translate
```

## Использование

```bash
# Перевод + оригинал
./screen_translate

# Только перевод
./screen_translate -short
```

Флаг `-short` удобно передавать через команду хоткея:
```bash
gsettings set ${BASE}:${PATH} command '/путь/screen_translate -short'
```

## Ограничения

- OCR работает лучше всего с чётким экранным текстом (интерфейсы, субтитры, документы)
- Требует подключения к интернету для перевода (OCR работает офлайн)
