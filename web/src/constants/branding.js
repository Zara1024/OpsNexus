import brandLogo from '@/assets/image/opsnexus-cyber.svg'
import loginLogo from '@/assets/image/opsnexus-cyber.svg'

export const BRANDING = Object.freeze({
  name: 'OpsNexus',
  slogan: 'AI 智能运维助手',
  description: '统一资产、容器、告警、工单与 AI 智能运维助手工作台',
  logo: brandLogo,
  loginLogo,
  loginFooterYear: '2026',
  sidebarBadge: 'AIOps Multi-Model Workspace',
  assistantPositioning: '支持接入 OpenAI / Claude / Gemini / DeepSeek / Qwen / Ollama / 本地模型，统一承载 Agent、工具编排、知识检索与诊断巡检。',
  aiModels: ['OpenAI', 'Claude', 'Gemini', 'DeepSeek', 'Qwen', 'Ollama', 'Local LLM'],
  aiCapabilities: ['Agent 协作', '工具编排', '知识检索', '诊断分析', '巡检辅助'],
  loginDescription: '统一纳管 CMDB、主机、应用发布、Kubernetes、监控告警、工单审计与 AI 助手能力，构建面向生产环境的多模型智能运维中枢。',
  loginHighlights: [
    {
      name: '统一资产与 CMDB',
      description: '纳管主机、应用、容器与数据库资源，形成可检索、可关联的运维底座。',
      points: ['CMDB / 主机 / 数据库', '资源拓扑与状态视图', '检索关联与统一纳管']
    },
    {
      name: 'Kubernetes 与发布',
      description: '覆盖集群、节点、命名空间与工作负载，支撑环境治理和发布协同。',
      points: ['集群 / 工作负载 / 网络存储', '环境治理与变更联动', '发布协同与风险收敛']
    },
    {
      name: '监控告警与自动化',
      description: '联动指标、事件、SSL 与通知渠道，支持告警收敛、推送与动作闭环。',
      points: ['指标 / 事件 / SSL 联动', '告警收敛与通知机器人', '自动化处置与闭环追踪']
    },
    {
      name: 'AI 巡检与平台协同',
      description: '融合多模型接入、Agent 协作、知识检索、巡检辅助与诊断分析，帮助快速定位问题并沉淀经验。',
      points: ['OpenAI / Claude / DeepSeek / Qwen', 'Agent / 工具编排 / 知识检索', '诊断分析 / 巡检辅助 / 经验复用']
    }
  ]
})
