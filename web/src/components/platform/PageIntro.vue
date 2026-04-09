<template>
  <article class="page-intro">
    <div v-if="title" class="page-intro__title">{{ title }}</div>
    <p v-if="text" class="page-intro__text">{{ text }}</p>
    <div v-if="$slots.default" class="page-intro__content">
      <slot />
    </div>
  </article>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: ''
  },
  text: {
    type: String,
    default: ''
  }
})
</script>
