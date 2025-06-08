<script lang="ts">
    import toast from "svelte-5-french-toast";
    import { isError } from '../utils'
    import {
		ListContainers,
		StartContainer,
		StopContainer,
		RestartContainer,
		RemoveContainer,
        KillContainer,
	} from "../../wailsjs/go/app/App";
    import type { app } from "../../wailsjs/go/models";
    import ConteinerInspect from "./ConteinerInspect.svelte";
    import ConteinerLogs from "./ConteinerLogs.svelte";
    import ContainerTerminal from "./ContainerTerminal.svelte";
    import CopyBtn from "./CopyBtn.svelte";

    let containers = $state<app.ContainerInfo[]>([]);
    let container = $state<app.ContainerInfo | null>(null)
    let action = $state<'logs' | 'inspect' | 'terminal'>('logs')
    let loading = $state<boolean>(false);
    let inAction = $state<boolean>(false);

    function handleViewLogs(c: app.ContainerInfo) {
		container = c
        action = 'logs'
	}

    function handleViewInspect(c: app.ContainerInfo) {
		container = c
        action = 'inspect'
	}

    function handleViewTerminal(c: app.ContainerInfo) {
        container = c
        action = 'terminal'
    }

    function handleClose() {
		container = null
	}

    async function loadContainers() {
        if (inAction || container !== null) {
            return
        }
		loading = true;
		try {
			containers = await ListContainers() ?? [];
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to load containers');
		} finally {
			loading = false;
		}
	}

    async function handleStartContainer(id: string) {
		try {
            inAction = true
			await StartContainer(id);
            toast.success('Container started');
			await loadContainers();
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to start container');
		} finally {
            inAction = false
        }
	}

	async function handleStopContainer(id: string) {
		try {
            inAction = true
			await StopContainer(id);
            toast.success('Container stopped');
			await loadContainers();
		} catch (e) {
			toast.error(isError(e) ? e.message : 'Failed to stop container');
		} finally {
            inAction = false
        }
	}

	async function handleRestartContainer(id: string) {
		try {
            inAction = true
			await RestartContainer(id);
            toast.success('Container restarted');
			await loadContainers();
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to restart container');
		} finally {
            inAction = false
        }
	}

    async function handleKillContainer(id: string) {
		try {
            inAction = true
			await KillContainer(id);
            toast.success('Container killed');
			await loadContainers();
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to kill container');
		} finally {
            inAction = false
        }
	}

	async function handleRemoveContainer(id: string) {
		try {
            inAction = true
			await RemoveContainer(id);
            toast.success('Container removed');
			await loadContainers();
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to remove container');
		} finally {
            inAction = false
        }
	}

    $effect(() => {
        loadContainers();
        const interval = setInterval(loadContainers, 3000);
        return () => clearInterval(interval);
    });
</script>

<button 
    class="w-full bg-latte-surface2 dark:bg-mocha-surface2 p-2 rounded hover:bg-latte-surface3 dark:hover:bg-mocha-surface3"
    onclick={() => loadContainers()}
    disabled={loading || inAction}
>
    {loading ? 'Loading...' : 'Refresh Containers'}
