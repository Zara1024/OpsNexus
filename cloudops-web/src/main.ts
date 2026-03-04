import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'uno.css'
import './styles/reset.scss'
import './styles/variables.scss'

import App from './App.vue'
import router from './router'
import { setupPermissionDirective } from '@/directives/permission'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
setupPermissionDirective(app)

app.mount('#app')
