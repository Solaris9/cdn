<script context="module" lang="ts">
    export interface FileResult {
        cdn_url: string;
        spaces_url: string;
        spaces_cdn: string;
        file_name: string;
        last_modified: Date;
        size: number;
    }

    interface FileResults {
        files: FileResult[];
        length: number;
    }
</script>

<script lang="ts">
    import type { Writable } from 'svelte/store';
    import File from './file.svelte';
    export let authorization: Writable<string>;

    async function getFiles(): Promise<FileResults> {
        const requestInit = {
            headers: { Authorization: $authorization },
        };
        // http://localhost:3001
        return await fetch(`/api/files`, requestInit).then((res) => res.json());
    }
</script>

{#await getFiles()}
    <div class="container">
        <div class="container-info">
            <h1>Loading...</h1>
        </div>
    </div>
{:then result}
    {#if result.length > 0}
        <div class="files">
            {#each result.files as file, i (i)}
                <File {file} {authorization} />
            {/each}
        </div>
    {:else}
        <div class="container">
            <div class="container-info">
                <h1>No images found</h1>
            </div>
        </div>
    {/if}
{:catch error}
    <div class="container">
        <div class="container-info">
            <h1>Error</h1>
            <p>{error}</p>
        </div>
    </div>
{/await}

<style>
    .files {
        max-height: calc(100vh - 60px);
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        justify-content: center;
        overflow-y: auto;
    }

    .container {
        height: 80%;
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .container-info {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }

    p {
        color: white;
        margin: 5px;
    }
</style>
