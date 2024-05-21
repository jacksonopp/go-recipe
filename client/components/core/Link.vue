<script setup lang="ts">
import type { NuxtLinkProps } from "#app";

const props = defineProps<{
  to: string | NuxtLinkProps;
}>();
</script>

<template>
  <NuxtLink
    class="text-gray-800 text-md hover:text-orange-500 transition-colors"
    :to="to"
  >
    <slot />
  </NuxtLink>
</template>

<style scoped></style>
