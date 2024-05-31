<script setup lang="ts">
import type { Comment } from "~/types/comment.type";

const props = defineProps<{
  comment: Comment;
}>();
</script>

<template>
  <div class="flex items-start gap-4">
    <div
      class="rounded-lg w-11 h-11 bg-[#55efc4] text-3xl flex items-center justify-center"
    >
      😋
    </div>
    <div class="grid gap-1 items-start text-sm">
      <div class="flex items-center gap-2">
        <div class="font-medium">{{ props.comment.name }}</div>
        <div class="text-sm text-gray-500 dark:text-gray-400">
          {{ props.comment.date }}
        </div>
      </div>
      <div>
        <p>
          {{ props.comment.content }}
        </p>
      </div>
      <div class="flex items-center gap-2 mt-1">
        <button
          class="bg-gray-100 border-2 border-gray-100 flex items-center gap-1 px-2 h-7 rounded-full hover:bg-white hover:border-gray-300 dark:bg-gray-800 dark:border-gray-800 dark:hover:bg-gray-950"
        >
          👍 {{ props.comment.likes }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
