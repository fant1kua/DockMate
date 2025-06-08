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
    import CopyBtn from "./CopyBtn.svelte";

    let containers = $state<app.ContainerInfo[]>([]);
    let container = $state<app.ContainerInfo | null>(null)
    let action = $state<'logs' | 'inspect'>('logs')
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

    function handleClose() {
		container = null
	}

    async function loadContainers() {
        if (inAction) {
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
                    <div>{container.image}</div>
                    
                    <div class="font-bold">Status:</div>
                    <div>{container.status}</div>
                    
                    <div class="font-bold">State:</div>
                    <div>{container.state}</div>
                </div>
                
                <div class="mt-4 flex gap-2">
                    <button 
                        class="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded"
                        onclick={() => handleStartContainer(container.id)}
                        disabled={inAction || container.state === 'running'}
                    >
                        Start
                    </button>
                    <button 
                        class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                        onclick={() => handleStopContainer(container.id)}
                        disabled={inAction || container.state !== 'running'}
                    >
                        Stop
                    </button>
                    <button 
                        class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded"
                        onclick={() => handleRestartContainer(container.id)}
                        disabled={inAction || container.state !== 'running'}
                    >
                        Restart
                    </button>
                    <button 
                        class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                        onclick={() => handleKillContainer(container.id)}
                        disabled={inAction || container.state !== 'running'}
                    >
                        Kill
                    </button>
                    <button 
                        class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                        onclick={() => handleRemoveContainer(container.id)}
                        disabled={inAction}
                    >
                        Remove
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
                </div>
            </div>
        {/each}
    {/if}
</div>
<ConteinerLogs  container={action === 'logs' ? container : null} onClose={handleClose} />
<ConteinerInspect  container={action === 'inspect' ? container : null} onClose={handleClose} />