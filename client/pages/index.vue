<script setup lang="ts">
definePageMeta({
  layout: "unauthenticated",
});

import HomepageHero from "~/components/homepage/HomepageHero.vue";
import HomepageTrendingRecipes from "~/components/homepage/HomepageTrendingRecipes.vue";
import HomepageFooter from "~/components/homepage/HomepageFooter.vue";
</script>

<template>
  <HomepageHero />
  <HomepageTrendingRecipes />
  <HomepageFooter />
</template>

<style scoped></style>
