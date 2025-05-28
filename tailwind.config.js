/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		"./views/**/*.templ", // Incluye todos los archivos .templ
		"./views/**/*.go", // Si usas clases en archivos Go
	],
	theme: {
		extend: {},
		fontFamily: {
			roboto: ["Roboto", "sans-serif"],
			// bangers: ["Bangers", "sans-serif"],
		},
	},
	plugins: [],
};
