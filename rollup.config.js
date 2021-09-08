import babel from '@rollup/plugin-babel';
import commonjs from '@rollup/plugin-commonjs';
import css from 'rollup-plugin-css-porter';
import json from '@rollup/plugin-json';
import resolve from '@rollup/plugin-node-resolve';

export default [
    {
        input: 'src/index.js',        
        output: { 
            file: 'public/main.js',
            format: 'iife',
            sourcemap: true,
            name: 'Viewer',
            globals: {
                'leaflet': 'L',
            },
        },
        plugins: [                      
            resolve({moduleDirectories: ['node_modules', 'src']}),
            commonjs(),
            json(),
            css({dest: 'public/main.css', minified: false}),
            babel({    
                babelHelpers: 'bundled',
                extensions: ['.js', '.mjs'],
                exclude: ['node_modules/@babel/**', 'node_modules/core-js/**'],
                include: ['src/**', 'node_modules/**'],
            }),
        ],
    },
    {
        input: 'node_modules/@liaozhai/tiles/src/tile.js',        
        output: { 
            file: 'public/tile.js',
            format: 'es',
            sourcemap: true,
        },
        plugins: [                      
            resolve({moduleDirectories: ['node_modules', 'src']}),
            commonjs(),
            json(),
            babel({    
                babelHelpers: 'bundled',
                extensions: ['.js', '.mjs'],
                exclude: ['node_modules/@babel/**', 'node_modules/core-js/**'],
                include: ['src/**', 'node_modules/**'],
            }),
        ],
    },     
];