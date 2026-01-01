/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                'status-ok': '#22c55e',
                'status-warning': '#f97316',
                'status-error': '#ef4444',
            }
        },
    },
    plugins: [],
}
