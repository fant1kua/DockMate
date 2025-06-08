<script lang="ts">
    import {
        ListImages,
	} from "../../wailsjs/go/app/App";
    import type { app } from "../../wailsjs/go/models";


    let images = $state<app.ImageInfo[]>([]);
    let loading = $state<boolean>(false);
	let error = $state<string | null>(null);

    function formatBytes(bytes: number): string {
        const units = ['B', 'KB', 'MB', 'GB', 'TB'];
        let size = bytes;
        let unitIndex = 0;
        while (size >= 1024 && unitIndex < units.length - 1) {
            size /= 1024;
            unitIndex++;
        }
        return `${size.toFixed(2)} ${units[unitIndex]}`;
    }

    async function loadImages() {
		loading = true;
		error = null;
		try {
			images = await ListImages() ?? [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load containers';
		} finally {
			loading = false;
		}
	}

    loadImages();
</script>

<button 
    class="w-full bg-latte-surface2 dark:bg-mocha-surface2 p-2 rounded hover:bg-latte-surface3 dark:hover:bg-mocha-surface3"
    onclick={() => loadImages()}
    disabled={loading}
>
    {loading ? 'Loading...' : 'Refresh Images'}
</button>
<div class="grid gap-4">
    {#if error}
        <div class="text-red-500 mb-4">{error}</div>
    {/if}
    
    {#if !images || images.length === 0}
        <div class="text-center text-gray-500">No containers found</div>
    {:else}
        {#each images as image}
            <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded">
                <div class="grid grid-cols-2 gap-2">
                    <div class="font-bold">ID:</div>
                    <div>{image.id}</div>

                    <div class="font-bold">Tags:</div>
                    <div>{image.tags.join(',')}</div>

                    <div class="font-bold">Size:</div>
                    <div>{formatBytes(image.size)} ({image.size} bytes)</div>

                    <div class="font-bold">Created at:</div>
                    <div>{image.createdAt}</div>
                </div>
            </div>
        {/each}
    {/if}
</div>