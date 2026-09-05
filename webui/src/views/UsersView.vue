<template>
  <section class="users-view">
    <header class="users-header">
      <div>
        <h1>Users</h1>
        <p class="muted">All profiles from the database</p>
      </div>
      <div class="actions">
        <button type="button" class="btn-refresh" @click="fetchUsers" :disabled="loading">
          {{ loading ? 'Refreshing…' : 'Refresh' }}
        </button>
      </div>
    </header>

    <ErrorMsg v-if="error" :msg="error" />

    <div v-if="loading" class="loading">Loading users…</div>
    <div v-else-if="hasUsers" class="users-list">
      <div v-for="user in users" :key="user.id" class="user-card">
        <img :src="photoSrc(user)" :alt="user.username" class="avatar" />
        <div class="info">
          <div class="username">{{ user.username }}</div>
        </div>
      </div>
    </div>
    <div v-else class="empty">
      No users found.
    </div>
  </section>
</template>

<script>
export default {
  name: 'UsersView',

  data() {
    return {
      loading: false,
      error: null,
      users: [],
    };
  },
  computed: {
    hasUsers() {
      return Array.isArray(this.users) && this.users.length > 0;
    },
  },
  methods: {
    letterAvatar(name, size = 64) {
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
      } catch {
        return '/nopfp.png';
      }
    },
    photoSrc(user) {
      try {
        if (user?.photo && typeof user.photo === 'string' && user.photo.length > 0) {
          return user.photo.startsWith('data:') ? user.photo : `data:image/png;base64,${user.photo}`;
        }
        return this.letterAvatar(user?.username || 'U', 56);
      } catch {
        return this.letterAvatar('U', 56);
      }
    },
    async fetchUsers() {
      this.loading = true;
      this.error = null;
      try {
        const token = localStorage.getItem('token');
        const cfg = token ? { headers: { Authorization: `Bearer ${token}` } } : {};
        const res = await this.$axios.get('/users', cfg);
        this.users = res?.data?.users || [];
      } catch (e) {
        this.error = e?.response?.data?.error || 'Failed to load users';
      } finally {
        this.loading = false;
      }
    },
  },
  mounted() {
    this.fetchUsers();
  },
};
</script>

<style scoped>
.users-view {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  background: var(--bg);
  color: var(--text);
  padding: 1.25rem;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  margin-top: 1rem;
}

.users-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.users-header h1 {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 700;
}

.muted {
  color: var(--muted, #6b7280);
  margin: 0;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.btn-refresh {
  padding: 0.5rem 1rem;
  background: var(--bg-alt);
  color: var(--accent);
  border: 1px solid var(--accent);
  border-radius: 999px;
  font-weight: 700;
  cursor: pointer;
}

.btn-refresh:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.users-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-top: 1rem;
}

.user-card {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 0.75rem;
  padding: 0.75rem;
  background: var(--bg-alt);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  align-items: center;
}

.avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid var(--border);
}

.info {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.username {
  font-weight: 700;
  font-size: 1rem;
}

.loading,
.empty {
  padding: 1rem;
  color: var(--muted, #6b7280);
}
</style>
