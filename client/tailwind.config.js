import { nextui } from "@nextui-org/theme";

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ["var(--font-sans)"],
        mono: ["var(--font-geist-mono)"],
      },
      screens: {
        betterhover: { raw: "(hover: hover)" },
      },
    },
  },
  darkMode: "class",
  plugins: [
    nextui({
      // themes: {
      // light: {
      //   colors: {
      //     primary: "#FF71D7",
      //   },
      // },
      // dark: {
      //   colors: {
      //     primary: "#CC3EA4",
      //   },
      // },
      // },
    }),
  ],
};
