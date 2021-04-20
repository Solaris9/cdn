<script lang="ts">
    import type { Writable } from 'svelte/store';
    export let authorization: Writable<string>;

    let token: string;
    let error: string;

    async function submit() {
        const requestInit: RequestInit = {
            method: 'post',
            body: JSON.stringify({ token }),
            headers: {
                'content-type': 'application/json',
            },
        };

        const res = await fetch(`/api/verify`, requestInit).then((res) => res.json());

        if (!res.success) return (error = res.message);
        if (error) error = null;

        authorization.set(token);
    }
</script>

<div class="container">
    <div class="login">
        <h1>Login</h1>
        <form>
            <input type="text" bind:value={token} />
            <input type="submit" on:click|preventDefault={submit} style="display: none" />
            <button type="submit" on:click|preventDefault={submit}>Submit</button>
        </form>
        {#if error}
            <p>{error}</p>
        {/if}
    </div>
</div>

<style>
    .container {
        height: 80%;
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .login {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;

        background-color: #1f2833;
        width: 400px;
        min-height: 250px;

        border-radius: 20px;
    }

    form {
        display: flex;
    }

    h1 {
        color: white;
        margin: 20px;
    }

    input,
    button {
        border: none;
        margin: 10px 0 10px 0;
        padding: 10px 15px 10px 15px;
    }

    input {
        border-radius: 20px 0 0 20px;
    }

    button {
        border-radius: 0 20px 20px 0;
    }

    p {
        color: white;
        background-color: lightcoral;
        padding: 10px 15px 10px 15px;
        border-radius: 20px;
    }
</style>
