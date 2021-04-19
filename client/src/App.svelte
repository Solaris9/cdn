<script context="module" lang="ts">
    import { writable } from 'svelte/store';
    import local from './utils/local-storage';

    export const modal = writable(null);
    export const authorization = local<string>('token');
</script>

<script lang="ts">
    import Modal, { bind } from 'svelte-simple-modal';

    import Upload from './components/upload.svelte';
    import Settings from './components/settings.svelte';
    import Files from './components/files.svelte';
    import Login from './components/login.svelte';

    const actions = [
        // {
        //     name: 'Upload',
        //     component: Upload,
        // },
        {
            name: 'Settings',
            component: Settings,
        },
        {
            name: 'Github',
            url: 'http://github.com/Solaris9/cdn',
        },
    ];

    const dark = {
        'background-color': '#1f2833',
    };

    function show(component: any) {
        modal.set(bind(component, { authorization }));
    }
</script>

<header>
    <div class="links">
        {#each actions as { component, url, name }}
            {#if url}
                <a class="item" href={url} target={'_blank'} rel={'noopener noreferrer'}>{name}</a>
            {:else}
                <button disabled={!$authorization} class="item" on:click={() => show(component)}>{name}</button>
            {/if}
        {/each}
    </div>
</header>

<Modal show={$modal} styleWindow={dark} styleCloseButton={dark}>
    {#if $authorization}
        <Files {authorization} />
    {:else}
        <Login {authorization} />
    {/if}
</Modal>

<style>
    header {
        background-color: #1f2833;
    }

    .links {
        padding: 10px;
        display: flex;
        justify-content: center;
    }

    .item {
        padding: 10px 40px 10px 40px;
    }

    button {
        border: none;
        text-decoration: none;
        background-color: unset;
    }

    button:not(:disabled) {
        color: white;
    }

    a {
        color: white;
        text-decoration: none;
    }
</style>
