<template>
  <PageContainer :tight="tight" :wide="wide">
    <div class="table-page">
      <slot name="header" />
      <slot name="toolbar" />
      <slot name="stats" />
      <SectionCard class="table-page__card" :title="sectionTitle" :subtitle="sectionSubtitle">
        <div class="table-page__table-shell">
          <slot />
        </div>
        <template v-if="$slots.footer" #footer>
          <div class="table-page__footer-shell">
            <slot name="footer" />
          </div>
        </template>
      </SectionCard>
    </div>
  </PageContainer>
</template>

<script setup>
import PageContainer from './PageContainer.vue'
import SectionCard from './SectionCard.vue'

defineProps({
  sectionTitle: {
    type: String,
    default: ''
  },
  sectionSubtitle: {
    type: String,
    default: ''
  },
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
