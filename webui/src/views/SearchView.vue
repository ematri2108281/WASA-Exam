<template>
  <section class="search">
    <div class="search-form">
      <input v-model="qUser" class="input" type="text" placeholder="Search user (username)" @input="doSearch" />
    </div>

    <ErrorMsg v-if="error" :msg="error" />

    <div class="results">
      <div class="col">
        <h3>Users</h3>
        <div v-if="hasSearched && users.length === 0" class="muted">No users</div>
        <div v-for="u in users" :key="u.id" class="row" @click="selectUser(u.id)" :class="{ selected: selectedUserId === u.id }">
          <div class="avatar"><img :src="getUserPhoto(u)" :alt="u.username + ' photo'" /></div>
          <div class="label">{{ u.username }}</div>
          <input type="radio" class="picker" :value="u.id" v-model="selectedUserId" @click.stop="selectUser(u.id)" aria-label="Select user" />
        </div>
      </div>
      <div class="action-bar" v-if="hasSearched && users.length">
        <div class="bar-inner">
          <button class="btn btn-primary start" :disabled="!canStartChat" @click="startChat()">Start Chat</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'SearchView',
  emits: ['chat-started'],
  data() {
    return {
      qUser: '',
      users: [],
      selectedUserId: '',
      error: null,
      loading: false,
      currentUserId: localStorage.getItem('userId') || '',
      hasSearched: false,
    };
  },
  computed: {
    canSearch() {
      return this.qUser && this.qUser.trim().length > 0;
    },
    canStartChat() {
      return !!this.selectedUserId && !this.loading;
    },
  },
  methods: {
    letterAvatar(name, size = 60) {
      try {
        const letter = (String(name || '').trim().charAt(0) || '?').toUpperCase();
        const bg = '#e8f0fe';
        const fg = 'var(--neon-cyan)';
        const svg = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns='http://www.w3.org/2000/svg' width='${size}' height='${size}' viewBox='0 0 ${size} ${size}'>
  <rect width='100%' height='100%' rx='${Math.floor(size/2)}' ry='${Math.floor(size/2)}' fill='${bg}'/>
  <text x='50%' y='53%' dominant-baseline='middle' text-anchor='middle' fill='${fg}' font-family='Segoe UI, Roboto, sans-serif' font-weight='700' font-size='${Math.floor(size*0.5)}'>${letter}</text>
</svg>`;
        return 'data:image/svg+xml;utf8,' + encodeURIComponent(svg);
      } catch { return '/nopfp.png'; }
    },
    getUserPhoto(user) {
      if (!user) return this.letterAvatar('');
      const photo = user.photo;
      if (typeof photo === 'string') {
        if (photo.startsWith('data:')) return photo;
        if (photo.trim()) return `data:image/png;base64,${photo}`;
      }
      const name = user.username || user.name || '';
      return this.letterAvatar(name);
    },
    async doSearch() {
      if (!this.canSearch) return;
      this.loading = true;
      this.hasSearched = true;
      this.error = null;
      this.users = [];
      try {
        const params = new URLSearchParams();
        if (this.qUser.trim()) params.set('user', this.qUser.trim());
        const res = await this.$axios.get(`/searchby?${params.toString()}`);
        const allUsers = res.data?.users || [];

        const currentId = localStorage.getItem('userId') || this.currentUserId || '';
        this.users = allUsers.filter(u => {
          if (!u || !u.id) return false;
          return u.id !== currentId;
        });
        if (!this.users.find(u => u.id === this.selectedUserId)) {
          this.selectedUserId = '';
        }
        
        console.log('Search results:', this.users);
      } catch (e) {
        console.error('Search failed', e);
        this.error = 'Search failed';
      } finally {
        this.loading = false;
      }
    },
    selectUser(userId) {
      this.selectedUserId = userId;
    },
    async startChat(userId) {
      const targetId = userId || this.selectedUserId;
      if (!targetId) return;
      this.error = null;
      try {
        const res = await this.$axios.post('/direct-conversations', { peerUserId: targetId });
        const id = res?.data?.id || res?.data?.conversationId;
        if (id) {
          this.$emit('chat-started', id);
        }
      } catch (e) {
        console.error('Failed to start chat', e);
        this.error = 'Failed to start chat';
      }
    },
  },
};
</script>

<style scoped>
.search { display: grid; gap: 1rem; padding-bottom: 96px; }
.search-header { border-bottom: 1px solid var(--border); padding-bottom: .5rem; }
.search-form { display: flex; gap: .5rem; align-items: center; }
.input { flex: 1; padding: .6rem .7rem; border: 1px solid var(--border); border-radius: var(--radius); background: var(--bg); color: var(--text); }
.btn { padding: .55rem .9rem; border-radius: var(--radius); border: none; background: var(--accent); color: #000; font-weight: 700; }
.results { display: block; width: 100%; }
.col { background: var(--bg-alt); border: 1px solid var(--border); border-radius: var(--radius); padding: .75rem; width: 100%; }
.col h3 { margin: 0 0 .5rem 0; text-align: left; }
.row { display: flex; align-items: center; gap: .75rem; padding: .6rem .75rem; text-align: left; min-height: 80px; }
.row:hover { background: var(--bg-hover); border-radius: .5rem; }
.avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  overflow: hidden;
  background: var(--bg);
  position: relative;
  flex-shrink: 0;
  line-height: 0;
  margin-left: .5rem;
}
.avatar img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  display: block;
  border-radius: 50%;
}
.label { color: var(--text); flex: 1; text-align: left; font-weight: 700; font-size: 1.05rem; }
.muted { color: var(--text-dim); text-align: left; padding: .5rem; }
.btn { padding: .55rem .9rem; border-radius: var(--radius); border: none; background: var(--accent); color: #000; font-weight: 700; white-space: nowrap; }

.btn.start {
  min-width: 120px;
  width: auto;
  padding: .65rem 1.4rem;
  display: inline-flex;
  justify-content: center;
  background: var(--neon-cyan); 
  color: #ffffff; 
  border: 2px solid var(--neon-cyan);
  box-shadow: 0 0 12px rgba(26, 115, 232, 0.3);
  font-weight: 700;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn.start:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.5);
  background: var(--neon-purple);
}

.btn.start:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: rgba(26, 115, 232, 0.3);
  box-shadow: none;
}

.picker {
  margin-left: auto;
  margin-right: .75rem;
  width: 22px;
  height: 22px;
  min-width: 22px;
  min-height: 22px;
  max-width: 22px;
  max-height: 22px;
  cursor: pointer;
  flex-shrink: 0;
  appearance: none;
  -webkit-appearance: none;
  border: 2px solid var(--border);
  border-radius: 50%;
  background: var(--bg-alt);
  position: relative;
  display: grid;
  place-items: center;
  padding: 0;
  transition: background 0.2s, border-color 0.2s;
  box-sizing: border-box;
}
.picker::after {
  content: '';
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--accent);
  transform: scale(0);
  transition: transform 0.15s ease;
}
.picker:checked {
  border-color: var(--accent);
}
.picker:checked::after {
  transform: scale(1);
}
.picker:hover {
  border-color: var(--accent);
}
.row.selected { background: var(--bg-hover); border-radius: .5rem; }
.action-bar {
  position: fixed;
  left: 50%;
  transform: translateX(-50%);
  bottom: 12px;
  display: flex;
  justify-content: center;
  padding: 0;
  background: transparent;
  z-index: 8;
  width: auto;
}
.bar-inner {
  background: var(--bg-alt);
  border: 1px solid var(--border);
  border-radius: 999px;
  box-shadow: 0 6px 20px rgba(0,0,0,0.35);
  padding: 0.4rem 0.8rem;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: auto;
  min-width: 140px;
  max-width: 260px;
}
</style>
