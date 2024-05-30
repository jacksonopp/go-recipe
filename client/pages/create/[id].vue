<script setup lang="ts">
import CreateRecipe from "~/layouts/create-recipe.vue";

definePageMeta({
  layout: "authenticated",
  middleware: ["auth"],
});
</script>

<template>
  <CreateRecipe title="Add your ingredients and steps"> hi</CreateRecipe>
</template>

<style scoped></style>
