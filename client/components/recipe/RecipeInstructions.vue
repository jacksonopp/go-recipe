<script setup lang="ts">
import type { Instruction } from "~/types/recipe.type";

const props = defineProps<{
  instructions: Instruction[];
}>();
</script>

<template>
  <div class="grid gap-4 mt-4">
    <div
      v-for="instruction in props.instructions"
      class="grid grid-cols-[auto_1fr] items-start gap-4"
    >
      <div
        class="rounded-full bg-gray-900 text-white w-8 h-8 flex items-center justify-center dark:bg-gray-50 dark:text-gray-900"
      >
        {{ instruction.step }}
      </div>
      <div>
        {{ instruction.contents }}
      </div>
    </div>
  </div>
</template>

<style scoped></style>
