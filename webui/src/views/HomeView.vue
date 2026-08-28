<template>
  <div class="home-container">
   

    <!-- Error Message -->
    <div v-if="errormsg" class="error-banner">
      {{ errormsg }}
      <button @click="errormsg = null" class="error-close">×</button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading conversations...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="!hasConversations" class="empty-state">
      <div class="empty-icon">💬</div>
      <h2>No conversations yet</h2>
      <p>Start chatting with someone!</p>
      <button class="neon-button" @click="openSearchModal">
        Start New Chat
      </button>
    </div>

    <!-- Chat Layout -->
    <div v-else class="chat-layout">
      <!-- Sidebar: Conversations List -->
      <aside class="conversations-sidebar" :class="{ 'forward-mode': forwardMode.active }">

  <!-- Header: Search o Forward Mode -->
  <div class="sidebar-header">
    
    <!-- Normal Search -->
    <div class="search-box">
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="search-icon">
        <circle cx="11" cy="11" r="8"></circle>
        <path d="m21 21-4.35-4.35"></path>
      </svg>
      <input
        v-model="searchQuery"
        class="neon-search"
        placeholder="Search chats..."
      />
    </div>
  </div>

  <!-- Lista delle conversazioni -->
  <div class="conversations-list">
    <div
      v-for="conv in filteredSortedConversations"
      :key="conv.id"
      class="conv-item"
      :class="{ active: selectedConvId === conv.id }"
      @click="selectConversation(conv.id)"
    >
      <div class="conv-avatar">
        <img :src="getConvPhoto(conv)" :alt="conv.name || 'Chat'" />
      </div>
      <div class="conv-info">
        <div class="conv-header">
          <span class="conv-name">{{ getConvName(conv) }}</span>
          <span class="conv-time">{{ getLastMessageTime(conv) }}</span>
        </div>
        <p class="conv-preview">{{ getLastMessagePreview(conv) }}</p>
      </div>
    </div>
  </div>
</aside>



      <!-- Main: Chat View -->
      <main class="chat-main">
         <!-- Forward Mode Banner -->
          <div v-if="selectedConvId" class="chat-shell">
  <!-- Badge Forward Mode -->
  <div v-if="forwardMode.active" class="forward-badge">
    📤 Forwarding
    <button class="exit-forward" @click="exitForwardMode">×</button>
  </div>

  <ConvView
    :key="selectedConvId"
    :conversationId="selectedConvId"
    :embedded="true"
    @conversation-deleted="onConversationDeleted"
    @group-updated="onGroupUpdated"
  />
</div>

        <div v-else class="no-selection">
          <div class="selection-icon">💭</div>
          <h3>Select a conversation</h3>
          <p>Choose a chat from the sidebar to start messaging</p>
        </div>
      </main>
    </div>

    <!-- Search Modal -->
    <div v-if="searchModalOpen" class="modal-overlay" @click.self="closeSearchModal">
      <div class="modal-card">
        <div class="modal-header">
          <h2>New Chat</h2>
          <button class="btn-close" @click="closeSearchModal">×</button>
        </div>
        <div class="modal-body">
          <SearchView @chat-started="onChatStarted" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>

import ConvView from './ConvView.vue'
import SearchView from './SearchView.vue'

