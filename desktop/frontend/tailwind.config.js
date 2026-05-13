/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        tactical: {
          950: '#07080d',
          900: '#0e111a',
          800: '#171b28',
          700: '#242a3a',
          red: '#ff4655',
          cyan: '#49f5d4',
        },
      },
      fontFamily: {
        sans: ['Inter', 'Segoe UI', 'Arial', 'sans-serif'],
      },
    },
  },
  plugins: [],
};
