<script setup lang="ts">
import type { Ingredient } from "~/types/recipe.type";

const props = defineProps<{
  ingredients: Ingredient[];
}>();
</script>

<template>
  <ul class="grid gap-4 mt-4">
    <li v-for="ingredient in ingredients">
      <div>
        {{ ingredient.quantity }} {{ ingredient.unit }}
        {{ ingredient.name }}
      </div>
    </li>
  </ul>
</template>

<style scoped></style>
