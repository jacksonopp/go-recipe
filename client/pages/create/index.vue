<script setup lang="ts">
import CreateRecipe from "~/layouts/create-recipe.vue";

useHead({
  title: "Create a new recipe",
});

definePageMeta({
  layout: "authenticated",
  middleware: ["auth"],
});
</script>

<template>
  <main class="max-w-screen-md mx-auto pt-8">
    <CreateRecipe title="Create a new recipe">
      <form class="w-1/2 mx-auto mt-8" @submit.prevent>
        <div class="flex flex-col gap-4 md:flex-row">
          <div class="w-full">
            <label for="title" class="my-2 text-sm">Title</label>
            <input
              id="title"
              class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              type="text"
            />
          </div>
        </div>
        <button
          class="w-full py-3 mt-4 text-white bg-gradient-to-r from-amber-600 to-orange-600 rounded-md"
          type="submit"
        >
          Create recipe
        </button>
      </form>
    </CreateRecipe>
  </main>
</template>

<style scoped></style>
