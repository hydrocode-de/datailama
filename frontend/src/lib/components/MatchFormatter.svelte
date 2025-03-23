<script lang="ts">
    const {match, query}:  {match: string, query: string} = $props();

    let matchWords = $derived(
        query.split(' ')
        .filter(chunk => chunk.trim().length > 2)
        .map(word => word.toLowerCase())
        .filter((word, index, arr) => arr.indexOf(word) === index)
    );

    let highlightedText = $derived(match.split(' ')
        .map(chunk => {
            if (matchWords.includes(chunk.toLowerCase())) {
                return `<span class="datailama-match">${chunk}</span>`
            } else {
                return chunk;
            }
        }).join(' ')
    )
</script>

<style>
    :global(.datailama-match) {
        font-weight: bold;
        color: #2563eb;
    }
</style>

<div class="datailama-match-formatter">
    {@html highlightedText}
</div>

