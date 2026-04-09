<template>
  <section class="page-toolbar">
    <div class="page-toolbar__filters">
      <slot />
    </div>
    <div v-if="$slots.actions" class="page-toolbar__actions">
      <slot name="actions" />
    </div>
  </section>
</template>

<script setup>
</script>
