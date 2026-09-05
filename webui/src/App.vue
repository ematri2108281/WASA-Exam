<template>
  <div class="app-wrapper">
    <!-- Unified Top Bar (only when NOT on login page) -->
    <header v-if="$route.path !== '/login'" class="app-topbar">
      <div class="header-left">
        <h1 class="neon-title">WASA</h1>
        <span class="subtitle">Chats</span>
      </div>
      
      <div class="header-actions">

  <!-- LEFT -->
  <div class="actions-left">

    <!-- Home -->
    <RouterLink to="/home" class="btn-icon" title="Home">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
        <polyline points="9 22 9 12 15 12 15 22"/>
      </svg>
    </RouterLink>

    <!-- New Chat -->
    <button
      class="btn-icon"
      title="New Chat"
      @click="emitNewChat"
    >
      <!-- chat bubble -->
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
      </svg>
    </button>

    <!-- New Group -->
    <RouterLink to="/groups/create" class="btn-icon" title="New Group">
      <!-- users + plus -->
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="9" cy="7" r="4"/>
        <path d="M3 21v-2a4 4 0 0 1 4-4h4"/>
        <path d="M16 11v6"/>
        <path d="M13 14h6"/>
      </svg>
    </RouterLink>

  </div>

  <!-- RIGHT -->
  <div class="actions-right">
     <!-- Profile: real user avatar -->
    <RouterLink to="/profile" class="btn-icon btn-icon-avatar" title="Profile">
      <img :src="userAvatarSrc" alt="Profile" class="topbar-avatar" />
    </RouterLink>


    <!-- Users (ICONA ORIGINALE TENUTA) -->
    <RouterLink to="/users" class="btn-icon" title="Users">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
        <circle cx="9" cy="7" r="4"/>
        <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
        <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
      </svg>
    </RouterLink>
    <!-- Logout -->
    <button @click="logout" class="btn-logout" title="Logout">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
        <polyline points="16 17 21 12 16 7"/>
        <line x1="21" y1="12" x2="9" y2="12"/>
      </svg>
    </button>

  </div>
</div>


    </header>

    <!-- Main Content -->
    <main :class="{ 'with-topbar': $route.path !== '/login', 'fullscreen': $route.path === '/login' }">
      <RouterView @new-chat="handleNewChat" @refresh="handleRefresh" />
    </main>
  </div>
</template>

<script setup>
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { onMounted, onUnmounted, ref, computed } from 'vue'
import axios from './services/axios.js'

const router = useRouter()

// ── User avatar for the top-bar ─────────────────────────────────────────────
const storedPhoto    = ref(localStorage.getItem('userPhoto') || '')
const storedUsername = ref(localStorage.getItem('username') || '')

function buildLetterAvatar(name, size = 32) {
  try {
    const letter = (String(name || '').trim().charAt(0) || '?').toUpperCase()
    const svg = `<svg xmlns='http://www.w3.org/2000/svg' width='${size}' height='${size}' viewBox='0 0 ${size} ${size}'>
  <rect width='100%' height='100%' rx='${Math.floor(size/2)}' ry='${Math.floor(size/2)}' fill='#e8f0fe'/>
  <text x='50%' y='53%' dominant-baseline='middle' text-anchor='middle' fill='var(--neon-cyan)' font-family='Segoe UI, Roboto, sans-serif' font-weight='700' font-size='${Math.floor(size*0.5)}'>${letter}</text>
</svg>`
    return 'data:image/svg+xml;utf8,' + encodeURIComponent(svg)
  } catch { return '' }
}

const userAvatarSrc = computed(() => {
  const p = storedPhoto.value || ''
  if (p && (p.startsWith('data:') || p.startsWith('http'))) return p
  return buildLetterAvatar(storedUsername.value, 32)
})

// Keep avatar in sync when ProfileView updates localStorage
function syncFromStorage() {
  storedPhoto.value    = localStorage.getItem('userPhoto') || ''
  storedUsername.value = localStorage.getItem('username') || ''
}

let syncInterval = null

onMounted(async () => {
  try {
    await axios.get('/liveness')
  } catch (e) {
    console.error('Liveness check failed', e)
  }
  syncFromStorage()
  window.addEventListener('storage', syncFromStorage)
  // Fallback polling so same-tab updates are caught too
  syncInterval = setInterval(syncFromStorage, 2000)
})

