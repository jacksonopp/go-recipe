<script setup lang="ts">
const props = defineProps({
  leftFull: {
    type: Boolean,
    default: false,
  },
  logoUrl: {
    type: String,
    default: "/",
  },
});
</script>

<template>
  <header class="bg-white shadow-sm">
    <div class="flex items-center py-4 md:mx-6 md:py-5">
      <div>
        <NuxtLink
          :to="props.logoUrl"
          class="text-2xl font-semibold text-orange-500 font-header"
        >
          Yumfeed
        </NuxtLink>
      </div>
      <div class="flex items-center justify-between gap-4 w-full">
        <div
          class="hidden items-center gap-4 ml-4 md:flex"
          :class="{ 'w-full': props.leftFull }"
        >
          <slot name="left" />
        </div>
        <div v-if="$slots.center">
          <slot name="center" />
        </div>
        <div v-if="$slots.right" class="flex items-center gap-4">
          <slot name="right" />
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped></style>
