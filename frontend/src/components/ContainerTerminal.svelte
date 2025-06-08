<script lang="ts">
    import { Terminal } from 'xterm';
    import { FitAddon } from 'xterm-addon-fit';
    import { WebLinksAddon } from 'xterm-addon-web-links';
    import { WebglAddon } from 'xterm-addon-webgl';
    import { onMount, onDestroy } from 'svelte';
    import type { app } from "../../wailsjs/go/models";
    import { StreamContainerLogs, StopContainerLogs, ExecContainer, CloseContainerSession, GetCurrentPath } from "../../wailsjs/go/app/App";
    import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";

    let { container, onClose } = $props<{
        container: app.ContainerInfo | null;
		onClose(): void;
    }>();

    let terminalElement = $state<HTMLDivElement>();
    let terminal: Terminal;
    let currentCommand = '';
    let currentPath = '/';

    async function updatePrompt() {
        if (container) {
            currentPath = await GetCurrentPath(container.id);
            terminal.write(`\r\n${currentPath} $ `);
        }
    }

    $effect(() => {
        if (!terminalElement || !container) return;

        terminal = new Terminal({
            cursorBlink: true,
            cols: 80,
            rows: 24,
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

        // Open terminal in the DOM
        terminal.open(terminalElement);
        fitAddon.fit();

        // Handle window resize
        const handleResize = () => {
            fitAddon.fit();
        };
        window.addEventListener('resize', handleResize);

        // Start streaming container logs
        StreamContainerLogs(container.id);
        EventsOn("logStream", async (line: string) => {
            terminal.writeln(line);
            // Update path after each command
            await updatePrompt();
        });

        // Handle user input
        terminal.onData((data) => {
            if (data === '\r') {
                // Enter key pressed
                terminal.write('\r\n');
                if (currentCommand.trim()) {
                    ExecContainer(container.id, currentCommand).catch((error) => {
                        terminal.writeln(`Error executing command: ${error}`);
                        updatePrompt();
                    });
                    currentCommand = '';
                } else {
                    updatePrompt();
                }
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

        // Initial prompt
        updatePrompt();

        return () => {
            window.removeEventListener('resize', handleResize);
            EventsOff('logStream');
            StopContainerLogs();
            CloseContainerSession(container.id);
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