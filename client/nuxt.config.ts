// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: {
    enabled: true,

    timeline: {
      enabled: true,
    },
  },
  css: ["~/assets/css/main.css"],

  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },

  routeRules: {
    "/api/**": { proxy: "http://127.0.0.1:8080/api/**" },
  },

  modules: [
    [
      "@nuxtjs/google-fonts",
      {
        families: {
          "Montserrat+Alternates": {
            wght: [200, 400, 600, 800],
            ital: [200, 400, 600, 800],
          },
          Montserrat: {
            wght: [200, 400, 600, 800],
            ital: [200, 400, 600, 800],
          },
        },
        subsets: ["latin"],
        display: "swap",
        prefetch: false,
        preconnect: false,
        preload: false,
        download: true,
        base64: false,
      },
    ],
  ],
});