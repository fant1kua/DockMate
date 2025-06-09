<script lang="ts">
    import { Terminal } from 'xterm';
    import { FitAddon } from 'xterm-addon-fit';
    import { WebLinksAddon } from 'xterm-addon-web-links';
    import { WebglAddon } from 'xterm-addon-webgl';
    import type { app } from "../../wailsjs/go/models";
    import { StreamContainerLogs, StopContainerLogs, ExecContainer } from "../../wailsjs/go/app/App";
    import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";

    let { container, onClose } = $props<{
        container: app.ContainerInfo | null;
		onClose(): void;
    }>();

    let terminalElement = $state<HTMLDivElement>();
    let terminal: Terminal;

    $effect(() => {
        if (!terminalElement) return;

        terminal = new Terminal({
            cursorBlink: true,
            cols: 80,
            rows: 14,
            theme: {
                background: '#1e1e1e',
                foreground: '#ffffff'
            },
        });

        // Initialize addons
        const fitAddon = new FitAddon();
        const webLinksAddon = new WebLinksAddon();
        const webglAddon = new WebglAddon();

        // Add addons to terminal
        terminal.loadAddon(fitAddon);
        terminal.loadAddon(webLinksAddon);
        terminal.loadAddon(webglAddon);

        // Open terminal in the container
        terminal.open(terminalElement);
        fitAddon.fit();

        terminal.open(terminalElement);

        if (container) {
            StreamContainerLogs(container.id);
            EventsOn("logStream", (line: string) => {
                terminal.writeln(line);
            });

            let currentCommand = ''

            terminal.onData((data) => {
                if (data === '\r') {
                    // Enter key pressed
                    terminal.write('\r\n');
                    if (currentCommand.trim()) {
                        console.log(currentCommand)
                        ExecContainer(container.id, currentCommand).catch((error) => {
                            terminal.writeln(`Error executing command: ${error}`);
                        });
                    }
                    currentCommand = '';
                } else if (data === '\u0003') {
                    // Ctrl+C
                    terminal.write('^C\r\n');
                    currentCommand = '';
                } else if (data === '\u007F') {
                    // Backspace
                    if (currentCommand.length > 0) {
                        currentCommand = currentCommand.slice(0, -1);
                        terminal.write('\b \b');
                    }
                } else {
                    // Regular character
                    currentCommand += data;
                    terminal.write(data);
                }
            });
        }

        return () => {
            EventsOff('logStream');
            StopContainerLogs();
            terminal.dispose();
        }
    });

</script>

{#if container}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
        <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded-lg w-3/4 h-3/4 flex flex-col">
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-xl font-bold">Terminal for {container.names[0]}</h2>
                <button 
                    class="text-gray-500 hover:text-gray-700"
                    onclick={onClose}
                >
                    ✕
                </button>
            </div>
            <div class="flex-1 overflow-auto bg-latte-surface2 dark:bg-mocha-surface2 p-4 rounded font-mono text-sm">
                <!-- <pre class="whitespace-pre-wrap">{logs}</pre> -->
				<div bind:this={terminalElement} class="terminal-container"></div>
            </div>
        </div>
    </div>
{/if}