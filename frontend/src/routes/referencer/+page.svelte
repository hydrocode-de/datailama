<script lang="ts">
    import type { Paper } from '$lib/models/paper';
    import PaperCard from '$lib/components/PaperCard.svelte';

    let papers = $state<Paper[]>([]);
    let text = $state<string>('');
    let showPapers = $state<boolean>(true);

    function handleTextChange(event: Event) {
        const textarea = event.target as HTMLTextAreaElement;
        text = textarea.value;
    }
</script>

<div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">DataILama Referencer</h1>

    <div class="flex gap-4">
        <!-- Main content area -->
        <div class="flex-1">
            <textarea
                class="w-full h-[600px] p-4 border rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Enter your text here..."
                oninput={handleTextChange}
            ></textarea>
        </div>

        <!-- Sliding papers panel -->
        <div class="transition-all duration-300 ease-in-out {showPapers ? 'w-[400px] opacity-100' : 'w-0 opacity-0'} overflow-hidden">
            <div class="bg-white rounded-lg shadow-lg p-4 h-[600px] overflow-y-auto">
                <h2 class="text-xl font-bold mb-4">Referenced Papers</h2>
                {#if papers.length > 0}
                    {#each papers as paper}
                        <PaperCard paper={paper} query={text} />
                    {/each}
                {:else}
                    <p class="text-gray-500">No papers referenced yet.</p>
                {/if}
            </div>
        </div>
    </div>
</div>
