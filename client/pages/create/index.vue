<script setup lang="ts">
useHead({
  title: "Create a new recipe",
});

definePageMeta({
  layout: "authenticated",
  middleware: ["auth"],
});
</script>

<template>Create</template>

<style scoped></style>
