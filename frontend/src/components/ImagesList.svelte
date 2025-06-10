<script lang="ts">
    import toast from 'svelte-5-french-toast';
    import { isError, formatBytes } from "../utils";
    import {
        List,
        Remove,
        CreateAndStart,
        StartWatching,
        StopWatching,
        Save,
	} from "@app/app/DockerImagesService";
    import type { app } from "@app/models";
    import { EventsOff, EventsOn } from "@runtime/runtime";
    import CopyBtn from "./CopyBtn.svelte";
    import Inspect from './Inspect.svelte';

    type IAction = 'inspect'

    let list = $state<app.ImageInfo[]>([]);
    let loading = $state<boolean>(false);
    let inAction = $state<boolean>(false);
    let image = $state<app.ImageInfo | null>(null)
    let action = $state<IAction>('inspect')

    async function load() {
        if (inAction) {
            return;
        }
        loading = true;
        try {
            list = await List();
        } catch (e) {
            toast.error(e instanceof Error ? e.message : 'Failed to load images');
        } finally {
            loading = false;
        }
    }

    async function handleDeleteImage(id:string) {
        try {
            inAction = true
			await Remove(id);
            toast.success('Image deleted');
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to delete image');
		} finally {
            inAction = false
        }
    }

    async function handleStartContainer(id:string) {
        try {
            inAction = true
			await CreateAndStart(id);
            toast.success('Container started successfully');
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to start container');
		} finally {
            inAction = false
        }
    }

    async function handleSave(id:string) {
        try {
            inAction = true
			await Save(id);
            toast.success('Container started successfully');
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to start container');
		} finally {
            inAction = false
        }
    }

    function handleAction(c: app.ImageInfo, act: IAction) {
        image = c
        action = act
    }

    function handleClose() {
		image = null
	}


    $effect(() => {
        load();
    });

    $effect(() => {
        EventsOn("docker:images", (l: app.ImageInfo[]) => {
          list = l
        });
        StartWatching();

        return () => {
            StopWatching();
            EventsOff('docker:images');
        }
    });
</script>

<button 
    class="w-full bg-latte-surface2 dark:bg-mocha-surface2 p-2 rounded hover:bg-latte-surface3 dark:hover:bg-mocha-surface3"
    onclick={() => load()}
    disabled={loading || inAction}
>
    {loading ? 'Loading...' : 'Refresh Images'}
</button>
<div class="grid gap-4">
    {#if !list || list.length === 0}
        <div class="text-center text-gray-500">No images found</div>
    {:else}
        {#each list as item}
            <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded">
                <div class="grid grid-cols-2 gap-2">
                    <div class="font-bold">ID:</div>
                    <div class="flex items-center gap-2">
                        <span class="truncate max-w-[200px]">{item.id}</span>
                        <CopyBtn value={item.id} />
                    </div>

                    <div class="font-bold">Tags:</div>
                    <div>{item.tags.join(',')}</div>

                    <div class="font-bold">Size:</div>
                    <div>{formatBytes(item.size)} ({item.size} bytes)</div>

                    <div class="font-bold">Created at:</div>
                    <div>{item.createdAt}</div>
                </div>

                <div class="mt-4 flex gap-2">
                    <button
                        aria-label="Start Container"
                        class="text-green-500 hover:text-green-600 px-2 py-1 rounded disabled:opacity-50"
                        onclick={() => handleStartContainer(item.id)}
                        disabled={inAction}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M8 5v14l11-7z"/></svg>
                    </button>

                    <button
                        aria-label="Save"
                        class="text-yellow-500 hover:text-yellow-600 px-2 py-1 rounded disabled:opacity-50"
                        onclick={() => handleSave(item.id)}
                        disabled={inAction}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M21 7v12q0 .825-.587 1.413T19 21H5q-.825 0-1.412-.587T3 19V5q0-.825.588-1.412T5 3h12zm-9 11q1.25 0 2.125-.875T15 15t-.875-2.125T12 12t-2.125.875T9 15t.875 2.125T12 18m-6-8h9V6H6z"/></svg>
                    </button>

                    <button
                        aria-label="Remove"
                        class="text-red-500 hover:text-red-600 px-2 py-1 rounded disabled:opacity-50"
                        onclick={() => handleDeleteImage(item.id)}
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
<Inspect type="image" id={action === 'inspect' ? image?.id : null} onClose={handleClose} />