export default {
  name: 'HomeView',
  components: { ConvView, SearchView },
  
  data() {
    return {
      errormsg: null,
      loading: false,
      conversations: [],
      selectedConvId: null,
      userId: localStorage.getItem('userId') || null,
      searchQuery: '',
      forwardMode: { 
        active: false, 
        fromConvId: '', 
        messageId: '', 
        search: '', 
        suggestions: [], 
        loading: false, 
        timer: null 
      },
      searchModalOpen: false,
      pollingInterval: null,
    }
  },
  
  computed: {
    hasConversations() {
      return Array.isArray(this.conversations) && this.conversations.length > 0
    },
    
    sortedConversations() {
      const arr = Array.isArray(this.conversations) ? [...this.conversations] : []
      const tsOf = (c) => {
        try {
          const lm = c?.lastMessage || c?.LastMessage || {}
          const t1 = lm?.timestamp ? new Date(lm.timestamp).getTime() : NaN
          const t2 = c?.createdAt ? new Date(c.createdAt).getTime() : NaN
          const t3 = c?.CreatedAt ? new Date(c.CreatedAt).getTime() : NaN
          const t4 = c?._snapshotLastTs || NaN
          return [t1, t2, t3, t4].find(v => Number.isFinite(v)) || 0
        } catch { 
          return 0 
        }
      }
      arr.sort((a, b) => tsOf(b) - tsOf(a))
      return arr
    },
    
    filteredSortedConversations() {
      const query = (this.searchQuery || '').trim().toLowerCase()
      if (!query) return this.sortedConversations
      return this.sortedConversations.filter(c => {
        const name = this.getConvName(c).toLowerCase()
        return name.includes(query)
      })
    }
  },
  
  methods: {
    getConvName(conv) {
    if (!conv) return 'Chat'
    
    const isGroup = conv.isGroup === true || (conv.membersIds?.length || 0) > 2
    if (isGroup) {
      return conv.name || conv.displayName || 'Group Chat'
    }
    
    try {
      const myUserId = String(this.userId || '')
      const participants = conv.participants || []

      const peer = participants.find(p => {
        const pId = String(p?.id || p?.userId || p?.ID || '')
        return pId && pId !== myUserId
      })
      
      if (peer?.username) {
        return peer.username
      }
      
      // Fallback con members
      const members = conv.members || []
      const peerMember = members.find(m => {
        const mId = String(m?.id || m?.userId || '')
        return mId && mId !== myUserId
      })
      
      if (peerMember?.username) {
        return peerMember.username
      }
      
      // Ultimo fallback
      return conv.name || conv.displayName || 'Direct Chat'
    } catch (error) {
      console.error('Error getting conversation name:', error)
      return conv.name || 'Chat'
    }
  },
    
    letterAvatar(name, size = 64) {
      try {
        const letter = (String(name || '').trim().charAt(0) || '?').toUpperCase()
        const bg = '#e8f0fe'
        const fg = 'var(--neon-cyan)'
        const svg = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns='http://www.w3.org/2000/svg' width='${size}' height='${size}' viewBox='0 0 ${size} ${size}'>
  <rect width='100%' height='100%' rx='${Math.floor(size/2)}' ry='${Math.floor(size/2)}' fill='${bg}'/>
  <text x='50%' y='53%' dominant-baseline='middle' text-anchor='middle' fill='${fg}' font-family='Segoe UI, Roboto, sans-serif' font-weight='700' font-size='${Math.floor(size*0.5)}'>${letter}</text>
</svg>`
        return 'data:image/svg+xml;utf8,' + encodeURIComponent(svg)
      } catch { 
        return '/nopfp.png' 
      }
    },
    
    async refresh(showLoading = true) {
      if (showLoading) {
        this.loading = true
      }
      this.errormsg = null
      
      try {
        const token = localStorage.getItem('token')
        const response = await this.$axios.get('/conversations', 
          token ? { headers: { Authorization: `Bearer ${token}` } } : {}
        )
        
        const serverConvs = response.data || []
        let merged = Array.isArray(serverConvs) ? [...serverConvs] : []
        
        try {
          const raw = localStorage.getItem('leftConversations') || '{}'
          const map = JSON.parse(raw) || {}
          const me = String(this.userId || '')
          const idsPresent = new Set(merged.map(c => String(c?.id || '')))
          
          for (const [convId, entry] of Object.entries(map)) {
            if (!entry || String(entry.userId || '') !== me) continue
            if (idsPresent.has(String(convId))) continue
            
            const snap = entry.conversation || {}
            const msgs = Array.isArray(snap.messages) ? snap.messages : []
            const lastMsg = msgs.length ? msgs[msgs.length - 1] : null
            const lastTs = lastMsg?.timestamp 
              ? new Date(lastMsg.timestamp).getTime() 
              : (entry.leftAt ? new Date(entry.leftAt).getTime() : Date.now())
            
            merged.push({
              id: convId,
              name: snap.name || 'Group',
              isGroup: true,
              profilePhoto: snap.profilePhoto || snap.photo || null,
              lastMessage: lastMsg ? { timestamp: lastMsg.timestamp } : undefined,
              _snapshotLastTs: lastTs
            })
          }
        } catch (e) {
          console.error('Failed to load left conversations:', e)
        }

        this.conversations = merged
        
        if (this.sortedConversations.length > 0 && !this.selectedConvId) {
          this.selectedConvId = this.sortedConversations[0].id
        }
      } catch (e) {
        console.error('Failed to load conversations:', e)
        this.errormsg = e?.response?.data?.error || e.message || 'Failed to load conversations'
      } finally {
        if (showLoading) {
          this.loading = false
        }
      }
    },
    
    openSearchModal() {
      this.searchModalOpen = true
    },
    
    closeSearchModal() {
      this.searchModalOpen = false
    },
    
    async onChatStarted(convId) {
      this.searchModalOpen = false
      try {
        await this.refresh()
      } catch (e) {
        console.error('Failed to refresh after chat started:', e)
      }
      this.selectedConvId = convId
      this.$router.push({ path: '/home', query: { conv: convId } })
    },
    
    selectConversation(convId) {
      if (this.forwardMode.active) {
        this.forwardToConversation(convId)
        return
      }
      this.selectedConvId = convId
      this.$router.push({ path: '/home', query: { conv: convId } })
    },
    
    async forwardToConversation(targetId) {
      if (!this.forwardMode.messageId || !this.forwardMode.fromConvId || !targetId) return
      
      try {
        const token = localStorage.getItem('token')
        await this.$axios.post(
          `/conversations/${this.forwardMode.fromConvId}/messages/${this.forwardMode.messageId}/forward`,
          { targetConversationId: targetId },
          token ? { headers: { Authorization: `Bearer ${token}` } } : {}
        )
        
        this.selectedConvId = targetId
        this.exitForwardMode()
        this.$router.push({ path: '/home', query: { conv: targetId } })
      } catch (e) {
        this.errormsg = e?.response?.data?.error || 'Failed to forward message'
      }
    },
    
    exitForwardMode() {
      this.forwardMode = { 
        active: false, 
        fromConvId: '', 
        messageId: '', 
        search: '', 
        suggestions: [], 
        loading: false, 
        timer: null 
      }
      this.searchQuery = ''
    },
    
    onForwardSearchInput() {
      const q = (this.forwardMode.search || '').trim()
      if (this.forwardMode.timer) clearTimeout(this.forwardMode.timer)
      if (!q) { 
        this.forwardMode.suggestions = []
        this.forwardMode.loading = false
        return 
      }
      
      this.forwardMode.loading = true
      this.forwardMode.timer = setTimeout(async () => {
        try {
          const token = localStorage.getItem('token')
          const headers = token ? { Authorization: `Bearer ${token}` } : {}
          const res = await this.$axios.get(`/searchby?user=${encodeURIComponent(q)}`, { headers })
          const all = res?.data?.users || []
          this.forwardMode.suggestions = all.filter(u => u && u.id && u.id !== this.userId).slice(0, 10)
        } catch (e) { 
          this.forwardMode.suggestions = [] 
        } finally { 
          this.forwardMode.loading = false 
        }
      }, 250)
    },
    
    async selectForwardUser(user) {
  if (!user || !this.forwardMode.messageId || !this.forwardMode.fromConvId) return
  
  try {
    const token = localStorage.getItem('token')
    const headers = token ? { Authorization: `Bearer ${token}` } : {}
    
    const myUserId = String(this.userId || '')
    const existingConv = this.conversations.find(conv => {
      // Salta i gruppi
      if (conv.isGroup === true) return false
      
      const participants = conv.participants || []
      // Controlla se c'è questo utente tra i partecipanti (chat 1-a-1)
      return participants.some(p => {
        const pId = String(p?.id || p?.userId || '')
        return pId === String(user.id)
      }) && participants.length === 2 // Assicurati che sia una chat diretta
    })
    
    let targetId
    
    if (existingConv) {
      // Usa la conversazione esistente
      targetId = existingConv.id || existingConv.conversationId
    } else {
      // Crea nuova conversazione solo se non esiste
      const response = await this.$axios.post('/direct-conversations', { peerUserId: user.id }, { headers })
      const conv = response?.data || {}
      targetId = conv.id || conv.ID || conv.conversationId || conv.conversation_id
    }
    
    if (!targetId) { 
      this.errormsg = 'Failed to create conversation'
      return 
    }
    
    // Forward il messaggio
    await this.$axios.post(
      `/conversations/${this.forwardMode.fromConvId}/messages/${this.forwardMode.messageId}/forward`,
      { targetConversationId: targetId },
      { headers }
    )
    
    await this.refresh()
    
    this.selectedConvId = targetId
    this.exitForwardMode()
    this.$router.push({ path: '/home', query: { conv: targetId } })
  } catch (e) {
    console.error('Failed to forward message:', e)
    this.errormsg = e?.response?.data?.error || 'Failed to forward message'
  }
},
    
    onConversationDeleted(payload) {
      const deletedId = String(payload?.deletedId || '')
      const nextFromChild = payload?.nextId ? String(payload.nextId) : ''
      
      this.conversations = (this.conversations || []).filter(c => {
        const cid = String(c?.id || c?.conversationId || c?.conversation_id || c?.ID || '')
        return cid && cid !== deletedId
      })
      
      this.selectedConvId = null
      
      this.$nextTick(async () => {
        await this.refresh()
        const fallback = this.sortedConversations[0]
        const nextId = nextFromChild || fallback?.id || fallback?.conversationId || fallback?.conversation_id || null
        this.selectedConvId = nextId || null
        
        if (nextId) {
          this.$router.push({ path: '/home', query: { conv: nextId } })
        } else {
          this.$router.push('/home')
        }
      })
    },
    
    async onGroupUpdated(payload) {
      await this.refresh()
      if (payload?.convId) {
        this.selectedConvId = payload.convId
      }
    },
    
    getConvPhoto(conv) {
    if (!conv) return this.letterAvatar('')
    
    // Per gruppi: usa conv.photo
    const isGroup = conv.isGroup === true || (conv.membersIds?.length || 0) > 2
  if (isGroup) {
    const photo = conv.groupPhoto || conv.profilePhoto || conv.photo || ''
    if (typeof photo === 'string') {
      if (photo.startsWith('data:')) return photo
      if (photo.trim()) return `data:image/png;base64,${photo}`
    }
    const name = conv.name || 'Group'
    return this.letterAvatar(name)
  }
    
    try {
      const myUserId = String(this.userId || '')
      const participants = conv.participants || []
      
      const peer = participants.find(p => {
        const pId = String(p?.id || p?.userId || p?.ID || '')
        return pId && pId !== myUserId
      })
      
      if (peer) {
        const photoB64 = peer.photo || peer.Photo || peer.profilePhoto || peer.ProfilePhoto
        
        if (photoB64) {
          if (photoB64.startsWith('data:')) return photoB64
          return `data:image/png;base64,${photoB64}`
        }
        
        const peerName = peer.username || peer.Username || peer.name || peer.Name || 'User'
        return this.letterAvatar(peerName)
      }
      
      const name = conv.name || 'Chat'
      return this.letterAvatar(name)
    } catch (error) {
      console.error('Error getting conversation photo:', error)
      return this.letterAvatar('?')
    }
  },
    
    getLastMessagePreview(conv) {
      try {
        const lm = conv?.lastMessage || conv?.LastMessage || {}
        const contentType = lm?.content?.type || lm?.Content?.Type || ''
        const contentValue = lm?.content?.value || lm?.Content?.Value || ''
        
        if (contentType.toLowerCase() === 'image') {
          return '📷 Photo'
        }
        
        if (contentValue) {
          try {
            const decoded = atob(contentValue)
            return decoded.length > 50 ? decoded.substring(0, 50) + '...' : decoded
          } catch {
            return contentValue.length > 50 ? contentValue.substring(0, 50) + '...' : contentValue
          }
        }
        
        return 'No messages yet'
      } catch {
        return 'No messages yet'
      }
    },
    
    getLastMessageTime(conv) {
      try {
        const lm = conv?.lastMessage || conv?.LastMessage || {}
        const timestamp = lm?.timestamp || lm?.Timestamp || lm?.createdAt || lm?.CreatedAt || ''
        if (!timestamp) return ''
        
        const date = new Date(timestamp)
        if (isNaN(date.getTime())) return ''
        
        const now = new Date()
        const isToday = date.getDate() === now.getDate() && 
                        date.getMonth() === now.getMonth() && 
                        date.getFullYear() === now.getFullYear()
        
        const isYesterday = () => {
          const yesterday = new Date(now)
          yesterday.setDate(yesterday.getDate() - 1)
          return date.getDate() === yesterday.getDate() && 
                 date.getMonth() === yesterday.getMonth() && 
                 date.getFullYear() === yesterday.getFullYear()
        }
        
        if (isToday) {
          return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
        } else if (isYesterday()) {
          return 'Yesterday'
        } else if ((now - date) < 7 * 24 * 60 * 60 * 1000) {
          return date.toLocaleDateString('en-US', { weekday: 'short' })
        } else {
          return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
        }
      } catch {
        return ''
      }
    }
  },
  
  mounted() {
    this.refresh()
    
    if (this.$route.query.conv) {
      this.selectedConvId = this.$route.query.conv
    }
    
    if (this.$route.query.forward && this.$route.query.from) {
      this.forwardMode.active = true
      this.forwardMode.messageId = String(this.$route.query.forward)
      this.forwardMode.fromConvId = String(this.$route.query.from)
    }
    
    window.addEventListener('app-new-chat', this.openSearchModal)
    window.addEventListener('app-refresh', () => this.refresh())
    
    this.pollingInterval = setInterval(() => {
      this.refresh(false)
    }, 3000)
  },
  
  beforeUnmount() {
    if (this.pollingInterval) {
      clearInterval(this.pollingInterval)
    }
    
    window.removeEventListener('app-new-chat', this.openSearchModal)
    window.removeEventListener('app-refresh', () => this.refresh())
  },
  
  watch: {
    '$route.query.conv'(newId) {
      if (newId) {
        this.selectedConvId = newId
      }
    },
    '$route.query.forward'(msgId) {
      if (msgId && this.$route.query.from) {
        this.forwardMode.active = true
        this.forwardMode.messageId = String(msgId)
        this.forwardMode.fromConvId = String(this.$route.query.from)
      }
    }
  }
}
</script>

<style scoped>
.home-container {
  position: relative;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow: hidden;
}



@keyframes float-up {
  0% {
    transform: translateY(0) translateX(0);
    opacity: 1;
  }
  100% {
    transform: translateY(-100vh) translateX(30px);
    opacity: 0;
  }
}

.error-banner {
  position: relative;
  z-index: 10;
  padding: 1rem 2rem;
  background: rgba(255, 0, 100, 0.15);
  border-bottom: 1px solid rgba(255, 0, 100, 0.4);
  color: #ff0064;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.error-close {
  background: transparent;
  border: none;
  color: #ff0064;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0 0.5rem;
}


.loading-container {
  position: relative;
  z-index: 10;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 1.5rem;
  color: var(--neon-cyan);
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(26, 115, 232, 0.2);
  border-top-color: var(--neon-cyan);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Empty State */
.empty-state {
  position: relative;
  z-index: 10;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 3rem;
  text-align: center;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1.5rem;
  filter: drop-shadow(0 0 20px rgba(26, 115, 232, 0.3));
}

.empty-state h2 {
  color: var(--text-primary);
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
}

.empty-state p {
  color: var(--text-secondary);
  margin-bottom: 2rem;
}

.neon-button {
  padding: 1rem 2rem;
  background: linear-gradient(90deg, var(--neon-cyan) 0%, var(--neon-purple) 100%);
  border: none;
  border-radius: 12px;
  color: #ffffff;
  font-size: 1rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.3);
}

.neon-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 0 30px rgba(26, 115, 232, 0.5);
}

/* Chat Layout */
.chat-layout {
  position: relative;
  z-index: 10;
  display: grid;
  grid-template-columns: 350px 1fr;
  flex: 1;
  overflow: hidden;
}


.conversations-sidebar {
  display: flex;
  flex-direction: column;
  background: #ffffff;
  border-right: 1px solid var(--border-color);
  backdrop-filter: none;
  overflow: hidden; 
}

.sidebar-header {
  flex-shrink: 0; 
  padding-top: 75px;
  background: #ffffff;
  backdrop-filter: none;
  border-bottom: 1px solid var(--border-color);
}

.conversations-list {
  flex: 1;         
  overflow-y: auto; 
  padding-top: 0;   
}


.search-box {
  position: sticky;
  top: 10px;
  padding: 1rem;
  background: #ffffff;
  backdrop-filter: none;
  border-bottom: 1px solid var(--border-color);
  z-index: 5;
}



.search-icon {
  position: absolute;
  left: 2rem;
  top: 50%;
  transform: translateY(-50%);
  color: var(--neon-cyan);
  pointer-events: none;
}

.neon-search {
  width: 100%;
  padding: 0.875rem 1rem 0.875rem 3rem;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  color: var(--text-primary);
  font-size: 0.95rem;
  transition: all 0.3s ease;
  outline: none;
}

.neon-search::placeholder {
  color: var(--text-dim);
}

.neon-search:focus {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 8px rgba(26, 115, 232, 0.2);
  background: #ffffff;
}
.conversations-sidebar.forward-mode {
  border-right: 4px solid var(--neon-cyan);
  box-shadow: 0 0 15px rgba(26, 115, 232, 0.2) inset;
  animation: neon-glow 1.5s infinite alternate;
}

@keyframes neon-glow {
  0% {
    box-shadow: 0 0 5px rgba(26, 115, 232, 0.15) inset;
  }
  50% {
    box-shadow: 0 0 20px rgba(26, 115, 232, 0.3) inset;
  }
  100% {
    box-shadow: 0 0 10px rgba(26, 115, 232, 0.2) inset;
  }
}
.conversations-sidebar.forward-mode .conv-item {
  border-left: 4px solid rgba(26, 115, 232, 0.5);
}

.conversations-sidebar.forward-mode .conv-item.active:hover {
  background-color: rgba(26, 115, 232, 0.1);
}


.forward-badge button.exit-forward {
  background: transparent;
  border: none;
  font-size: 14px;
  cursor: pointer;
  color: var(--text-primary);
}

.forward-badge button.exit-forward:hover {
  color: #ff004f;
}


.btn-exit {
  background: transparent;
  border: none;
  color: var(--neon-purple);
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
}

.forward-loading {
  margin-top: 0.5rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
  text-align: center;
}
.forward-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  cursor: pointer;
  border-bottom: 1px solid rgba(26, 115, 232, 0.10);
  transition: background 0.2s ease;
}

.forward-item:last-child {
  border-bottom: none;
}

.forward-item:hover {
  background: rgba(26, 115, 232, 0.10);
}

.user-avatar {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--neon-cyan), var(--neon-purple));
  border-radius: 50%;
  color: var(--bg-primary);
  font-weight: 700;
  font-size: 0.875rem;
}

.forward-hint {
  margin-top: 0.75rem;
  color: var(--text-secondary);
  font-size: 0.75rem;
  text-align: center;
}


.conversations-list {
  flex: 1;
  overflow-y: auto;
  
}

.conv-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.25rem;
  cursor: pointer;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.2s ease;
}

.conv-item:hover {
  background: var(--bg-hover);
}

.conv-item.active {
  background: rgba(26, 115, 232, 0.1);
  border-left: 3px solid var(--neon-cyan);
  box-shadow: inset 0 0 20px rgba(26, 115, 232, 0.05);
}

.conv-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  border: 2px solid rgba(26, 115, 232, 0.25);
}

.conv-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.conv-info {
  flex: 1;
  min-width: 0;
}

.conv-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
}

.conv-name {
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conv-time {
  font-size: 0.75rem;
  color: var(--text-secondary);
  flex-shrink: 0;
}

.conv-preview {
  font-size: 0.875rem;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin: 0;
}

/* Chat Main */
.chat-main {
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow: hidden;
}

.chat-shell {
  position: relative; 
  width: 100%;
  height: 100%;
  overflow: hidden;
}


.no-selection {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  padding: 2rem;
}

.selection-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  filter: drop-shadow(0 0 15px rgba(26, 115, 232, 0.25));
}

.no-selection h3 {
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.no-selection p {
  color: var(--text-secondary);
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(8px);
}

.modal-card {
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  background: #ffffff;
  border: 1px solid var(--border-color);
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 4px 24px rgba(26, 115, 232, 0.15);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h2 {
  color: var(--neon-cyan);
  font-size: 1.25rem;
  margin: 0;
}

.btn-close {
  background: transparent;
  border: none;
  color: var(--neon-cyan);
  font-size: 2rem;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.btn-close:hover {
  color: #ff0064;
  transform: rotate(90deg);
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  max-height: calc(80vh - 80px);
}


@media (max-width: 768px) {
  .chat-layout {
    grid-template-columns: 1fr;
  }
  
  .conversations-sidebar {
    display: none;
  }
}
</style>
