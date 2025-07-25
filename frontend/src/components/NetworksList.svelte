<script lang="ts">
    import toast from 'svelte-5-french-toast';
    import { isError } from "../utils";
    import {
        List,
        Remove,
        StartWatching,
        StopWatching,
	} from "@app/app/DockerNetworksService";
    import type { app } from "@app/models";
    import { EventsOff, EventsOn } from "@runtime/runtime";
    import CopyBtn from "./CopyBtn.svelte";
    import Inspect from './Inspect.svelte';

    type IAction = 'inspect'

    let list = $state<app.NetworkInfo[]>([]);
    let loading = $state<boolean>(false);
	let error = $state<string | null>(null);
    let inAction = $state<boolean>(false);
    let network = $state<app.NetworkInfo | null>(null)
    let action = $state<IAction>('inspect')

    async function load() {
        if (inAction) {
            return;
        }
        loading = true;
        try {
            list = await List();
        } catch (e) {
            toast.error(e instanceof Error ? e.message : 'Failed to load networks');
        } finally {
            loading = false;
        }
    }

    async function handleDeleteVolume(name:string) {
        try {
            inAction = true
			await Remove(name);
            toast.success('Network deleted');
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to delete network');
		} finally {
            inAction = false
        }
    }

    function handleAction(c: app.NetworkInfo, act: IAction) {
        network = c
        action = act
    }

    function handleClose() {
		network = null
	}

    $effect(() => {
        load();
    });

    $effect(() => {
        EventsOn("docker:networks", (l: app.NetworkInfo[]) => {
          list = l
        });
        StartWatching();

        return () => {
            StopWatching();
            EventsOff('docker:networks');
        }
    });
</script>

<button 
    class="w-full bg-latte-surface2 dark:bg-mocha-surface2 p-2 rounded hover:bg-latte-surface3 dark:hover:bg-mocha-surface3"
    onclick={() => load()}
    disabled={loading}
>
    {loading ? 'Loading...' : 'Refresh Networks'}
</button>
<div class="grid gap-4">
    {#if error}
        <div class="text-red-500 mb-4">{error}</div>
    {/if}
    
    {#if !list || list.length === 0}
        <div class="text-center text-gray-500">No networks found</div>
    {:else}
        {#each list as item}
            <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded">
                <div class="grid grid-cols-2 gap-2">
                    <div class="font-bold">ID:</div>
                    <div class="flex items-center gap-2">
                        <span class="truncate max-w-[200px]">{item.id}</span>
                        <CopyBtn value={item.id} />
                    </div>

                    <div class="font-bold">Name:</div>
                    <div>{item.name}</div>
                </div>
                <div class="mt-4 flex gap-2">
                    <button
                        aria-label="Remove"
                        class="text-red-500 hover:text-red-600 px-2 py-1 rounded disabled:opacity-50"
                        onclick={() => handleDeleteVolume(item.name)}
                        disabled={inAction}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="currentColor" fill-rule="evenodd" d="M5.75 3V1.5h4.5V3zm-1.5 0V1a1 1 0 0 1 1-1h5.5a1 1 0 0 1 1 1v2h2.5a.75.75 0 0 1 0 1.5h-.365l-.743 9.653A2 2 0 0 1 11.148 16H4.852a2 2 0 0 1-1.994-1.847L2.115 4.5H1.75a.75.75 0 0 1 0-1.5zm-.63 1.5h8.76l-.734 9.538a.5.5 0 0 1-.498.462H4.852a.5.5 0 0 1-.498-.462z" clip-rule="evenodd"/></svg>
                    </button>

                    <button 
                        class="bg-green-500 hover:bg-red-600 text-white px-2 py-1 rounded"
                        onclick={() => handleAction(item, 'inspect')}
                    >
                        Inspect
                    </button>
                </div>
            </div>
        {/each}
    {/if}
</div>
<Inspect type="network" id={action === 'inspect' ? network?.id : null} onClose={handleClose} />