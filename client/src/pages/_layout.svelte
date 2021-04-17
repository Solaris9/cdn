<script lang="ts">
    import { isActive } from '@roxi/routify';
    import { onMount } from 'svelte';

    const links = [
        {
            name: 'Upload',
            url: './dashboard/upload',
        },
        {
            name: 'Dashboard',
            url: './dashboard',
        },
        {
            name: 'Github',
            url: 'http://github.com/Solaris9/cdn',
            external: true,
        },
    ];

    onMount(() => {
        console.log($isActive(''));
    });
</script>

<header>
    <a class:active={$isActive('./')} href="/">Home</a>
    <div class="other-links">
        {#each links as { external, url, name }, i (i)}
            <a
                class:active={!external && $isActive(url)}
                href={url}
                target={external ? '_blank' : undefined}
                rel={external ? 'noopener noreferrer' : undefined}>{name}</a
            >
        {/each}
    </div>
</header>

<slot />

<footer />

<style>
    header {
        background-color: #1f2833;
        padding: 20px;
        display: flex;
        justify-content: space-between;
    }

    a {
        color: white;
        text-decoration: none;
    }

    .active {
        color: blue;
        text-decoration: underline;
    }

    .other-links a:nth-child(2) {
        margin: 0 20px 0 20px;
    }
</style>
