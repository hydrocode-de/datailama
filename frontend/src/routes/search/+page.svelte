<script lang="ts">
    import type { Paper } from '$lib/models/paper';
    
    let title = $state('');
    let author = $state('');
    let order = $state('citations_year');
    let direction = $state('desc');
    let isLoading = $state(false);
    let isExpanded = $state(false);

    let papers = $state<Paper[]>([]);

    async function search() {
        isLoading = true;
        const response = await fetch(`/api/paper/search?title=${title}&author=${author}&order=${order}&direction=${direction}&limit=10`);
        const data = await response.json();
        console.log(data);
        if (data.count && data.count > 0) {
            papers = data.paper;
        } else {
            papers = [];
        }
        isLoading = false;
    }
</script>

<div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">DataILama Search</h1>

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

    {#if papers.length > 0}
        <div class="rounded-lg p-6 mt-8">
            <h2 class="text-xl font-bold mb-4">Results</h2>
                {#each papers as paper}
                    <div class="bg-white rounded-lg shadow-md p-4 mb-4 hover:shadow-lg transition-shadow">
                        <div class="flex flex-col lg:flex-row gap-4">
                            <!-- Left column: DOI and Journal link -->
                            <div class="lg:w-1/5 flex flex-col">
                                {#if paper.doi}
                                    <a href="https://doi.org/{paper.doi}" target="_blank" rel="noopener noreferrer" 
                                       class="inline-block bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full mb-2">
                                        DOI: {paper.doi}
                                    </a>
                                {/if}
                                {#if paper.journal}
                                    <a href={paper.url} target="_blank" rel="noopener noreferrer" 
                                       class="text-sm text-gray-600 hover:text-blue-600">
                                        {paper.journal}
                                    </a>
                                {/if}
                            </div>

                            <!-- Middle column: Title and Author -->
                            <div class="lg:w-3/5">
                                <h3 class="text-lg font-semibold text-gray-900 mb-1">{paper.title}</h3>
                                <p class="text-sm text-gray-600">{paper.author}</p>
                            </div>

                            <!-- Right column: Year and Citations -->
                            <div class="lg:w-1/5 text-right">
                                <div class="text-sm font-medium text-gray-500">{paper.published}</div>
                                <div class="text-xl font-bold text-gray-900">{paper.citations} / {paper.citations_year.toFixed(2)}</div>
                                <div class="text-xs text-gray-500">citations / per year</div>
                            </div>
                        </div>
                    </div>
                {/each} 
        </div>
    {/if}
</div>