</button>
<div class="grid gap-4"> 
    {#if !containers || containers.length === 0}
        <div class="text-center text-gray-500">No containers found</div>
    {:else}
        {#each containers as container}
            <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded">
                <div class="grid grid-cols-2 gap-2">
                    <div class="font-bold">ID:</div>
                    <div class="flex items-center gap-2">
                        <span class="truncate max-w-[200px]">{container.id}</span>
                        <CopyBtn value={container.id} />
                    </div>
                    
                    <div class="font-bold">Names:</div>
                    <div>{container.names.join(', ')}</div>
                    
                    <div class="font-bold">Image:</div>
                    <div class="flex items-center gap-2">
                        <span class="truncate max-w-[200px]">{container.image}</span>
                        <CopyBtn value={container.image} />
                    </div>
                    
                    <div class="font-bold">Status:</div>
                    <div>{container.status}</div>
                    
                    <div class="font-bold">State:</div>
                    <div>{container.state}</div>
                </div>
                
                <div class="mt-4 flex gap-2">
                    <button
                        aria-label="Start Container"
                        class="text-green-500 hover:text-green-600 px-3 py-1 rounded disabled:opacity-50"
                        onclick={() => handleStartContainer(container.id)}
                        disabled={inAction}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M8 5v14l11-7z"/></svg>
                    </button>
                    <button 
                        aria-label="Stop"
                        class="text-red-500 hover:text-red-600 px-3 py-1 rounded disabled:opacity-50"
                        onclick={() => handleStopContainer(container.id)}
                        disabled={inAction || container.state !== 'running'}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="currentColor" d="M8 0a8 8 0 1 0 0 16A8 8 0 0 0 8 0m0 14.5a6.5 6.5 0 1 1 0-13a6.5 6.5 0 0 1 0 13M5 5h6v6H5z"/></svg>                    </button>
                    <button
                        aria-label="Restart"
                        class="text-yellow-500 hover:text-yellow-600 px-3 py-1 rounded disabled:opacity-50"
                        onclick={() => handleRestartContainer(container.id)}
                        disabled={inAction || container.state !== 'running'}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M12 4c2.1 0 4.1.8 5.6 2.3c3.1 3.1 3.1 8.2 0 11.3c-1.8 1.9-4.3 2.6-6.7 2.3l.5-2c1.7.2 3.5-.4 4.8-1.7c2.3-2.3 2.3-6.1 0-8.5C15.1 6.6 13.5 6 12 6v4.6l-5-5l5-5zM6.3 17.6C3.7 15 3.3 11 5.1 7.9l1.5 1.5c-1.1 2.2-.7 5 1.2 6.8q.75.75 1.8 1.2l-.6 2q-1.5-.6-2.7-1.8"/></svg>
                    </button>
                    <button
                        aria-label="Kill"
                        class="text-red-500 hover:text-red-600 px-3 py-1 rounded disabled:opacity-50"
                        onclick={() => handleKillContainer(container.id)}
                        disabled={inAction || container.state !== 'running'}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 64 64"><path fill="currentColor" d="M62 10.571L53.429 2L32 23.429L10.571 2L2 10.571L23.429 32L2 53.429L10.571 62L32 40.571L53.429 62L62 53.429L40.571 32z"/></svg>                    </button>
                    <button
                        aria-label="Remove"
                        class="text-red-500 hover:text-red-600 px-3 py-1 rounded disabled:opacity-50"
                        onclick={() => handleRemoveContainer(container.id)}
                        disabled={inAction}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="currentColor" fill-rule="evenodd" d="M5.75 3V1.5h4.5V3zm-1.5 0V1a1 1 0 0 1 1-1h5.5a1 1 0 0 1 1 1v2h2.5a.75.75 0 0 1 0 1.5h-.365l-.743 9.653A2 2 0 0 1 11.148 16H4.852a2 2 0 0 1-1.994-1.847L2.115 4.5H1.75a.75.75 0 0 1 0-1.5zm-.63 1.5h8.76l-.734 9.538a.5.5 0 0 1-.498.462H4.852a.5.5 0 0 1-.498-.462z" clip-rule="evenodd"/></svg>
                    </button>
                    <button 
                        class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded"
                        onclick={() => handleViewLogs(container)}
                    >
                        View Logs
                    </button>
                    <button 
                        class="bg-green-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                        onclick={() => handleViewInspect(container)}
                    >
                        Inspect
                    </button>
                    <button 
                        class="bg-purple-500 hover:bg-purple-600 text-white px-3 py-1 rounded"
                        onclick={() => handleViewTerminal(container)}
                        disabled={container.state !== 'running'}
                    >
                        Terminal
                    </button>
                </div>
            </div>
        {/each}
    {/if}
</div>
<ConteinerLogs  container={action === 'logs' ? container : null} onClose={handleClose} />
<ConteinerInspect  container={action === 'inspect' ? container : null} onClose={handleClose} />
<ContainerTerminal container={action === 'terminal' ? container : null} onClose={handleClose} />