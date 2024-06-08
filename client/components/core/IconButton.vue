<script setup lang="ts">
import type { PropType } from "vue";

type ButtonStyle = "primary" | "danger";

const props = defineProps({
  buttonStyle: {
    type: String as PropType<ButtonStyle>,
    default: "primary",
  },
  type: {
    type: String as PropType<"button" | "submit" | "reset">,
    default: "button",
  },
  altText: {
    type: String,
    default: "",
    required: true,
  },
});

defineEmits(["click"]);
</script>

<template>
  <button
    :class="{
      'text-gray-800 hover:text-gray-800/80 size-[40px] transition-colors grid place-items-center focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2':
        props.buttonStyle === 'primary',
      'text-red-500 hover:text-red-500/80 size-[40px] transition-colors grid place-items-center focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2':
        props.buttonStyle === 'danger',
    }"
    :type="type"
    @click="$emit('click')"
  >
    <span class="sr-only">{{ props.altText }}</span>
    <slot />
  </button>
</template>

<style scoped></style>
