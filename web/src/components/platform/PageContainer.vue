<template>
  <div :class="['platform-page', { 'platform-page--tight': tight }]">
    <div :class="['platform-page__inner', { 'platform-page__inner--wide': wide }]">
      <slot />
    </div>
  </div>
</template>

<script setup>
defineProps({
  tight: {
    type: Boolean,
    default: false
  },
  wide: {
    type: Boolean,
    default: false
  }
})
</script>
