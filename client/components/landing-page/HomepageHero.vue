<script setup lang="ts"></script>

<template>
  <div class="bg-orange-100/80 px-4 py-12 md:px-6 lg:py-20">
    <div class="container mx-auto grid grid-cols-1 gap-8 lg:grid-cols-2">
      <div class="flex flex-col items-start justify-center space-y-6">
        <h1
          class="text-4xl font-bold tracking-tight text-gray-800 sm:text-5xl lg:text-6xl"
        >
          Discover and Share Delicious Recipes
        </h1>
        <p class="text-lg text-gray-600">
          Join our vibrant community of food enthusiasts and explore a world of
          culinary inspiration. Share your favorite recipes, connect with
          like-minded foodies, and discover new flavors to delight your taste
          buds.
        </p>
        <div class="flex items-center space-x-4">
          <NuxtLink
            to="/join"
            class="rounded-md bg-orange-500 px-6 py-3 text-lg font-medium text-white shadow-md transition-colors hover:bg-orange-600 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"
          >
            Join Now
          </NuxtLink>
          <a
            href="#"
            class="text-lg font-medium text-orange-500 hover:underline"
          >
            Explore Recipes
          </a>
        </div>
      </div>
      <div class="relative">
        <img
          src="/placeholder.svg"
          alt="Hero RecipeType"
          class="mx-auto h-full w-full rounded-lg object-cover shadow-lg max-w-[640px]"
          width="640"
          height="480"
          style="aspect-ratio: 640 / 480; object-fit: cover"
        />
        <div
          class="absolute bottom-0 left-0 w-full rounded-b-lg bg-white/90 px-6 py-4 backdrop-blur-sm"
        >
          <div class="flex items-center justify-between">
            <h2 class="text-2xl font-bold text-gray-800">
              Creamy Mushroom Risotto
            </h2>
            <button
              class="rounded-md bg-orange-500 px-4 py-2 text-sm font-medium text-white shadow-md transition-colors hover:bg-orange-600 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"
            >
              Save RecipeType
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
