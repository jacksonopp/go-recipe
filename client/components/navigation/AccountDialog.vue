<script setup lang="ts">
const emit = defineEmits(["update:isOpen"]);
const props = defineProps<{ isOpen: boolean }>();

const dialog = ref<HTMLDialogElement | null>(null);

const close = () => {
  emit("update:isOpen", false);
};

defineExpose({ dialog });

watch(
  () => props.isOpen,
  (isOpen) => {
    if (isOpen) {
      dialog.value?.showModal();
    } else {
      dialog.value?.close();
    }
  },
);

onMounted(() => {
  dialog.value?.addEventListener("click", (event) => {
    let rect = event.target?.getBoundingClientRect();
    if (
      rect.left > event.clientX ||
      rect.right < event.clientX ||
      rect.top > event.clientY ||
      rect.bottom < event.clientY
    ) {
      dialog.value?.close();
      emit("update:isOpen", false);
    }
  });
});

const { logout } = useAuth({ autoLogin: false });
</script>

<template>
  <dialog
    @close="close"
    class="absolute top-14 right-8 m-0 ml-auto w-48 rounded backdrop:bg-transparent p-2 shadow-lg backdrop:backdrop-blur-[3px]"
    ref="dialog"
  >
    <ul class="flex flex-col gap-2">
      <li class="rounded p-2 hover:bg-gray-400/40">
        <NuxtLink to="/settings" class="w-full text-left inline-block"
          >Settings
        </NuxtLink>
      </li>
      <li class="rounded p-2 hover:bg-red-400/10">
        <button class="w-full text-left text-red-800" @click="logout">
          Log Out
        </button>
      </li>
    </ul>
  </dialog>
</template>

<style scoped></style>
