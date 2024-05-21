<script setup lang="ts">
const formData = ref({
  name: "",
  email: "",
  password: "",
  passwordConfirm: "",
});

const submitForm = () => {
  console.log(formData.value);
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
            Join our Recipe Community
          </h1>
          <p class="text-lg text-gray-200 md:text-xl lg:text-2xl">
            Connect with home cooks, share your favorite recipes, and discover
            new dishes to try.
          </p>
          <form v-on:submit.prevent="submitForm">
            <div class="flex flex-col gap-4 md:flex-row">
              <div class="flex w-full flex-col">
                <label for="name" class="my-2 text-sm text-white">Name</label>
                <input
                  v-model="formData.name"
                  class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  type="text"
                />
              </div>
              <div class="flex w-full flex-col">
                <label for="name" class="my-2 text-sm text-white">Name</label>
                <input
                  v-model="formData.email"
                  class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  type="email"
                />
              </div>
            </div>
            <div class="flex flex-col gap-4 md:flex-row">
              <div class="flex w-full flex-col">
                <label for="name" class="my-2 text-sm text-white"
                  >Password</label
                >
                <input
                  v-model="formData.password"
                  class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  type="password"
                />
              </div>
              <div class="flex w-full flex-col">
                <label for="name" class="my-2 text-sm text-white"
                  >Confirm Password</label
                >
                <input
                  v-model="formData.passwordConfirm"
                  class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  type="password"
                />
              </div>
            </div>
            <button
              class="inline-flex w-full mt-4 items-center justify-center whitespace-nowrap text-sm ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 h-10 bg-white text-orange-600 hover:bg-gray-200 rounded-md py-3 px-6 font-medium"
            >
              Join Now
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
