<script setup lang="ts">
import type { PropType } from "vue";

const props = defineProps({
  buttonStyle: {
    type: String as PropType<"primary" | "outline" | "danger" | "success">,
    default: "primary",
  },
  buttonSize: {
    type: String as PropType<"small" | "medium" | "large">,
    default: "medium",
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  type: {
    type: String as PropType<"button" | "submit" | "reset">,
    default: "button",
  },
});

defineEmits(["click"]);
</script>

<template>
  <button
    @click="$emit('click')"
    :type="props.type"
    :disabled="props.disabled"
    class="inline-flex w-full items-center justify-center whitespace-nowrap rounded-md"
    :class="{
      'bg-orange-500 text-white transition-colors ring-offset-background hover:bg-orange-600/90 focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50':
        props.buttonStyle === 'primary',
      'border transition-colors ring-offset-background border-input bg-background hover:bg-gray-200/20 hover:text-accent-foreground focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50':
        props.buttonStyle === 'outline',
      'px-4 h-8 font-medium text-sm': props.buttonSize === 'small',
      'px-4 py-2 h-10 font-medium text-sm': props.buttonSize === 'medium',
      'px-4 py-4 h-12 font-semibold text-md': props.buttonSize === 'large',
    }"
  >
    <slot></slot>
  </button>
</template>

<style scoped></style>
