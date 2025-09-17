/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#0ea5e9',
          600: '#0284c7',
        },
        danger: '#ef4444',
        success: '#21c55d',
      },
      backgroundColor: {
        dark: '#0e1420',
        'dark-panel': '#111a2b',
        'dark-card': '#0f1626',
      },
      textColor: {
        dark: '#e6edf3',
        'dark-muted': '#8b99b0',
      },
      borderColor: {
        dark: 'rgba(255, 255, 255, 0.08)',
      },
      boxShadow: {
        dark: '0 10px 25px rgba(0, 0, 0, 0.35)',
      },
    },
  },
  plugins: [],
}
