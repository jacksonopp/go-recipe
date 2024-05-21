<script setup lang="ts">
definePageMeta({
  layout: "unauthenticated",
  middleware: ["before-login"],
});

import { FetchError } from "ofetch";

const fd = {
  username: { value: "", dirty: false },
  password: { value: "", dirty: false },
};

type FdKeys = keyof typeof fd;

const formData = reactive(fd);

const serverErrorMessage = ref("");

const submitForm = async () => {
  try {
    await $fetch("/api/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: formData.username.value,
        password: formData.password.value,
      }),
    });
  } catch (e) {
    const error = e as FetchError;
    switch (error.data.code) {
      default:
        serverErrorMessage.value = "something went wrong";
        break;
    }
  }
};
</script>

<template>
  <div
    class="w-full bg-gradient-to-r from-amber-600 to-orange-600 py-20 h-screen-nav md:py-32 lg:py-40"
  >
    <div class="container mx-auto px-4 md:px-6 lg:px-8">
      <div class="mx-auto max-w-sm sm:max-w-2xl">
        <div class="flex flex-col justify-center space-y-6">
          <h1 class="text-4xl font-bold text-white md:text-5xl lg:text-6xl">
            Welcome back foodie!
          </h1>
          <p
            class="inline-flex items-center text-lg text-gray-200 md:text-xl lg:text-2xl"
          >
            Ready to cook up something amazing? Bon appétit!
            <span class="pl-4 text-5xl">🍽️</span>
          </p>
          <div v-if="!!serverErrorMessage">
            <p class="text-sm text-red-100">
              {{ serverErrorMessage }}
            </p>
          </div>
          <form v-on:submit.prevent="submitForm">
            <div class="flex flex-col gap-4 md:flex-row">
              <!--name formgroup-->
              <div class="w-full">
                <label for="name" class="my-2 text-sm text-white"
                  >Username</label
                >
                <input
                  v-model="formData.username.value"
                  class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  type="text"
                />
              </div>
            </div>
            <div class="flex flex-col md:flex-row md:gap-4">
              <div class="w-full">
                <label for="name" class="my-2 text-sm text-white"
                  >Password</label
                >
                <input
                  v-model="formData.password.value"
                  class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  type="password"
                />
              </div>
            </div>
            <button
              class="mt-4 inline-flex h-10 w-full items-center justify-center whitespace-nowrap rounded-md bg-white px-6 py-3 text-sm font-medium text-orange-600 transition-colors ring-offset-background hover:bg-gray-200 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
            >
              Log In
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
