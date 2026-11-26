/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './app/**/*.{vue,js,ts}',
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.{vue,js,ts}',
    './pages/**/*.{vue,js,ts}',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
      colors: {
        brand: {
          primary:   '#1565C0',  // Biru utama
          accent:    '#05BAF1',  // Biru logo — ikon, link, badge
          highlight: '#BDD64A',  // Hijau logo — status, dekorasi
          'page-bg': '#F5F6FA',  // Background halaman
          // Tints
          'primary-tint':  '#EBF3FF',
          'primary-border':'#C7DCFF',
          'primary-light': '#F5F8FF',
          'primary-dark':  '#0D3B7A',
          'accent-tint':   '#E5F9FE',
          'accent-border': '#B3EDFD',
          'highlight-tint':   '#F8FBE9',
          'highlight-border': '#D9E89C',
          'highlight-text':   '#4D7C0F',
        },
      },
    },
  },
  plugins: [],
}
