<script lang="ts">
    import type { Paper } from '$lib/models/paper';

    const { 
        onPapers = (papers: Paper[], count: number) => {}
    } = $props();

    let isLoading = $state(false);
    let isDirty = $state(false);
    let prompt = $state('');
    let debounceTimer = $state<number | null>(null);

    function debouncedSearch() {
        // Clear any existing timer
        if (debounceTimer !== null) {
            clearTimeout(debounceTimer);
        }
        
        // Set isLoading immediately to show feedback
        isLoading = true;
        
        // Set a new timer
        debounceTimer = setTimeout(() => {
            search();
            debounceTimer = null;
        }, 1000); // 1 second debounce
    }

    async function search() {
        isLoading = true;

        const requestUrl = `/api/paper/search/body?prompt=${encodeURIComponent(prompt)}&limit=10`;
        console.log(requestUrl);
        const response = await fetch(requestUrl);
        const data = await response.json();
        console.log(data);

        if (data.count && data.count > 0) {
            onPapers(data.paper, data.count, prompt);
        } else {
            onPapers([], 0, prompt);
        }

        isLoading = false;
    }
</script>

<div class="bg-white rounded-lg shadow-md">
    <div class="w-full p-4 flex items-center justify-between">
        <div class="flex-1">
            <input 
                type="text" 
                name="prompt" 
                id="prompt" 
                class="w-full text-lg rounded-md border-gray-300 focus:border-gray-500" 
                bind:value={prompt}
                oninput={debouncedSearch}
                placeholder="Search articles..."
            />
        </div>
        <div class="flex items-center gap-4 ml-4">
            {#if isLoading}
                <svg class="animate-spin h-5 w-5 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
            {/if}
        </div>
    </div>
</div>