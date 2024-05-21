<script setup lang="ts"></script>

<template>
  <div class="bg-white px-4 py-12 md:px-6 lg:py-20">
    <div class="container mx-auto">
      <div class="mb-8 flex items-center justify-between">
        <h2 class="text-3xl font-bold text-gray-800">Trending Recipes</h2>
        <a href="#" class="text-lg font-medium text-orange-500 hover:underline">
          View All
        </a>
      </div>
      <div
        class="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
      >
        <div
          class="rounded-lg bg-orange-100/80 p-4 shadow-md transition-transform hover:-translate-y-1 hover:shadow-lg"
        >
          <img
            src="/placeholder.svg"
            alt="Trending Recipe 1"
            class="mb-4 h-48 w-full rounded-lg object-cover"
            width="300"
            height="200"
            style="aspect-ratio: 300 / 200; object-fit: cover"
          />
          <h3 class="text-xl font-bold text-gray-800">
            Creamy Garlic Shrimp Pasta
          </h3>
        </div>
        <div
          class="rounded-lg bg-orange-100/80 p-4 shadow-md transition-transform hover:-translate-y-1 hover:shadow-lg"
        >
          <img
            src="/placeholder.svg"
            alt="Trending Recipe 2"
            class="mb-4 h-48 w-full rounded-lg object-cover"
            width="300"
            height="200"
            style="aspect-ratio: 300 / 200; object-fit: cover"
          />
          <h3 class="text-xl font-bold text-gray-800">
            Grilled Salmon with Lemon Dill Sauce
          </h3>
        </div>
        <div
          class="rounded-lg bg-orange-100/80 p-4 shadow-md transition-transform hover:-translate-y-1 hover:shadow-lg"
        >
          <img
            src="/placeholder.svg"
            alt="Trending Recipe 3"
            class="mb-4 h-48 w-full rounded-lg object-cover"
            width="300"
            height="200"
            style="aspect-ratio: 300 / 200; object-fit: cover"
          />
          <h3 class="text-xl font-bold text-gray-800">
            Vegetarian Chili with Sweet Potatoes
          </h3>
        </div>
        <div
          class="rounded-lg bg-orange-100/80 p-4 shadow-md transition-transform hover:-translate-y-1 hover:shadow-lg"
        >
          <img
            src="/placeholder.svg"
            alt="Trending Recipe 4"
            class="mb-4 h-48 w-full rounded-lg object-cover"
            width="300"
            height="200"
            style="aspect-ratio: 300 / 200; object-fit: cover"
          />
          <h3 class="text-xl font-bold text-gray-800">
            Crispy Baked Chicken Thighs
          </h3>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
