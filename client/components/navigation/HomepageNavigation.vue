<script setup lang="ts">
import NavigationBase from "~/components/navigation/NavigationBase.vue";
</script>

<template>
  <NavigationBase>
    <template v-slot:left>
      <a
        href="#"
        class="text-gray-800 text-md hover:text-orange-500 transition-colors"
      >
        Recipes
      </a>
      <a
        href="#"
        class="text-gray-800 text-md hover:text-orange-500 transition-colors"
      >
        Community
      </a>
    </template>
    <template v-slot:right>
      <NuxtLink
        class="rounded-full bg-orange-500 px-3 py-2 font-semibold text-white hover:bg-orange-600 transition-colors"
        to="/login"
        >Sign In
      </NuxtLink>
      <NuxtLink
        class="rounded-full bg-gray-200 px-3 py-2 font-semibold text-orange-700 hover:bg-gray-300 transition-colors"
        to="/join"
        >Sign Up
      </NuxtLink>
    </template>
  </NavigationBase>
</template>

<style scoped></style>
