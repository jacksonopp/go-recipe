<script setup lang="ts">
const props = defineProps<{
  title: string;
}>();
</script>

<template>
  <main class="max-w-screen-md mx-auto pt-8">
    <h1 class="text-center font-header font-semibold text-3xl">
      {{ props.title }}
    </h1>
    <slot />
  </main>
</template>

<style scoped></style>
