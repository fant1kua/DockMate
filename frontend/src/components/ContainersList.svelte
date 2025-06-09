<script lang="ts">
    import toast from 'svelte-5-french-toast';
    import { isError } from "../utils";
    import {
        List,
        StartWatching,
        StopWatching,
        Start,
        Stop,
        Restart,
        Remove,
        Kill,
    } from '@app/app/DockerContainersService';
    import type { app } from "@app/models";
    import { EventsOff, EventsOn } from "@runtime/runtime";
    import CopyBtn from "./CopyBtn.svelte";
    import ConteinerLogs from "./ConteinerLogs.svelte";
    import ConteinerInspect from "./ConteinerInspect.svelte";
    import ContainerTerminal from "./ContainerTerminal.svelte";

    type IAction = 'logs' | 'inspect' | 'terminal'

    let list = $state<app.ContainersGroup[]>([]);
    let loading = $state<boolean>(false);
    let container = $state<app.ContainerInfo | null>(null)
    let action = $state<IAction>('logs')
    let inAction = $state<boolean>(false);

    async function load() {
        if (inAction) {
            return;
        }
        loading = true;
        try {
            list = await List();
        } catch (e) {
            toast.error(e instanceof Error ? e.message : 'Failed to load containers');
        } finally {
            loading = false;
        }
    }

    async function handleStartContainer(id: string) {
        try {
            inAction = true;
            await Start(id);
            toast.success('Container started');
        } catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to start container');
        } finally {
            inAction = false;
        }
    }

    async function handleStopContainer(id: string) {
        try {
            inAction = true;
            await Stop(id);
            toast.success('Container stopped');
        } catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to stop container');
        } finally {
            inAction = false;
        }
    }

    async function handleRestartContainer(id: string) {
        try {
            inAction = true;
            await Restart(id);
            toast.success('Container restarted');
        } catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to restart container');
        } finally {
            inAction = false;
        }
    }

    async function handleRemoveContainer(id: string) {
        try {
            inAction = true;
            await Remove(id);
            toast.success('Container removed');
        } catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to remove container');
        } finally {
            inAction = false;
        }
    }

    async function handleKillContainer(id: string) {
        try {
            inAction = true;
            await Kill(id);
            toast.success('Container killed');
        } catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to kill container');
        } finally {
            inAction = false;
        }
    }

    function handleAction(c: app.ContainerInfo, act: IAction) {
        container = c
        action = act
    }

    function handleClose() {
		container = null
	}

    $effect(() => {
        load();
    });

    $effect(() => {
        EventsOn("docker:containers", (l: app.ContainersGroup[]) => {
          list = l
        });

        StartWatching();

        return () => {
            EventsOff('docker:containers');
            StopWatching();
        }
    });
</script>

<button 
    class="w-full bg-latte-surface2 dark:bg-mocha-surface2 p-2 rounded hover:bg-latte-surface3 dark:hover:bg-mocha-surface3"
    onclick={() => load()}
    disabled={loading || inAction}
>
    {loading ? 'Loading...' : 'Refresh Containers'}
</button>

<div class="grid gap-4">
    {#if !list || list.length === 0}
        <div class="text-center text-gray-500">No containers found</div>
    {:else}
        {#each list as project}
            <details class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded" open>
                <summary class="text-xl font-bold mb-1">{project.name}</summary>
                {#if !project.containers || project.containers.length === 0}
                    <div class="text-center text-gray-500">No containers in this project</div>
                {:else}
                    {#each project.containers as container}
                        <div class="bg-latte-surface2 dark:bg-mocha-surface2 p-4 rounded mb-2">
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
                                <div class="flex items-center gap-2">
                                    <span class={`inline-block w-2 h-2 rounded-full ${
                                        container.state === 'running' ? 'bg-green-500' : 'bg-red-500'
                                    }`}></span>
                                    {container.status}
                                </div>
                            </div>

                            <div class="mt-4 flex gap-2">
                                {#if container.state !== 'running'}
                                    <button
                                        aria-label="Start"
                                        class="text-green-500 hover:text-green-600 px-3 py-1 rounded disabled:opacity-50"
                                        onclick={() => handleStartContainer(container.id)}
                                        disabled={inAction}
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M8 5v14l11-7z"/></svg>
                                    </button>
                                {:else}
                                    <button
                                        aria-label="Stop"
                                        class="text-yellow-500 hover:text-yellow-600 px-3 py-1 rounded disabled:opacity-50"
                                        onclick={() => handleStopContainer(container.id)}
                                        disabled={inAction}
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M6 6h12v12H6z"/></svg>
                                    </button>
                                    <button
                                        aria-label="Restart"
                                        class="text-blue-500 hover:text-blue-600 px-3 py-1 rounded disabled:opacity-50"
                                        onclick={() => handleRestartContainer(container.id)}
                                        disabled={inAction}
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M17.65 6.35A7.958 7.958 0 0 0 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/></svg>
                                    </button>
                                    <button 
                                        class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded"
                                        onclick={() => handleAction(container, 'logs')}
                                    >
                                        View Logs
                                    </button>
                                    <button 
                                        class="bg-green-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                                        onclick={() => handleAction(container, 'inspect')}
                                    >
                                        Inspect
                                    </button>
                                    <button 
                                        class="bg-purple-500 hover:bg-purple-600 text-white px-3 py-1 rounded"
                                        onclick={() => handleAction(container, 'terminal')}
                                    >
                                        Terminal
                                    </button>
                                {/if}
                                <button
                                    aria-label="Kill"
                                    class="text-orange-500 hover:text-orange-600 px-3 py-1 rounded disabled:opacity-50"
                                    onclick={() => handleKillContainer(container.id)}
                                    disabled={inAction}
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10s10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8s8 3.59 8 8s-3.59 8-8 8zm4.59-12.42L10 14.17l-2.59-2.58L6 13l4 4l8-8z"/></svg>
                                </button>
                                <button
                                    aria-label="Remove"
                                    class="text-red-500 hover:text-red-600 px-3 py-1 rounded disabled:opacity-50"
                                    onclick={() => handleRemoveContainer(container.id)}
                                    disabled={inAction}
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="currentColor" fill-rule="evenodd" d="M5.75 3V1.5h4.5V3zm-1.5 0V1a1 1 0 0 1 1-1h5.5a1 1 0 0 1 1 1v2h2.5a.75.75 0 0 1 0 1.5h-.365l-.743 9.653A2 2 0 0 1 11.148 16H4.852a2 2 0 0 1-1.994-1.847L2.115 4.5H1.75a.75.75 0 0 1 0-1.5zm-.63 1.5h8.76l-.734 9.538a.5.5 0 0 1-.498.462H4.852a.5.5 0 0 1-.498-.462z" clip-rule="evenodd"/></svg>
                                </button>
                            </div>
                        </div>
                    {/each}
                {/if}
            </details>
        {/each}
    {/if}
</div>
<ConteinerLogs  container={action === 'logs' ? container : null} onClose={handleClose} />
<ConteinerInspect  container={action === 'inspect' ? container : null} onClose={handleClose} />
<ContainerTerminal container={action === 'terminal' ? container : null} onClose={handleClose} />