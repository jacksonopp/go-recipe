<script setup lang="ts">

</script>

<template>
  <div class="bg-orange-100/80 px-4 py-12 md:px-6 lg:py-20">
    <div class="container mx-auto grid grid-cols-1 gap-8 lg:grid-cols-3">
      <div class="col-span-2">
        <div class="mb-8">
          <h2 class="text-3xl font-bold text-gray-800">Search Recipes</h2>
          <p class="text-lg text-gray-600">Find the perfect recipe for any occasion or craving.</p>
        </div>
        <div class="relative">
          <input
              placeholder="Search for recipes..."
              class="w-full rounded-md border border-gray-200 px-4 py-3 text-lg focus:border-orange-500 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"
              type="text"
          />
          <button class="absolute right-2 top-1/2 -translate-y-1/2 rounded-md bg-orange-500 px-4 py-2 text-lg font-medium text-white shadow-md transition-colors hover:bg-orange-600 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2">
            Search
          </button>
        </div>
      </div>
      <div class="rounded-lg bg-white p-6 shadow-md">
        <h2 class="mb-4 text-2xl font-bold text-gray-800">Cuisine Categories</h2>
        <ul class="space-y-2">
          <li>
            <a href="#" class="text-lg text-orange-500 hover:underline">
              American
            </a>
          </li>
          <li>
            <a href="#" class="text-lg text-orange-500 hover:underline">
              Italian
            </a>
          </li>
          <li>
            <a href="#" class="text-lg text-orange-500 hover:underline">
              Mexican
            </a>
          </li>
          <li>
            <a href="#" class="text-lg text-orange-500 hover:underline">
              Asian
            </a>
          </li>
          <li>
            <a href="#" class="text-lg text-orange-500 hover:underline">
              Mediterranean
            </a>
          </li>
          <li>
            <a href="#" class="text-lg text-orange-500 hover:underline">
              Vegetarian
            </a>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>