onUnmounted(() => {
  window.removeEventListener('storage', syncFromStorage)
  if (syncInterval) clearInterval(syncInterval)
})

// Event handlers for HomeView communication
const emitNewChat = () => {
  window.dispatchEvent(new CustomEvent('app-new-chat'))
}

const emitRefresh = () => {
  window.dispatchEvent(new CustomEvent('app-refresh'))
}

const handleNewChat = () => {
  // Handle if needed
}

const handleRefresh = () => {
  // Handle if needed
}

const logout = () => {
  try {
    const preserve = localStorage.getItem('leftConversations')
    localStorage.clear()
    if (preserve) localStorage.setItem('leftConversations', preserve)
  } catch (e) {
    console.error('Error clearing localStorage:', e)
  }
  
  try { 
    delete axios.defaults.headers.common['Authorization'] 
  } catch (e) {
    console.error('Error clearing axios headers:', e)
  }
  
  router.push('/login')
}
</script>

<style>
/* App Wrapper */
/* App Wrapper */
.app-wrapper {
  width: 100%;
  min-height: 100vh;
  background: var(--bg-primary);
  display: flex;
  flex-direction: column;
}

/* Unified Top Bar */
.app-topbar {
  position: sticky;
  top: 0;
  z-index: 200;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  height: 64px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  box-shadow: 0 1px 4px rgba(26, 115, 232, 0.12);
}
.actions-left,
.actions-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  margin-left: 2rem;
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 1rem;
}

.neon-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
  letter-spacing: normal;
}

.subtitle {
  color: var(--text-secondary);
  font-size: 1rem;
  font-weight: 400;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.btn-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 50%;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
}

.btn-icon:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

/* Avatar variant – no extra bg tint, just the image */
.btn-icon-avatar {
  padding: 0;
  background: transparent !important;
}

.topbar-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--border-color);
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.btn-icon-avatar:hover .topbar-avatar {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 8px rgba(26, 115, 232, 0.35);
}

.btn-icon-avatar.router-link-active .topbar-avatar {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 10px rgba(26, 115, 232, 0.4);
}

.btn-icon.router-link-active {
  color: var(--neon-cyan);
  background: var(--bg-hover);
}

.btn-logout {
  padding: 0.5rem;
  background: transparent;
  border: none;
  border-radius: 50%;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  transition: all 0.2s ease;
}

.btn-logout:hover {
  background: var(--bg-hover);
  color: #f15c6d; /* slightly softer red */
}

/* Main Content - MODIFIED */
main {
  flex: 1;
  width: 100%;
  min-height: 0; /* AGGIUNTO: permette di shrinkare correttamente */
}

main.fullscreen {
  height: 100vh;
}

main.with-topbar {
  height: calc(100vh - 76px); /* CHANGED: da min-height a height */
  /* overflow rimosso - lo scroll è gestito da .messages-neon */
}

/* Responsive */
@media (max-width: 768px) {
  .app-topbar {
    padding: 1rem 1.5rem;
    height: 68px; /* AGGIUNTO */
  }

  main.with-topbar {
    height: calc(100vh - 68px); /* AGGIORNATO */
  }

  .header-left {
    gap: 0.5rem;
  }

  .neon-title {
    font-size: 1.5rem;
    letter-spacing: 0.2rem;
  }

  .subtitle {
    font-size: 0.875rem;
  }

  .header-actions {
    gap: 0.5rem;
  }

  .btn-icon {
    width: 40px;
    height: 40px;
  }

  .btn-logout {
    padding: 0.6rem 1rem;
  }
}

@media (max-width: 480px) {
  .app-topbar {
    height: 64px; /* AGGIUNTO */
  }

  main.with-topbar {
    height: calc(100vh - 64px); /* AGGIORNATO */
  }

  .subtitle {
    display: none;
  }
  
  .header-actions {
    gap: 0.35rem;
  }
  
  .btn-icon {
    width: 36px;
    height: 36px;
  }
  
  .btn-logout {
    padding: 0.6rem 0.8rem;
  }
}
</style>
