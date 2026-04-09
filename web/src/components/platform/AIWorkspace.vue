<template>
  <PageContainer wide>
    <div class="ai-workspace">
      <aside class="ai-workspace__rail">
        <slot name="rail" />
      </aside>
      <section class="ai-workspace__main">
        <PageHeader :eyebrow="eyebrow" :title="title" :subtitle="subtitle">
          <template v-if="$slots.intro" #intro>
            <slot name="intro" />
          </template>
          <template v-if="$slots.actions" #actions>
            <slot name="actions" />
          </template>
        </PageHeader>
        <div class="ai-workspace__content">
          <slot />
        </div>
      </section>
    </div>
  </PageContainer>
</template>

<script setup>
import PageContainer from './PageContainer.vue'
import PageHeader from './PageHeader.vue'

defineProps({
  eyebrow: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    required: true
  },
  subtitle: {
    type: String,
    default: ''
  }
})
</script>
