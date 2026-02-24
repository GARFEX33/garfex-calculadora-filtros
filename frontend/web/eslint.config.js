import js from '@eslint/js';
import ts from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';
import sveltePlugin from 'eslint-plugin-svelte';
import svelteParser from 'svelte-eslint-parser';
import globals from 'globals';

/** @type {import('eslint').Linter.Config[]} */
export default [
	// ── Archivos ignorados ──────────────────────────────────────────────────
	{
		ignores: [
			'.svelte-kit/**',
			'build/**',
			'dist/**',
			'node_modules/**',
			'*.config.js',
			'*.config.ts'
		]
	},

	// ── Base JS recomendado ─────────────────────────────────────────────────
	js.configs.recommended,

	// ── TypeScript — archivos .ts y .js ────────────────────────────────────
	{
		files: ['**/*.ts', '**/*.js'],
		languageOptions: {
			parser: tsParser,
			parserOptions: {
				project: './tsconfig.json',
				extraFileExtensions: ['.svelte']
			},
			globals: {
				...globals.browser,
				...globals.es2022,
				...globals.node
			}
		},
		plugins: {
			'@typescript-eslint': ts
		},
		rules: {
			// Reglas TS recomendadas
			...ts.configs['recommended'].rules,
			...ts.configs['recommended-requiring-type-checking'].rules,

			// Preferencias explícitas
			'@typescript-eslint/no-explicit-any': 'error',
			'@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
			'@typescript-eslint/consistent-type-imports': ['error', { prefer: 'type-imports' }],
			'@typescript-eslint/no-floating-promises': 'error',
			'@typescript-eslint/no-misused-promises': 'error',
			'@typescript-eslint/await-thenable': 'error',

			// JS base
			'no-console': ['warn', { allow: ['warn', 'error'] }],
			'prefer-const': 'error',
			eqeqeq: ['error', 'always'],
			'no-var': 'error'
		}
	},

	// ── Svelte — archivos .svelte ───────────────────────────────────────────
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parser: svelteParser,
			parserOptions: {
				parser: tsParser,
				project: './tsconfig.json',
				extraFileExtensions: ['.svelte']
			},
			globals: {
				...globals.browser,
				...globals.es2022
			}
		},
		plugins: {
			svelte: sveltePlugin,
			'@typescript-eslint': ts
		},
		rules: {
			// Svelte recomendado
			...sveltePlugin.configs.recommended.rules,

			// TypeScript dentro de .svelte
			...ts.configs['recommended'].rules,
			'@typescript-eslint/no-explicit-any': 'error',
			'@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
			'@typescript-eslint/consistent-type-imports': ['error', { prefer: 'type-imports' }],

			// Svelte 5: no usar sintaxis legacy
			'svelte/no-reactive-reassign': 'error',
			'svelte/valid-compile': 'error',

			// JS base
			eqeqeq: ['error', 'always'],
			'no-var': 'error',
			'no-console': ['warn', { allow: ['warn', 'error'] }],
			// prefer-const OFF en .svelte: $props() destructuring es let por diseño en Svelte 5
			'prefer-const': 'off'
		},
		processor: sveltePlugin.processors['.svelte']
	}
];
