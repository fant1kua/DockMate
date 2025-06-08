<script lang="ts">
    import toast from "svelte-5-french-toast";
    import {
        ListVolumes,
        DeleteVolume,
	} from "../../wailsjs/go/app/App";
    import type { app } from "../../wailsjs/go/models";
    import CopyBtn from "./CopyBtn.svelte";
    import { isError } from "../utils";

    let volumes = $state<app.VolumeInfo[]>([]);
    let loading = $state<boolean>(false);
	let error = $state<string | null>(null);
    let inAction = $state<boolean>(false);

    async function loadVolumes() {
        if (inAction) {
            return;
        }
		loading = true;
		error = null;
		try {
			volumes = await ListVolumes() ?? [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load containers';
		} finally {
			loading = false;
		}
	}

    async function handleDeleteVolume(name:string) {
        try {
            inAction = true
			await DeleteVolume(name);
            toast.success('Image delete');
			await loadVolumes();
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to delete volume');
		} finally {
            inAction = false
        }
    }

    $effect(() => {
        loadVolumes();
        const interval = setInterval(loadVolumes, 3000);
        return () => clearInterval(interval);
    });
</script>

<button 
    class="w-full bg-latte-surface2 dark:bg-mocha-surface2 p-2 rounded hover:bg-latte-surface3 dark:hover:bg-mocha-surface3"
    onclick={() => loadVolumes()}
    disabled={loading}
>
    {loading ? 'Loading...' : 'Refresh Volumes'}
</button>
<div class="grid gap-4">
    {#if error}
        <div class="text-red-500 mb-4">{error}</div>
    {/if}
    
    {#if !volumes || volumes.length === 0}
        <div class="text-center text-gray-500">No volumes found</div>
    {:else}
        {#each volumes as volume}
            <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded">
                <div class="grid grid-cols-2 gap-2">
                    <div class="font-bold">Name:</div>
                    <div class="flex items-center gap-2">
                        <span class="truncate max-w-[200px]">{volume.name}</span>
                        <CopyBtn value={volume.name} />
                    </div>

                    <div class="font-bold">Created at:</div>
                    <div>{volume.createdAt}</div>
                </div>
                <div class="mt-4 flex gap-2">
                    <button 
                        class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                        onclick={() => handleDeleteVolume(volume.name)}
                        disabled={inAction}
                    >
                        Delete
                    </button>
                </div>
            </div>
        {/each}
    {/if}
</div>