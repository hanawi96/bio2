/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				// Map Tailwind colors to CSS variables
				primary: 'var(--color-primary)',
				'primary-light': 'var(--color-primary-light)',
				success: 'var(--color-success)',
				warning: 'var(--color-warning)',
				danger: 'var(--color-danger)',
				
				bg: {
					DEFAULT: 'var(--color-bg)',
					secondary: 'var(--color-bg-secondary)',
					tertiary: 'var(--color-bg-tertiary)',
				},
				
				text: {
					DEFAULT: 'var(--color-text)',
					secondary: 'var(--color-text-secondary)',
					tertiary: 'var(--color-text-tertiary)',
				},
				
				border: 'var(--color-border)',
				separator: 'var(--color-separator)',
			},
			
			borderRadius: {
				sm: 'var(--radius-sm)',
				DEFAULT: 'var(--radius-md)',
				md: 'var(--radius-md)',
				lg: 'var(--radius-lg)',
				xl: 'var(--radius-xl)',
				full: 'var(--radius-full)',
			},
			
			spacing: {
				1: 'var(--space-1)',
				2: 'var(--space-2)',
				3: 'var(--space-3)',
				4: 'var(--space-4)',
				5: 'var(--space-5)',
				6: 'var(--space-6)',
				8: 'var(--space-8)',
			},
			
			boxShadow: {
				sm: 'var(--shadow-sm)',
				DEFAULT: 'var(--shadow-md)',
				md: 'var(--shadow-md)',
				lg: 'var(--shadow-lg)',
				card: 'var(--shadow-card)',
			},
			
			fontFamily: {
				sans: 'var(--font-sans)',
				mono: 'var(--font-mono)',
			},
			
			fontSize: {
				xs: 'var(--text-xs)',
				sm: 'var(--text-sm)',
				base: 'var(--text-base)',
			},
		},
	},
	plugins: [],
}
