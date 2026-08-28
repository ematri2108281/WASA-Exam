import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import api from './services/axios.js'
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'

import './assets/dashboard.css'
import './assets/main.css'

/**
 * Clear any leftover auth state from a previous session.
 * This prevents stale tokens from persisting across hard reloads.
 */
function clearSession() {
	const keys = ['token', 'username', 'userId', 'userPhoto']
	keys.forEach(k => {
		try { localStorage.removeItem(k) } catch (_) { /* ignore */ }
	})
	try {
		delete api.defaults.headers.common['Authorization']
	} catch (_) { /* ignore */ }
}

clearSession()

/**
 * Register global components so they are available in every view
 * without needing to import them manually.
 */
function registerGlobals(app) {
	app.component('ErrorMsg', ErrorMsg)
	app.component('LoadingSpinner', LoadingSpinner)
}

const app = createApp(App)
app.config.globalProperties.$axios = api
registerGlobals(app)
app.use(router)
app.mount('#app')
