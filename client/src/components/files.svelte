<script lang="ts">
    interface FileResult {
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

    import type { Writable } from 'svelte/store';
    export let authorization: Writable<string>;

    async function getFiles(): Promise<FileResults> {
        const requestInit = {
            headers: { Authorization: $authorization },
        };

        return await fetch(`http://localhost:3001/api/files`, requestInit).then((res) => res.json());
    }

    function getDate(file: FileResult): string {
        const date = new Date(file.last_modified);
        return `${date.toDateString()} ${date.toLocaleTimeString()}`;
    }

    function getFileSize(size: number): string {
        if (size < 1024) {
            return `${size}B`;
        } else if (size < 1048576) {
            let num = size / 1024;
            return `${num.toFixed(2)}KiB`;
        } else if (size < 1073741824) {
            let num = size / 1048576;
            return `${num.toFixed(2)}MiB`;
        } else {
            let num = size / 1073741824;
            return `${num.toFixed(2)}GiB`;
        }
    }

    function del(file: string) {
        const requestInit: RequestInit = {
            method: 'delete',
            headers: { Authorization: $authorization },
        };

        fetch(`http://localhost:3001/api/files/${file}`, requestInit).then(() =>
            document.getElementById(file).remove()
        );
        // .catch(() => {

        // })
    }

    function copy(text: string) {
        const elem = document.createElement('textarea');
        elem.style.display = 'fixed';
        elem.textContent = text;

        document.body.appendChild(elem);

        elem.select();
        document.execCommand('copy');

        document.body.removeChild(elem);
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
                <div class="file" id={file.file_name}>
                    <div class="file-info">
                        <span>
                            <p class="file-name">{file.file_name}</p>
                            <p class="file-size">{getFileSize(file.size)}</p>
                        </span>
                        <span>
                            <p class="file-created">{getDate(file)}</p>
                        </span>
                        <span class="file-options">
                            <button class="file-option delete" on:click={async () => del(file.file_name)}>Delete</button
                            >
                            <button class="file-option" on:click={() => copy(file.cdn_url)}>Copy Link</button>
                        </span>
                    </div>
                    <div class="file-image">
                        <img src={file.spaces_url} alt={file.file_name} />
                    </div>
                </div>
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
        overflow-y: scroll;
    }

    .file {
        max-width: 300px;
        border-radius: 20px;
        margin: 15px;
        flex-grow: 4;
        box-shadow: 5px 5px 5px black;
    }

    .file-info {
        display: flex;
        flex-direction: column;

        background-color: #1a1d27;
        border: 3px solid #111218;
        border-radius: 15px 15px 0 0;
        padding: 0 10px 30px 0;
    }

    .file-image {
        display: flex;
        justify-content: center;

        margin-top: -20px;
        border-radius: 20px;
        background-color: #252936;
    }

    .file-options {
        margin-top: 10px;
        margin-left: 10px;
        display: flex;
        align-items: stretch;
        justify-content: space-around;
    }

    .file-option {
        flex-grow: 1;
        border: none;
        padding: 5px;
    }

    .file-option.delete {
        background-color: lightcoral;
    }

    .file-option:first-child {
        border-radius: 20px 0 0 20px;
    }

    .file-option:last-child {
        border-radius: 0 20px 20px 0;
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

    span {
        display: flex;
        justify-content: space-between;
    }

    p {
        color: white;
        margin: 5px;
    }
</style>
