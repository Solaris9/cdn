import svelte from 'rollup-plugin-svelte';
import commonjs from '@rollup/plugin-commonjs';
import dev from 'rollup-plugin-dev';
import resolve from '@rollup/plugin-node-resolve';
import livereload from 'rollup-plugin-livereload';
import { terser } from 'rollup-plugin-terser';
import sveltePreprocess from 'svelte-preprocess';
import typescript from '@rollup/plugin-typescript';
import css from 'rollup-plugin-css-only';

const production = !process.env.ROLLUP_WATCH;

function serve() {
    let server;

    function toExit() {
        if (server) server.kill(0);
    }

    return {
        writeBundle() {
            if (server) return;
            server = require('child_process').spawn('npm', ['run', 'start', '--', '--dev'], {
                stdio: ['ignore', 'inherit', 'inherit'],
                shell: true,
            });

            process.on('SIGTERM', toExit);
            process.on('exit', toExit);
        },
    };
}

export default {
    input: 'src/main.ts',
    output: {
        sourcemap: true,
        format: 'iife',
        name: 'app',
        file: 'public/build/bundle.js',
        inlineDynamicImports: true,
    },
    plugins: [
        svelte({
            preprocess: sveltePreprocess({ sourceMap: !production }),
            compilerOptions: {
                dev: !production,
            },
        }),
        css({ output: 'bundle.css' }),
        resolve({
            browser: true,
            dedupe: ['svelte', 'svelte/transition', 'svelte/internal'],
        }),
        commonjs(),
        typescript({
            rootDir: './src',
            sourceMap: production,
            inlineSources: !production,
        }),
        !production && serve(),
        !production && livereload('public'),
        production && terser(),
        !production &&
            dev({
                dirs: ['public'],
                port: 3000,
                proxy: { '/api/*': 'http://localhost:3000/' },
                spa: true,
            }),
    ],
    watch: {
        clearScreen: false,
    },
};
