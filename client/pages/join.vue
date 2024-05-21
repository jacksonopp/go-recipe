<script setup lang="ts">
definePageMeta({
  layout: "unauthenticated",
  middleware: ["before-login"],
});

import { FetchError } from "ofetch";

const router = useRouter();

const fd = {
  username: { value: "", dirty: false },
  // email: { value: "", dirty: false },
  password: { value: "", dirty: false },
  passwordConfirm: { value: "", dirty: false },
};

const serverErrorMessage = ref("");

type FdKeys = keyof typeof fd;

const formData = reactive(fd);

const markDirty = (field: FdKeys) => {
  formData[field].dirty = true;
};

const isUsernameValid = computed(() => {
  const pattern = /^[a-zA-Z0-9]{3,20}$/;

  // Test the username against the pattern
  return pattern.test(formData.username.value);
});

const isPasswordLongEnough = computed(() => {
  return formData.password.value.length >= 8;
});

const doPasswordsMatch = computed(() => {
  return formData.password.value === formData.passwordConfirm.value;
});

const buttonDisabled = computed(() => {
  return (
    !isPasswordLongEnough.value ||
    !doPasswordsMatch.value ||
    !isUsernameValid.value
  );
});

const submitForm = async () => {
  try {
    await $fetch("/api/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: formData.username.value,
        password: formData.password.value,
        passwordConfirm: formData.passwordConfirm.value,
      }),
    });

    router.push("/login");
  } catch (e: any) {
    const err: FetchError = e;
    console.log(err.data);
    switch (err.data.code) {
      case 409:
        serverErrorMessage.value = `Username "${formData.username.value}" already exists`;
        break;
      default:
        serverErrorMessage.value = "An error occurred";
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
          <h1
            class="text-4xl font-header font-semibold text-white md:text-5xl lg:text-6xl"
          >
            Join our Recipe Community
          </h1>
          <p class="text-lg font-header text-gray-200 md:text-xl lg:text-2xl">
            Connect with home cooks, share your favorite recipes, and discover
            new dishes to try.
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
                  @input="() => markDirty('username')"
                  v-model="formData.username.value"
                  class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  type="text"
                />
                <p
                  class="text-sm text-red-100"
                  v-if="formData.username.dirty && !isUsernameValid"
                >
                  Username can only contain numbers and letters. It must be
                  between 3 and 20 characters.
                </p>
              </div>
            </div>
            <div class="flex flex-col md:flex-row md:gap-4">
              <div class="w-full">
                <label for="name" class="my-2 text-sm text-white"
                  >Password</label
                >
                <input
                  @input="() => markDirty('password')"
                  v-model="formData.password.value"
                  class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  type="password"
                />
                <p
                  class="text-sm text-red-100"
                  v-if="formData.password.dirty && !isPasswordLongEnough"
                >
                  Password must be at least 8 characters long.
                </p>
              </div>
              <div class="w-full">
                <label for="name" class="my-2 text-sm text-white"
                  >Confirm Password</label
                >
                <input
                  @input="() => markDirty('passwordConfirm')"
                  v-model="formData.passwordConfirm.value"
                  class="flex h-10 w-full flex-1 rounded-md file:border-0 border file:bg-transparent bg-white/20 px-4 py-3 text-sm file:text-sm file:font-medium text-white placeholder:text-gray-200 border-input ring-offset-background focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  type="password"
                />
                <p
                  class="text-sm text-red-100"
                  v-if="formData.passwordConfirm.dirty && !doPasswordsMatch"
                >
                  Passwords do not match.
                </p>
              </div>
            </div>
            <button
              :disabled="buttonDisabled"
              class="mt-4 inline-flex h-10 w-full items-center justify-center whitespace-nowrap rounded-md bg-white px-6 py-3 text-sm font-medium text-orange-600 transition-colors ring-offset-background hover:bg-gray-200 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
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
