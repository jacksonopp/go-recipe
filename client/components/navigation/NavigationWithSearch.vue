<script setup lang="ts">
import Link from "~/components/core/Link.vue";
import Bell from "~/components/icons/Bell.vue";
import ChatBubble from "~/components/icons/ChatBubble.vue";
</script>

<template>
  <NavigationBase left-full>
    <template v-slot:left>
      <Link to="/home">Home</Link>
      <Link to="/explore">Explore</Link>
      <Link to="/create">Create</Link>
      <div class="relative w-full">
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
    </template>
    <template v-slot:right>
      <div class="text-gray-600 flex gap-2">
        <Bell />
        <ChatBubble />
      </div>
    </template>
  </NavigationBase>
</template>

<style scoped></style>
