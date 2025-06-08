<script lang="ts">
    import { ContainerInspect } from "../../wailsjs/go/app/App";
    import { app } from "../../wailsjs/go/models";
  	import 'xterm/css/xterm.css';

    let { container, onClose } = $props<{
        container: app.ContainerInfo | null;
		onClose(): void;
    }>();

    let inspectData = $state<string>('');

    async function loadInspectData() {
        if (!container) return;
        try {
            const data = await ContainerInspect(container.id);
            inspectData = JSON.stringify(JSON.parse(data), null, 2);
        } catch (e) {
            console.error('Failed to inspect container:', e);
        }
    }

    $effect(() => {
        if (container) {
            loadInspectData();
        }
    });
</script>

{#if container}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
        <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded-lg w-3/4 h-3/4 flex flex-col">
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-xl font-bold">Inspect {container.names[0]}</h2>
                <button 
                    class="text-gray-500 hover:text-gray-700"
                    onclick={onClose}
                >
                    ✕
                </button>
            </div>
            <div class="flex-1 overflow-auto bg-latte-surface2 dark:bg-mocha-surface2 p-4 rounded font-mono text-sm">
                <pre class="whitespace-pre-wrap">{inspectData}</pre>
            </div>
        </div>
    </div>
{/if}