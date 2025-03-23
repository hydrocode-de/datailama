<script lang="ts">
    import type { Paper } from '$lib/models/paper';
    
    // Props with callbacks
    const { 
        initialTitle = '', 
        initialAuthor = '', 
        initialOrder = 'citations_year', 
        initialDirection = 'desc',
        onPapers = (papers: Paper[], count: number) => {} 
    } = $props();
    
    // State
    let title = $state(initialTitle);
    let author = $state(initialAuthor);
    let order = $state(initialOrder);
    let direction = $state(initialDirection);
    let isLoading = $state(false);
    let isExpanded = $state(false);
    
    async function search() {
        isLoading = true;
        const requestUrl = `/api/paper/search/title?title=${encodeURIComponent(title)}&author=${encodeURIComponent(author)}&order=${encodeURIComponent(order)}&direction=${encodeURIComponent(direction)}&limit=10`;
        console.log(requestUrl);
        const response = await fetch(requestUrl);
        const data = await response.json();
        console.log(data);
        
        // Call the callback with search results
        if (data.count && data.count > 0) {
            onPapers(data.paper, data.count);
        } else {
            onPapers([], 0);
        }
        
        isLoading = false;
    }
</script>

<div class="bg-white rounded-lg shadow-md">
    <!-- Basic Search Header -->
    <div class="w-full p-4 flex items-center justify-between">
        <div class="flex-1">
            <input 
                type="text" 
                name="title" 
                id="title" 
                class="w-full text-lg rounded-md border-gray-300 focus:border-gray-500" 
                bind:value={title}
                oninput={search}
                placeholder="Search by title..."
            />
        </div>
        <div class="flex items-center gap-4 ml-4">
            {#if isLoading}
                <svg class="animate-spin h-5 w-5 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
            {/if}
            <button 
                class="p-2 hover:bg-gray-50 rounded-full transition-colors"
                onclick={() => isExpanded = !isExpanded}
                aria-label={isExpanded ? "Collapse advanced search" : "Expand advanced search"}
            >
                <svg 
                    class="w-6 h-6 text-gray-500 transform transition-transform {isExpanded ? 'rotate-180' : ''}" 
                    xmlns="http://www.w3.org/2000/svg" 
                    viewBox="0 0 20 20" 
                    fill="currentColor"
                >
                    <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
            </button>
        </div>
    </div>

    <!-- Advanced Search Content -->
    {#if isExpanded}
        <div class="border-t border-gray-200 p-4">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                    <label for="author" class="block text-sm font-medium text-gray-700">Author</label>
                    <input 
                        type="text" 
                        name="author" 
                        id="author" 
                        class="mt-1 block w-full text-lg rounded-md border-gray-300 focus:border-indigo-500 focus:ring-indigo-500" 
                        bind:value={author}
                        oninput={search}
                    />
                </div>
                <div>
                    <label for="order" class="block text-sm font-medium text-gray-700">Order by</label>
                    <select 
                        name="order" 
                        id="order" 
                        class="mt-1 block w-full text-lg rounded-md border-gray-300 focus:border-indigo-500 focus:ring-indigo-500"
                        bind:value={order}
                        onchange={search}
                    >
                        <option value="citations_year">Citations per Year</option>
                        <option value="citations">Total Citations</option>
                    </select>
                </div>
                <div>
                    <label for="direction" class="block text-sm font-medium text-gray-700">Direction</label>
                    <select 
                        name="direction" 
                        id="direction" 
                        class="mt-1 block w-full text-lg rounded-md border-gray-300 focus:border-indigo-500 focus:ring-indigo-500"
                        bind:value={direction}
                        onchange={search}
                    >
                        <option value="desc">Descending</option>
                        <option value="asc">Ascending</option>
                    </select>
                </div>
            </div>
        </div>
    {/if}
</div> 