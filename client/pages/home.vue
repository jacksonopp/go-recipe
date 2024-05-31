<script setup lang="ts">
import MasonryGrid from "~/components/homepage/MasonryGrid.vue";

definePageMeta({
  layout: "authenticated",
  middleware: ["auth"],
});
</script>

<template>
  <main class="mx-auto pt-8 bg-orange-50/70">
    <MasonryGrid />
  </main>
</template>

<style scoped></style>
