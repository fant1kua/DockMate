<script lang="ts">
    import { isError, formatBytes } from "../utils";
    import {
        DeleteImage,
        ListImages,
	} from "../../wailsjs/go/app/App";
    import type { app } from "../../wailsjs/go/models";

    import toast from 'svelte-5-french-toast';
    import CopyBtn from "./CopyBtn.svelte";

    let images = $state<app.ImageInfo[]>([]);
    let loading = $state<boolean>(false);
    let inAction = $state<boolean>(false);

    async function loadImages() {
        if (inAction) {
            return;
        }
		loading = true;
		try {
			images = await ListImages() ?? [];
		} catch (e) {
			toast.error(e instanceof Error ? e.message : 'Failed to load containers');
		} finally {
			loading = false;
		}
	}

    async function handleDeleteImage(id:string) {
        try {
            inAction = true
			await DeleteImage(id);
            toast.success('Image delete');
			await loadImages();
		} catch (e) {
            toast.error(isError(e) ? e.message : 'Failed to delete image');
		} finally {
            inAction = false
        }
    }

    $effect(() => {
        loadImages();
        const interval = setInterval(loadImages, 3000);
        return () => clearInterval(interval);
    });
</script>

<button 
    class="w-full bg-latte-surface2 dark:bg-mocha-surface2 p-2 rounded hover:bg-latte-surface3 dark:hover:bg-mocha-surface3"
    onclick={() => loadImages()}
    disabled={loading || inAction}
>
    {loading ? 'Loading...' : 'Refresh Images'}
</button>
<div class="grid gap-4">
    {#if !images || images.length === 0}
        <div class="text-center text-gray-500">No images found</div>
    {:else}
        {#each images as image}
            <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded">
                <div class="grid grid-cols-2 gap-2">
                    <div class="font-bold">ID:</div>
                    <div class="flex items-center gap-2">
                        <span class="truncate max-w-[200px]">{image.id}</span>
                        <CopyBtn value={image.id} />
                    </div>

                    <div class="font-bold">Tags:</div>
                    <div>{image.tags.join(',')}</div>

                    <div class="font-bold">Size:</div>
                    <div>{formatBytes(image.size)} ({image.size} bytes)</div>

                    <div class="font-bold">Created at:</div>
                    <div>{image.createdAt}</div>
                </div>

                <div class="mt-4 flex gap-2">
                    <button 
                        class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                        onclick={() => handleDeleteImage(image.id)}
                        disabled={inAction}
                    >
                        Delete
                    </button>
                </div>
            </div>
        {/each}
    {/if}
</div>