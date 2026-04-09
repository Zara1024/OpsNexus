<template>
  <section :class="['section-card', { 'section-card--dense': dense }]">
    <header v-if="title || subtitle || $slots.header || $slots.actions" class="section-card__header">
      <div class="section-card__title-wrap">
        <slot name="header">
          <div v-if="title" class="section-card__title">{{ title }}</div>
          <div v-if="subtitle" class="section-card__subtitle">{{ subtitle }}</div>
        </slot>
      </div>
      <div v-if="$slots.actions" class="section-card__actions">
        <slot name="actions" />
      </div>
    </header>
    <div class="section-card__body">
      <slot />
    </div>
    <footer v-if="$slots.footer" class="section-card__footer">
      <slot name="footer" />
    </footer>
  </section>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  dense: {
    type: Boolean,
    default: false
  }
})
</script>
