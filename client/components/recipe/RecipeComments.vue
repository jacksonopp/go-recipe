<script setup lang="ts">
import type { Comment } from "~/types/comment.type";

const comments: Comment[] = new Array(3).fill(null).map(() => ({
  name: "Sarah",
  date: "2 days ago",
  content:
    "This Bolognese recipe is absolutely delicious! The flavors are so rich and comforting. I'll definitely be making this again.",
  likes: 5,
}));
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold tracking-tighter sm:text-3xl">Comments</h2>
    <div class="grid gap-4 mt-4">
      <RecipeComment v-for="comment in comments" :comment="comment" />
    </div>
  </div>
</template>

<style scoped></style>
