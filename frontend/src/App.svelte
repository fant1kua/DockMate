<script lang="ts">
	import { QuitApp, MaximiseApp, MinimiseApp, ListContainers, StartContainer, StopContainer, RestartContainer } from "../wailsjs/go/app/App";
	import type { app } from "../wailsjs/go/models";

	let title = "DockMate";
	let containers: app.ContainerInfo[] = [];
	let loading = false;
	let error: string | null = null;

	async function loadContainers() {
		loading = true;
		error = null;
		try {
			containers = await ListContainers();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load containers';
		} finally {
			loading = false;
		}
	}

	function QuitButton() {
		QuitApp();
	}

	function MaximiseButton() {
		MaximiseApp();
	}

	function MinimiseButton() {
		MinimiseApp();
	}

	async function handleStartContainer(id: string) {
		try {
			await StartContainer(id);
			await loadContainers(); // Refresh the list
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to start container';
		}
	}

	async function handleStopContainer(id: string) {
		try {
			await StopContainer(id);
			await loadContainers(); // Refresh the list
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to stop container';
		}
	}

	async function handleRestartContainer(id: string) {
		try {
			await RestartContainer(id);
			await loadContainers(); // Refresh the list
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to restart container';
		}
	}

	// Load containers when component mounts
	loadContainers();
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
		class="grid grid-cols-[300px_1fr] border-l border-b border-r border-latte-surface0 dark:border-mocha-surface0 overflow-y-auto"
	>
		<nav class="bg-latte-surface1 dark:bg-mocha-surface1 p-2">
			<button 
				class="w-full bg-latte-surface2 dark:bg-mocha-surface2 p-2 rounded hover:bg-latte-surface3 dark:hover:bg-mocha-surface3"
				onclick={() => loadContainers()}
				disabled={loading}
			>
				{loading ? 'Loading...' : 'Refresh Containers'}
			</button>
		</nav>

		<content class="p-2">
			{#if error}
				<div class="text-red-500 mb-4">{error}</div>
			{/if}
			
			{#if containers.length === 0}
				<div class="text-center text-gray-500">No containers found</div>
			{:else}
				<div class="grid gap-4">
					{#each containers as container}
						<div class="bg-latte-surface1 dark:bg-mocha-surface1 p-4 rounded">
							<div class="grid grid-cols-2 gap-2">
								<div class="font-bold">ID:</div>
								<div>{container.id}</div>
								
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
									disabled={container.state === 'running'}
								>
									Start
								</button>
								<button 
									class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
									onclick={() => handleStopContainer(container.id)}
									disabled={container.state !== 'running'}
								>
									Stop
								</button>
								<button 
									class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded"
									onclick={() => handleRestartContainer(container.id)}
									disabled={container.state !== 'running'}
								>
									Restart
								</button>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</content>
	</div>
</main>
