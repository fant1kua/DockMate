<script lang="ts">
	import {Toaster} from 'svelte-5-french-toast';
	import {
		QuitApp,
		MaximiseApp,
		MinimiseApp,
	} from "@app/app/App";
    import ContainersList from "./components/ContainersList.svelte";
    import ImagesList from "./components/ImagesList.svelte";
    import NetworksList from "./components/NetworksList.svelte";
    import VolumesList from "./components/VolumesList.svelte";

	type IPage = 'containers' | 'images' | 'volumes' | 'networks'

	let title = "DockMate";
	let page = $state<IPage>('containers')

	function QuitButton() {
		QuitApp();
	}

	function MaximiseButton() {
		MaximiseApp();
	}

	function MinimiseButton() {
		MinimiseApp();
	}

	function handleSetPage(p: IPage) {
		page = p
	}

</script>

<main
	class="w-full h-full bg-latte-base dark:bg-mocha-base text-latte-text dark:text-mocha-text grid grid-rows-[auto_1fr]"
>
	<drag
		class="w-full bg-latte-surface0 dark:bg-mocha-surface0 cursor-grabbing"
		style="--wails-draggable:drag"
	>
		<div class="w-full grid grid-cols-3 select-none">
			<div></div>
			<div class="text-center py-1">{title}</div>
			<div class="text-right">
				<button
					class="font-bold px-2 py-1"
					onclick={() => {
						MinimiseButton();
					}}>_</button
				><button
					class="font-bold px-2 py-1"
					onclick={() => {
						MaximiseButton();
					}}>O</button
				><button
					class="font-bold px-2 py-1"
					onclick={() => {
						QuitButton();
					}}>X</button
				>
			</div>
		</div>
	</drag>

	<div
		class="grid grid-cols-[60px_1fr] border-l border-b border-r border-latte-surface0 dark:border-mocha-surface0 overflow-y-auto"
	>
		<nav class="bg-latte-surface1 dark:bg-mocha-surface1 p-2">
		<div class="flex flex-col gap-2">
			<button
				aria-label="Containers"
				class="p-2 hover:bg-latte-surface2 dark:hover:bg-mocha-surface2 rounded flex items-center gap-2 {page === 'containers' ? 'bg-latte-surface2 dark:bg-mocha-surface2' : ''}"
				onclick={() => handleSetPage('containers')}
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><path fill="currentColor" d="M29 12h-9V3h9zm-7-2h5V5h-5z"/><path fill="currentColor" d="M20 15v2h7v10H17V3H5c-1.103 0-2 .898-2 2v22c0 1.103.897 2 2 2h22c1.103 0 2-.897 2-2V15zM6.414 17H15v8.586zM15 15H6.414L15 6.414zM13.586 5L5 13.586V5zM5 18.414L13.586 27H5z"/></svg>
			</button>
			<button
				aria-label="Images"
				class="p-2 hover:bg-latte-surface2 dark:hover:bg-mocha-surface2 rounded flex items-center gap-2 {page === 'images' ? 'bg-latte-surface2 dark:bg-mocha-surface2' : ''}"
				onclick={() => handleSetPage('images')}
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><path fill="currentColor" d="M13.983 11.978h3.994V16h-3.994zm-3.994 3.994h3.994v3.994h-3.994zm3.994 0h3.994v3.994h-3.994zm3.994-3.994h3.994V16h-3.994zm-3.994 0h3.994V16h-3.994zm-3.994-3.994h3.994v3.994h-3.994zm3.994 0h3.994v3.994h-3.994zm3.994 0h3.994v3.994h-3.994z"/><path fill="currentColor" d="M29.655 14.745V6.25l-6.164-3.55-6.164 3.55v3.55l6.164 3.55l6.164-3.55zm-6.164-7.1l4.109 2.365l-4.109 2.365l-4.109-2.365l4.109-2.365zm-6.164 7.1v8.495l6.164 3.55l6.164-3.55v-3.55l-6.164-3.55l-6.164 3.55zm0 8.495v-3.55l6.164-3.55l6.164 3.55v3.55l-6.164 3.55l-6.164-3.55z"/></svg>
			</button>
			<button
				aria-label="Volumes"
				class="p-2 hover:bg-latte-surface2 dark:hover:bg-mocha-surface2 rounded flex items-center gap-2 {page === 'volumes' ? 'bg-latte-surface2 dark:bg-mocha-surface2' : ''}"
				onclick={() => handleSetPage('volumes')}
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 48 48"><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" d="m37.876 29.287l-4.627-2.644L37.875 24l4.625 2.643v5.286L24 42.5L5.5 31.929v-5.286L24 16.073l9.25 5.285l-4.625 2.643L24 21.358L5.5 31.929" stroke-width="1"/><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" d="M42.5 16.071L24 26.642l-4.624-2.643l-4.626 2.643L24 31.928l18.5-10.571zL24 5.5L5.5 16.071v5.286L10.125 24l4.626-2.643l-4.627-2.644" stroke-width="1"/><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" d="M37.876 18.715L24 10.785L5.5 21.358m37 5.285L24 37.214l-13.876-7.93M24 37.214V42.5m0-15.858v5.286m0-15.855v5.285M24 5.5v5.286" stroke-width="1"/></svg>
			</button>
			<button
				aria-label="Networks"
				class="p-2 hover:bg-latte-surface2 dark:hover:bg-mocha-surface2 rounded flex items-center gap-2 {page === 'networks' ? 'bg-latte-surface2 dark:bg-mocha-surface2' : ''}"
				onclick={() => handleSetPage('networks')}
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="currentColor" fill-rule="evenodd" d="M5.15 1.23c-.152.368-.152.833-.152 1.77s0 1.4.152 1.77c.203.49.592.879 1.08 1.08c.368.152.833.152 1.77.152v2H3.5a1.5 1.5 0 0 0-1.5 1.5v.513c-.327.017-.561.054-.765.139c-.49.203-.879.592-1.08 1.08c-.152.368-.152.834-.152 1.77s0 1.4.152 1.77c.203.49.592.879 1.08 1.08c.368.152.833.152 1.77.152s1.4 0 1.77-.152c.49-.203.879-.592 1.08-1.08c.152-.368.152-.834.152-1.77s0-1.4-.152-1.77a2 2 0 0 0-1.08-1.08c-.368-.152-.833-.152-1.77-.152v-.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 .5.5v.5c-.932 0-1.4 0-1.77.152c-.49.203-.879.592-1.08 1.08c-.152.368-.152.834-.152 1.77s0 1.4.152 1.77c.203.49.592.879 1.08 1.08c.368.152.834.152 1.77.152s1.4 0 1.77-.152c.49-.203.879-.592 1.08-1.08c.152-.368.152-.834.152-1.77s0-1.4-.152-1.77a2 2 0 0 0-1.08-1.08c-.204-.084-.438-.122-.765-.139v-.513a1.5 1.5 0 0 0-1.5-1.5h-3.5v-2.01c.327-.017.561-.054.765-.139c.49-.203.879-.592 1.08-1.08c.152-.368.152-.833.152-1.77s0-1.4-.152-1.77a2 2 0 0 0-1.08-1.08C9.407.001 8.942.001 8.005.001s-1.4 0-1.77.152c-.49.203-.879.592-1.08 1.08zm-4.04 10.3c-.109.214-.109.494-.109 1.05v.8c0 .56 0 .84.109 1.05c.096.188.249.341.437.437c.214.109.494.109 1.05.109h.8c.56 0 .84 0 1.05-.109c.188-.096.341-.249.437-.437c.109-.214.109-.494.109-1.05v-.8c0-.56 0-.84-.109-1.05a1 1 0 0 0-.437-.437c-.214-.109-.494-.109-1.05-.109h-.8c-.56 0-.84 0-1.05.109a1 1 0 0 0-.437.437M11 12.58c0-.56 0-.84.109-1.05c.096-.188.249-.341.437-.437c.214-.109.494-.109 1.05-.109h.8c.56 0 .84 0 1.05.109c.188.096.341.249.437.437c.109.214.109.494.109 1.05v.8c0 .56 0 .84-.109 1.05a1 1 0 0 1-.437.437c-.214.109-.494.109-1.05.109h-.8c-.56 0-.84 0-1.05-.109a1 1 0 0 1-.437-.437C11 14.216 11 13.936 11 13.38zM6.11 1.48c-.109.214-.109.494-.109 1.05v.8c0 .56 0 .84.109 1.05c.096.188.249.341.437.437c.214.109.494.109 1.05.109h.8c.56 0 .84 0 1.05-.109c.188-.096.341-.249.437-.437c.109-.214.109-.494.109-1.05v-.8c0-.56 0-.84-.109-1.05a1 1 0 0 0-.437-.437C9.233.934 8.953.934 8.397.934h-.8c-.56 0-.84 0-1.05.109a1 1 0 0 0-.437.437" clip-rule="evenodd"/></svg>
			</button>
		</div>
		</nav>

		<content class="p-2">
			{#if page === 'containers'}
				<ContainersList />
			{:else if page === 'images'}
				<ImagesList />
			{:else if page === 'volumes'}
				<VolumesList />
			{:else if page === 'networks'}
				<NetworksList />
			{/if}
		</content>
	</div>
</main>

<Toaster />
