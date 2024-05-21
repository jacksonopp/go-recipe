<script setup lang="ts"></script>

<template>
  <header class="bg-white shadow-sm">
    <div
      class="container mx-auto flex items-center justify-between px-4 py-4 md:px-6 md:py-5"
    >
      <NuxtLink to="/" class="text-2xl font-bold text-orange-500">
        Recipe Sharing
      </NuxtLink>
      <nav class="hidden items-center space-x-6 md:flex">
        <NuxtLink
          to="/demo"
          class="text-lg font-medium text-gray-800 hover:text-orange-500"
          >Demo
        </NuxtLink>
        <a
          href="#"
          class="text-lg font-medium text-gray-800 hover:text-orange-500"
        >
          Home
        </a>
        <a
          href="#"
          class="text-lg font-medium text-gray-800 hover:text-orange-500"
        >
          Recipes
        </a>
        <a
          href="#"
          class="text-lg font-medium text-gray-800 hover:text-orange-500"
        >
          Community
        </a>
        <a
          href="#"
          class="text-lg font-medium text-gray-800 hover:text-orange-500"
        >
          About
        </a>
      </nav>
      <div class="relative w-full max-w-md md:w-auto">
        <input
          placeholder="Search recipes..."
          class="w-full rounded-md border border-gray-300 py-2 pl-4 text-lg focus:border-orange-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus-visible:ring-orange-500"
          type="text"
        />
        <button
          class="absolute top-1/2 right-0 -translate-y-1/2 rounded-r-md bg-orange-500 px-4 py-2 text-lg font-medium text-white transition-colors hover:bg-orange-600 focus:outline-none focus:ring-2 focus:ring-orange-600 focus:ring-offset-2"
        >
          Search
        </button>
      </div>
    </div>
  </header>
</template>

<style scoped></style>
