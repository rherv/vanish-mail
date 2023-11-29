/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/app/templates/*.html"], // This is where your HTML templates / JSX files are located
  theme: {
    extend: {
      fontFamily: {
        sans: ["Iosevka Aile Iaso", "sans-serif"],
        mono: ["Iosevka Curly Iaso", "monospace"],
        serif: ["Iosevka Etoile Iaso", "serif"],
      },
    },
  },
  plugins: [],
};