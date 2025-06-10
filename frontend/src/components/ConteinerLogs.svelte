<script lang="ts">
    import { Terminal } from 'xterm';
  	import 'xterm/css/xterm.css';
    import { FitAddon } from 'xterm-addon-fit';
    import { WebLinksAddon } from 'xterm-addon-web-links';
    import {
        StartWatching,
        StopWatching,
    } from "@app/app/DockerLogsService";
    import type { app } from "@app/models";
    import { EventsOff, EventsOn } from "@runtime/runtime";

    let { container, onClose } = $props<{
        container: app.ContainerInfo | null;
		onClose(): void;
    }>();

  	let containerElement = $state<HTMLDivElement>();

    $effect(() => {
        if (!containerElement) return;

        const terminal = new Terminal({
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

        // Add addons to terminal
        terminal.loadAddon(fitAddon);
        terminal.loadAddon(webLinksAddon);

        // Open terminal in the container
        requestAnimationFrame(() => {
            terminal.open(containerElement);
            fitAddon.fit();
            terminal.focus();
        });

        EventsOn("docker:logs", (data: string) => {
            terminal.writeln(data);
        });
        StartWatching(container.id);

        return () => {
            StopWatching();
            EventsOff('docker:logs');
            terminal.dispose();
        }
    });

</script>

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
    <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded-lg w-3/4 h-3/4 flex flex-col">
        <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-bold">Logs for {container.names[0]}</h2>
            <button 
                class="text-gray-500 hover:text-gray-700"
                onclick={onClose}
            >
                ✕
            </button>
        </div>
        <div class="flex-1 overflow-auto bg-latte-surface2 dark:bg-mocha-surface2 p-4 rounded font-mono text-sm">
            <!-- <pre class="whitespace-pre-wrap">{logs}</pre> -->
            <div bind:this={containerElement} class="terminal-container"></div>
        </div>
    </div>
</div>