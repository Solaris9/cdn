<script lang="ts">
    import type { Writable } from 'svelte/store';

    export let authorization: Writable<string>;

    function removetoken() {
        authorization.set(null);
        window.location.reload();
    }
</script>

<h1>Settings</h1>

<hr />
<p class="note">This will log you out and reload the page.</p>
<button on:click={removetoken}>Logout</button>

<style>
    h1 {
        margin: unset;
    }

    h1,
    p {
        color: white;
    }

    hr {
        margin: 30px 0 10px 0;
    }

    button {
        padding: 10px;
        border-radius: 10px;
        border: none;
    }

    .note {
        padding: 5px 10px 5px 10px;
        background-color: #add8e633;
        border-radius: 10px;
        border-left: 10px solid #add8e654;
    }
</style>
