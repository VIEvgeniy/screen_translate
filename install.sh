#!/bin/bash
# Установка зависимостей для screen_translate на ALT Linux
set -e

echo "=== Установка системных пакетов ==="
sudo apt-get install -y \
    tesseract \
    tesseract-langpack-en \
    tesseract-langpack-ru \
    gnome-screenshot \
    zenity \
    golang

echo ""
echo "=== Сборка ==="
go build -o screen_translate .

echo ""
echo "=== Настройка хоткея Alt+T в GNOME ==="
BINARY_PATH="$(realpath ./screen_translate)"
BASE="org.gnome.settings-daemon.plugins.media-keys.custom-keybinding"
KPATH="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/translate-keybind/"

# Читаем существующие хоткеи и добавляем наш слот
CURRENT=$(gsettings get org.gnome.settings-daemon.plugins.media-keys custom-keybindings)
if echo "$CURRENT" | grep -q "translate-keybind"; then
    echo "Слот translate-keybind уже существует, обновляем команду..."
else
    # Добавляем новый слот к существующему списку
    NEW=$(echo "$CURRENT" | sed "s|]|, '${KPATH}']|" | sed "s|\[@\]|['${KPATH}']|")
    gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "$NEW"
fi

gsettings set ${BASE}:${KPATH} name    'Экранный переводчик'
gsettings set ${BASE}:${KPATH} command "${BINARY_PATH} -short"
gsettings set ${BASE}:${KPATH} binding '<Alt>t'

echo ""
echo "=== Готово! ==="
echo "Хоткей Alt+T зарегистрирован."
echo "Бинарник: ${BINARY_PATH}"
echo ""
echo "Для теста запусти вручную:"
echo "  ${BINARY_PATH}"
