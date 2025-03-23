<script lang="ts">
    import type { Paper } from '$lib/models/paper';
    import MatchFormatter from './MatchFormatter.svelte';
    
    // Props
    const { paper, query }: { paper: Paper, query: string } = $props();
    
    // State
    let isExpanded = $state(false);
    
    function toggleExpand() {
        isExpanded = !isExpanded;
    }
</script>

<div class="bg-white rounded-lg shadow-md p-4 mb-4 hover:shadow-lg transition-shadow">
    <!-- Match text if available -->
    {#if paper.match}
        <div class="flex flex-col text-gray-800 text-sm mb-3 p-3 bg-gray-50 rounded-md border-l-4 border-blue-400">
            <MatchFormatter match={paper.match} {query} />
        </div>
    {/if}
    
    <!-- Preview mode - Collapsed -->
    {#if !isExpanded}
        <div class="flex justify-between items-center">
            <div class="flex-1">
                <h3 class="text-lg font-semibold text-gray-900 truncate">{paper.title}</h3>
                <p class="text-sm text-gray-600">{paper.author} et al. • {new Date(paper.published).getFullYear()}</p>
            </div>
            <button 
                class="ml-2 p-2 text-gray-500 hover:bg-gray-100 rounded-full transition-colors"
                onclick={toggleExpand}
                aria-label="Expand details"
            >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                </svg>
            </button>
        </div>
    {:else}
        <!-- Full details - Expanded -->
        <div class="mb-2 flex justify-between">
            <div class="flex-1">
                <h3 class="text-lg font-semibold text-gray-900 mb-1">{paper.title}</h3>
            </div>
            <button 
                class="ml-2 p-2 text-gray-500 hover:bg-gray-100 rounded-full transition-colors"
                onclick={toggleExpand}
                aria-label="Collapse details"
            >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
                </svg>
            </button>
        </div>
        
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
                        {paper.journal}&nbsp;
                    </a>
                {/if}
            </div>

            <!-- Middle column: Author -->
            <div class="lg:w-3/5">
                <p class="text-sm text-gray-600">{paper.author} et al.</p>
            </div>

            <!-- Right column: Year and Citations -->
            <div class="lg:w-1/5 text-right">
                <div class="text-sm font-medium text-gray-500">{new Date(paper.published).getFullYear()}</div>
                <div class="text-xl font-bold text-gray-900">{paper.citations} / {paper.citations_year.toFixed(2)}</div>
                <div class="text-xs text-gray-500">citations / per year</div>
            </div>
        </div>
    {/if}
</div> 