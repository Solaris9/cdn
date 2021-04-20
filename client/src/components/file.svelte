<script lang="ts">
    import type { Writable } from 'svelte/store';
    import type { FileResult } from './files.svelte';

    export let file: FileResult;
    export let authorization: Writable<string>;

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

        fetch(`/api/files/${file}`, requestInit).then(() => document.getElementById(file).remove());
    }

    function copy(text: string) {
        const elem = document.createElement('textarea');
        elem.style.display = 'fixed';
        elem.textContent = text;

        document.body.appendChild(elem);

        elem.select();
        document.execCommand('copy');

        document.body.removeChild(elem);
        elem.remove();
    }
</script>

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
            <button class="file-option delete" on:click={async () => del(file.file_name)}>Delete</button>
            <button class="file-option" on:click={() => copy(file.cdn_url)}>Copy Link</button>
        </span>
    </div>
    <div class="file-image">
        <img src={file.spaces_url} alt={file.file_name} />
    </div>
</div>

<style>
    .file {
        background-color: #252936;

        max-width: 400px;
        border-radius: 10px;
        margin: 15px;
        flex-grow: 1;
        box-shadow: 5px 5px 5px black;
    }

    .file-info {
        padding: 0 10px 0 10px;
        display: flex;
        flex-direction: column;
    }

    .file-options {
        margin: 10px 0 10px 0;
        display: flex;
        align-items: stretch;
        justify-content: space-around;
    }

    .file-option {
        flex-grow: 1;
        border: none;
        padding: 5px;
    }

    .file-image {
        background-color: #1a1d27;

        margin: 5px;
        height: calc(100% - 120px);
        border-radius: 10px;

        display: flex;
        align-items: center;
        justify-content: center;
    }

    img {
        max-width: 100%;
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

    span {
        display: flex;
        justify-content: space-between;
    }

    p {
        color: white;
        margin: 5px;
    }
</style>
