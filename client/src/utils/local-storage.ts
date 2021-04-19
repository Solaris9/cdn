import { writable as internal, get } from 'svelte/store'

export default <T>(key: string, initialValue: T = null) => {
    const store = internal(initialValue);
    const { subscribe, set } = store;

    const json = typeof localStorage != 'undefined' ? localStorage.getItem(key) : null;
    if (json) set(JSON.parse(json));

    function updateStorage(key: string, value: T) {
        if (typeof (localStorage) == 'undefined') return;

        localStorage.setItem(key, JSON.stringify(value));
    }

    return {
        set(value: T) {
            updateStorage(key, value)
            set(value)
        },
        update(cb: (value: T) => T) {
            const value = cb(get(store))

            updateStorage(key, value)
            set(value)
        },
        subscribe
    }
}