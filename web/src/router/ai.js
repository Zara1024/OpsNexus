import AIAssistant from '@/views/ai/AIAssistant.vue'
import AIDiagnosis from '@/views/ai/AIDiagnosis.vue'

const routes = [
  {
    path: '/ai/assistant',
    component: AIAssistant,
    meta: { sTitle: 'AI 智能运维助手', tTitle: '助手工作台' }
  },
  {
    path: '/ai/diagnosis',
    component: AIDiagnosis,
    meta: { sTitle: 'AI 智能运维助手', tTitle: '诊断分析台' }
  }
]

export default routes
