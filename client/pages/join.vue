<script setup lang="ts">
const fd = {
  username: { value: "", dirty: false },
  // email: { value: "", dirty: false },
  password: { value: "", dirty: false },
  passwordConfirm: { value: "", dirty: false },
};

type FdKeys = keyof typeof fd;

const formData = reactive(fd);

const markDirty = (field: FdKeys) => {
  console.log("###", field);
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

const submitForm = () => {};
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
            <div class="flex flex-col md:gap-4 md:flex-row">
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
