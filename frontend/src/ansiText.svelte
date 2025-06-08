<script>
    // Import runes from 'svelte' (this is a Svelte 5 feature)
    // You might need to adjust your svelte.config.js to enable runes if not already done.
    // For Svelte 5, these are automatically available if your setup supports it.
    // However, explicitly importing them here is good practice or may be required
    // depending on your build setup.
    // In many Svelte 5 setups, $derived and $state are globally available within .svelte files
    // if runes are enabled, so explicit import might not be strictly necessary,
    // but it won't hurt and provides clarity.
    // For clarity, I'm including `from 'svelte'` but often it's implicitly available.
    // import { $derived } from 'svelte/compiler'; // <-- This is generally incorrect for Svelte 5 runtime.
                                                 // The correct way is just to use `$derived` directly.

    // Props are still declared with `export let`.
    // In Svelte 5, these become "signals" automatically.
    export let ansiLogString = '';

    // Функція для застосування стилів до тексту (не реактивна, звичайна JS функція)
    function applyStyles(text, styles) {
        if (!text) return '';

        // Екранування HTML-спеціальних символів для безпеки
        const escapedText = text
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');

        let styleString = '';
        for (const prop in styles) {
            if (styles.hasOwnProperty(prop)) {
                styleString += `${prop}: ${styles[prop]};`;
            }
        }

        if (styleString) {
            return `<span style="${styleString}">${escapedText}</span>`;
        } else {
            return escapedText;
        }
    }

    // Головна функція для конвертації ANSI в HTML (не реактивна, звичайна JS функція)
    function convertAnsiToHtml(ansiString) {
        let html = '';
        let currentStyle = {};
        const ansiRegex = /\x1b\[([0-9;]*)m/g;
        let lastIndex = 0;

        const foregroundColors = {
            30: 'black', 31: 'red', 32: 'green', 33: 'yellow',
            34: 'blue', 35: 'magenta', 36: 'cyan', 37: 'white',
            90: 'gray', 91: '#FF0000', 92: '#00FF00', 93: '#FFFF00',
            94: '#0000FF', 95: '#FF00FF', 96: '#00FFFF', 97: '#FFFFFF'
        };

        const backgroundColors = {
            40: 'black', 41: 'red', 42: 'green', 43: 'yellow',
            44: 'blue', 45: 'magenta', 46: 'cyan', 47: 'white',
            100: 'gray', 101: '#FF0000', 102: '#00FF00', 103: '#FFFF00',
            104: '#0000FF', 105: '#FF00FF', 106: '#00FFFF', 107: '#FFFFFF'
        };

        let match;
        while ((match = ansiRegex.exec(ansiString)) !== null) {
            if (match.index > lastIndex) {
                const text = ansiString.substring(lastIndex, match.index);
                html += applyStyles(text, currentStyle);
            }

            const codes = match[1].split(';').map(Number);

            codes.forEach(code => {
                if (code === 0) {
                    currentStyle = {};
                } else if (code === 1) {
                    currentStyle.fontWeight = 'bold';
                } else if (code === 3) {
                    currentStyle.fontStyle = 'italic';
                } else if (code === 4) {
                    currentStyle.textDecoration = 'underline';
                } else if (code === 7) { // Інвертування кольорів
                    const tempFg = currentStyle.color;
                    const tempBg = currentStyle.backgroundColor;
                    currentStyle.color = tempBg || 'black';
                    currentStyle.backgroundColor = tempFg || 'white';
                } else if (code >= 30 && code <= 37) {
                    currentStyle.color = foregroundColors[code];
                } else if (code >= 40 && code <= 47) {
                    currentStyle.backgroundColor = backgroundColors[code];
                } else if (code >= 90 && code <= 97) {
                    currentStyle.color = foregroundColors[code];
                } else if (code >= 100 && code <= 107) {
                    currentStyle.backgroundColor = backgroundColors[code];
                }
            });

            lastIndex = ansiRegex.lastIndex;
        }

        if (lastIndex < ansiString.length) {
            html += applyStyles(ansiString.substring(lastIndex), currentStyle);
        }

        return html.replace(/\n/g, '<br/>');
    }

    // $$$$ Svelte 5 Runes change $$$$
    // Використовуємо $derived для реактивних обчислень.
    // $derived автоматично "підписується" на сигнали (в даному випадку, ansiLogString)
    // і переобчислюється, коли ці сигнали змінюються.
    const htmlOutput = $derived(convertAnsiToHtml(ansiLogString));
</script>

<div class="log-viewer">
    {@html htmlOutput}
</div>

<style>
    /* Стилі залишаються незмінними */
    .log-viewer {
        font-family: monospace;
        white-space: pre-wrap;
        background-color: #282c34;
        color: #abb2bf;
        padding: 15px;
        border-radius: 8px;
        overflow-y: auto;
        max-height: 500px;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    }
</style>