<script setup lang="ts">
definePageMeta({
  layout: "authenticated",
  middleware: ["auth"],
});

const route = useRoute();

type Response = {
  id: number;
  username: string;
  recipes: {
    id: number;
    created_at: string;
    name: string;
  }[];
};

const { data, pending, error } = useFetch<Response>(
  `/api/user/${route.params.name}`,
);
</script>

<template>
  <pre>User {{ JSON.stringify(data, null, 2) }}</pre>
  <p v-if="pending">Loading...</p>
  <p v-else-if="error">An error occurred</p>
  <ul v-else>
    <li v-for="recipe in data?.recipes">
      <NuxtLink :to="`/recipe/${recipe.id}`">{{ recipe.name }} </NuxtLink>
    </li>
  </ul>
</template>

<style scoped></style>
