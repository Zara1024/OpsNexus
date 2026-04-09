<template>
  <section class="stat-strip">
    <article
      v-for="item in items"
      :key="item.key || item.label"
      class="stat-card"
      :data-tone="item.tone || 'neutral'"
    >
      <div class="stat-card__label">{{ item.label }}</div>
      <div class="stat-card__value">
        <span v-if="item.prefix" class="stat-card__prefix">{{ item.prefix }}</span>
        {{ item.value }}
        <span v-if="item.unit" class="stat-card__unit">{{ item.unit }}</span>
      </div>
      <div v-if="item.hint || item.meta" class="stat-card__hint">
        {{ item.hint || item.meta }}
      </div>
    </article>
  </section>
</template>

<script setup>
defineProps({
  items: {
    type: Array,
    default: () => []
  }
})
</script>
