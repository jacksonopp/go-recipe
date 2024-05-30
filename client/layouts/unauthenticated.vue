<script setup lang="ts">
import HomepageNavigation from "~/components/navigation/HomepageNavigation.vue";
import { useLoader } from "~/composables/useLoader";

const { loading } = useLoader();

const dialog = ref<HTMLDialogElement | null>(null);

onBeforeMount(() => {
  console.log(loading.value);
});

watch(loading, (isLoading) => {
  console.log(isLoading);
  if (isLoading) {
    dialog.value?.showModal();
  } else {
    dialog.value?.close();
  }
});
</script>

<template>
  <dialog ref="dialog">Loading</dialog>
  <HomepageNavigation />
  <slot />
</template>

<style scoped></style>
