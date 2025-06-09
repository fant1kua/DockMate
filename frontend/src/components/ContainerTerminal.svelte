<script lang="ts">
    import { Terminal } from 'xterm';
    import { FitAddon } from 'xterm-addon-fit';
    import { WebLinksAddon } from 'xterm-addon-web-links';
    import {
        StartInteractiveTerminal,
        SendToTerminal,
        CloseTerminal,
    } from "@app/app/DockerContainersTerminal";
    import type { app } from "@app/models";
    import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";

    let { container, onClose } = $props<{
        container: app.ContainerInfo | null;
		onClose(): void;
    }>();

    let terminalElement = $state<HTMLDivElement>();
    let terminal: Terminal;

    $effect(() => {
        if (!terminalElement || !container) return;

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

        // Add addons to terminal
        terminal.loadAddon(fitAddon);
        terminal.loadAddon(webLinksAddon);

        // Open terminal in the container
        terminal.open(terminalElement);
        fitAddon.fit();

        terminal.open(terminalElement);
        
        

        StartInteractiveTerminal(container.id, "111");

        let currentCommand = ''

        terminal.onKey((e) => {
            const char = e.key;
            const ev = e.domEvent;
            console.log(ev)

            if (ev.key === "Enter") {
                terminal.write('\r\n');
                if (currentCommand.trim()) {
                    console.log(currentCommand)
                    SendToTerminal("111", `${currentCommand}\n`).catch((error) => {
                        terminal.writeln(`Error executing command: ${error}`);
                    });
                    currentCommand = ''
                }
            } else if (ev.key === "Backspace") {
                // Do not delete prompt
                if (currentCommand !== '') {
                    terminal.write('\b \b');
                    currentCommand = currentCommand.slice(0, -1);
                }
            } else if (ev.ctrlKey && ev.key === 'c') {
                terminal.write('^C\r\n');
                currentCommand = '';
            } else {
                currentCommand += char;
                terminal.write(char);
            }
        });
        // terminal.onData((data) => {
        //     if (data === '\r') {
        //         // Enter key pressed
        //         terminal.write('\r\n');
        //         if (currentCommand.trim()) {
        //             console.log(currentCommand)
        //             SendToTerminal("111", `${currentCommand}\n`).catch((error) => {
        //                 terminal.writeln(`Error executing command: ${error}`);
        //             });
        //         }
        //         currentCommand = '';
        //     } else if (data === '\u0003') {
        //         // Ctrl+C
        //         terminal.write('^C\r\n');
        //         currentCommand = '';
        //     } else if (data === '\u007F') {
        //         // Backspace
        //         if (currentCommand.length > 0) {
        //             currentCommand = currentCommand.slice(0, -1);
        //             terminal.write('\b \b');
        //         }
        //     } else {
        //         // Regular character
        //         currentCommand += data;
        //         terminal.write(data);
        //     }
        // });

        EventsOn("docker:output", (line: string) => {
            console.info('output', line)
            terminal.writeln(line);
        });

        return () => {
            EventsOff('docker:output');
            CloseTerminal("111");
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