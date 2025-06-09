<script lang="ts">
    import toast from "svelte-5-french-toast";
    import { Inspect as ContainerInspect } from "@app/app/DockerContainersService";
    import { Inspect as ImageInspect } from "@app/app/DockerImagesService";
    import { Inspect as VolumeInspect } from "@app/app/DockerVolumesService";
    import { Inspect as NetworkInspect } from "@app/app/DockerNetworksService";
    
    let { type, id, onClose } = $props<{
        type: 'container' | 'image' | 'volume' | 'network'
        id: string | null;
		onClose(): void;
    }>();

    let inspectData = $state<string>('');

    async function loadContainerData() {
        if (!id) return;
        try {
            const data = await ContainerInspect(id);
            inspectData = JSON.stringify(JSON.parse(data), null, 2);
        } catch (e) {
            toast.error('Failed to inspect container');
        }
    }

    async function loadImageData() {
        if (!id) return;
        try {
            const data = await ImageInspect(id);
            inspectData = JSON.stringify(JSON.parse(data), null, 2);
        } catch (e) {
            toast.error('Failed to inspect iamge');
        }
    }

    async function loadVolumeData() {
        if (!id) return;
        try {
            const data = await VolumeInspect(id);
            inspectData = JSON.stringify(JSON.parse(data), null, 2);
        } catch (e) {
            toast.error('Failed to inspect volume');
        }
    }

    async function loadNetworkData() {
        if (!id) return;
        try {
            const data = await NetworkInspect(id);
            inspectData = JSON.stringify(JSON.parse(data), null, 2);
        } catch (e) {
            toast.error('Failed to inspect network');
        }
    }

    $effect(() => {
        if (!id) {
            return;
        }

        switch (type) {
            case 'container': {
                loadContainerData();
            } break;
            case 'image': {
                loadImageData();
            } break;
            case 'volume': {
                loadVolumeData();
            } break;
            case 'network': {
                loadNetworkData();
            } break;
        }
    });
</script>

{#if id}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
        <div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded-lg w-3/4 h-3/4 flex flex-col">
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-xl font-bold">Inspect</h2>
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