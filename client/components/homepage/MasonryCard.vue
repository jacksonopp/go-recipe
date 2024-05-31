<script setup lang="ts">
const props = defineProps<{
  image: string;
  name: string;
  title?: string;
  description?: string;
  alt?: string;
  recipeId: number;
}>();
</script>

<template>
  <div class="grid gap-4 bg-white shadow rounded-xl">
    <div>
      <img
        class="h-auto max-w-full rounded-t-xl"
        :src="props.image"
        :alt="props.alt ?? ''"
      />
    </div>
    <div class="px-4 pb-4">
      <NuxtLink
        :to="`/recipe/${props.recipeId}`"
        class="font-header font-semibold block hover:text-orange-800"
        >{{ props.title }}
      </NuxtLink>
      <NuxtLink
        :to="`/user/${encodeURI(props.name)}`"
        class="text-sm text-gray-500 hover:text-orange-800 block"
        >{{ props.name }}
      </NuxtLink>
    </div>
  </div>
</template>

<style scoped></style>
