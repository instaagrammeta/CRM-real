import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './composables/**/*.{js,ts}',
    './plugins/**/*.{js,ts}',
    './app.vue',
    './error.vue',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      colors: {
        // Donezo green palette (matches original Flask templates)
        brand: {
          50:  '#f0f9f4',
          100: '#e8f5ee',
          200: '#d4edda',
          300: '#a8d8b9',
          400: '#6db98a',
          500: '#3a9c63',
          600: '#1f7a4d',
          700: '#0f5f3a',
          800: '#0a4329',
          900: '#04321e',
        },
        ink: {
          DEFAULT: '#1a1a2e',
          soft:    '#4a5568',
          mute:    '#8b9dc3',
        },
        page: '#f1f4f8',
        border: '#eef2f6',
      },
      borderRadius: {
        'lg': '14px',
        'xl': '20px',
        '2xl': '24px',
        '3xl': '28px',
      },
      boxShadow: {
        card:  '0 1px 3px rgba(0,0,0,0.04)',
        hover: '0 8px 24px -8px rgba(0,0,0,0.12)',
        'card-lg': '0 12px 32px -12px rgba(0,0,0,0.16)',
      },
      keyframes: {
        'fade-in': {
          '0%': { opacity: '0', transform: 'translateY(4px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        'slide-in-right': {
          '0%': { opacity: '0', transform: 'translateX(40px)' },
          '100%': { opacity: '1', transform: 'translateX(0)' },
        },
        'pulse-soft': {
          '0%,100%': { opacity: '1' },
          '50%':     { opacity: '0.5' },
        },
      },
      animation: {
        'fade-in': 'fade-in 200ms ease-out',
        'slide-in-right': 'slide-in-right 280ms ease-out',
        'pulse-soft': 'pulse-soft 2s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
