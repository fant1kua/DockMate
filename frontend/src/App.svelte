<script lang="ts">
	import {
		QuitApp,
		MaximiseApp,
		MinimiseApp,
	} from "../wailsjs/go/app/App";

    import ContainersList from "./components/ContainersList.svelte";
    import ImagesList from "./components/ImagesList.svelte";

	let title = "DockMate";
	let page = $state<'containers' | 'images'>('containers')
	

	function QuitButton() {
		QuitApp();
	}

	function MaximiseButton() {
		MaximiseApp();
	}

	function MinimiseButton() {
		MinimiseApp();
	}

	function handleSetPage(p: 'containers' | 'images') {
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
		<div class="w-full grid grid-cols-3">
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
				aria-label="Containers"
				class="p-2 hover:bg-latte-surface2 dark:hover:bg-mocha-surface2 rounded flex items-center gap-2 {page === 'images' ? 'bg-latte-surface2 dark:bg-mocha-surface2' : ''}"
				onclick={() => handleSetPage('images')}
			>
			<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20"><path fill="currentColor" d="M4 15v-3H2V2h12v3h2v3h2v10H6v-3zm7-12c-1.1 0-2 .9-2 2h4a2 2 0 0 0-2-2m-7 8V6H3v5zm7-3h4a2 2 0 1 0-4 0m-5 6V9H5v5zm9-1a2 2 0 1 0 .001-3.999A2 2 0 0 0 15 13m2 4v-2c-5 0-5-3-10-3v5z"/></svg>	
			</button>
		</div>
		</nav>

		<content class="p-2">
			{#if page === 'containers'}
				<ContainersList />
			{:else if page === 'images'}
				<ImagesList />
			{/if}
		</content>
	</div>
</main>
