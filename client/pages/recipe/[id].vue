<script setup lang="ts">
import type { RecipeType } from "~/types/recipe.type";

definePageMeta({
  layout: "authenticated",
  // middleware: ["auth"],
});

const route = useRoute();

const { data: recipe } = useFetch<RecipeType>(`/api/recipe/${route.params.id}`);
</script>

<template>
  <pre>{{ JSON.stringify(data, null, 2) }}</pre>
  <div class="flex flex-col min-h-[100dvh]">
    <section class="w-full py-12 md:py-24 lg:py-32 bg-orange-600">
      <div class="container px-4 md:px-6">
        <div class="grid gap-8 lg:grid-cols-2 items-center">
          <div class="space-y-4 order-2 lg:order-1">
            <h1
              class="font-header text-white text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl"
            >
              {{ recipe?.name }}
            </h1>
            <p
              class="max-w-[600px] md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed"
            >
              {{ recipe?.description }}
            </p>
          </div>
          <img
            src="/placeholder.svg"
            width="800"
            height="500"
            alt="Spaghetti Bolognese"
            class="aspect-video rounded-xl object-cover order-1 lg:order-2"
          />
        </div>
      </div>
    </section>
    <section class="w-full py-12 md:py-24 lg:py-32">
      <div class="container px-4 md:px-6">
        <div class="grid gap-8">
          <div>
            <h2 class="text-2xl font-bold tracking-tighter sm:text-3xl">
              Ingredients
            </h2>
            <RecipeIngredients :ingredients="recipe!.ingredients" />
          </div>
          <div>
            <h2 class="text-2xl font-bold tracking-tighter sm:text-3xl">
              Instructions
            </h2>
            <RecipeInstructions :instructions="recipe!.instructions" />
          </div>
          <RecipeComments />
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped></style>
