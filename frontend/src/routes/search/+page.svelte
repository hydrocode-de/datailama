<script lang="ts">
    import type { Paper } from '$lib/models/paper';
    import TitleSearch from '$lib/components/TitleSearch.svelte';
    import PaperBodySearch from '$lib/components/PaperBodySearch.svelte';
    import PaperCard from '$lib/components/PaperCard.svelte';

    let papers = $state<Paper[]>([]);
    let query = $state<string>('Design soil Moisture');
    let activeTab = $state('title'); // 'title' or 'body'

    function handlePapers(paperResults: Paper[], count: number, query?: string) {
        papers = paperResults;
        query = query || '';
    }
    
    function setActiveTab(tab: string) {
        activeTab = tab;
        papers = []; // Clear results when switching tabs
    }
</script>

<div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">DataILama Search</h1>

    <!-- Tabs -->
    <div class="border-b border-gray-200 mb-4">
        <nav class="-mb-px flex">
            <button
                class="py-2 px-4 border-b-2 font-medium text-sm mr-8 transition-colors {activeTab === 'title' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
                onclick={() => setActiveTab('title')}
            >
                Title Search
            </button>
            <button
                class="py-2 px-4 border-b-2 font-medium text-sm transition-colors {activeTab === 'body' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
                onclick={() => setActiveTab('body')}
            >
                Full-Text Search
            </button>
        </nav>
    </div>
    
    <!-- Search components -->
    {#if activeTab === 'title'}
        <TitleSearch onPapers={handlePapers} />
    {:else}
        <PaperBodySearch onPapers={handlePapers} />
    {/if}

    {#if papers.length > 0 || true}
        <div class="rounded-lg p-6 mt-8">
            <h2 class="text-xl font-bold mb-4">Results</h2>
            <!-- <PaperCard paper={{
                id: 1,
                title: 'Long title that should be long enough to break the layout as a test',
                doi: '10.1000/182',
                author: 'John Doe',
                published: new Date('2021-01-01'),
                match: 'This is a test match which is really longer needed to design a proper layout. We can also change the camelCase to PascalCase about Soil moisure, long text.',
                journal: 'Test Journal',
                url: 'https://www.google.com',
                citations: 50,
                citations_year: 10.249593894,
                cosine_distance: 0.4829389473
            }} query={query} /> -->
            {#each papers as paper}
                <PaperCard paper={paper} query={query} />
            {/each} 
        </div>
    {/if}
</div>