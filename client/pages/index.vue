<script setup lang="ts">
definePageMeta({
  layout: "unauthenticated",
});

import HomepageHero from "~/components/landing-page/HomepageHero.vue";
import HomepageTrendingRecipes from "~/components/landing-page/HomepageTrendingRecipes.vue";
import HomepageFooter from "~/components/landing-page/HomepageFooter.vue";
</script>

<template>
  <HomepageHero />
  <HomepageTrendingRecipes />
  <HomepageFooter />
</template>

<style scoped></style>
