<template>
  <div class="login-container">
    <div class="login-card">
      <!-- Brand -->
      <div class="logo-section">
        <div class="logo-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
        </div>
        <h1 class="brand-name">WasaText</h1>
        <p class="tagline">Messaging app project</p>
      </div>

      <!-- Login Form -->
      <form @submit.prevent="handleLogin" class="login-form">
        <!-- Username Input -->
        <div class="field-group">
          <label for="username-input" class="field-label">Username</label>
          <div class="input-wrapper">
            <svg class="input-icon" xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
            <input
              id="username-input"
              v-model="username"
              type="text"
              class="text-input"
              placeholder="Enter your username"
              required
              minlength="3"
              maxlength="20"
              :disabled="loading"
              @input="clearError"
            />
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="error-message" role="alert">
          <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
          </svg>
          {{ error }}
        </div>

        <!-- Login Button -->
        <button type="submit" id="login-btn" class="login-btn" :disabled="loading">
          <span v-if="!loading">Sign In</span>
          <span v-else class="btn-loading">
            <span class="spinner"></span>
            Signing in…
          </span>
        </button>

        <!-- Info Text -->
        <p class="info-text">
          Enter any username to sign in or create a new account automatically.
        </p>
      </form>
    </div>

    <!-- Footer -->
    <footer class="login-footer">
      <p>WasaText &bull; 2026</p>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'LoginView',
  data() {
    return {
      username: '',
      error: null,
      loading: false
    }
  },
  methods: {
    async handleLogin() {
      const name = this.username.trim()
      if (!name) { this.error = "Username cannot be empty"; return }
      if (name.length < 3) { this.error = "Username must be at least 3 characters"; return }
      if (name.length > 20) { this.error = "Username must be 20 characters or less"; return }

      this.loading = true
      this.error = null

      try {
        const response = await this.$axios.post('/login', { name })
        const data = response.data
        const token = data.identifier || data.token
        const userId = data.id
        const username = data.username
        const photo = data.photo

        if (!token || !userId) throw new Error('Invalid server response')

        localStorage.setItem('token', token)
        localStorage.setItem('userId', userId)
        localStorage.setItem('username', username || name)

        if (photo) {
          localStorage.setItem('userPhoto', photo.startsWith('data:image') ? photo : `data:image/png;base64,${photo}`)
        } else {
          localStorage.setItem('userPhoto', '')
        }

        if (this.$axios.defaults?.headers) {
          this.$axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        }

        this.$router.push('/home')

      } catch (err) {
        console.error('Login error:', err)
        if (err.response) {
          const status = err.response.status
          const errorData = err.response.data
          this.error = status === 400 ? (errorData?.error || 'Invalid username format') :
                       status === 409 ? 'Username already taken' :
                       status === 401 ? 'Authentication failed' :
                       status === 500 ? 'Server error. Please try again later.' :
                       errorData?.error || `Server error (${status})`
        } else if (err.request) {
          this.error = 'Cannot connect to server. Is it running?'
        } else {
          this.error = err.message || 'Login failed'
        }
      } finally {
        this.loading = false
      }
    },
    clearError() { this.error = null }
  },
  mounted() {
    try { localStorage.clear() } catch (e) { console.error(e) }
    try { delete this.$axios.defaults.headers.common['Authorization'] } catch (e) {}
  }
}
</script>

<style scoped>
/* ── Layout ── */
.login-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: var(--bg-primary);
  padding: 2rem 1rem;
}

/* ── Card ── */
.login-card {
  width: 100%;
  max-width: 420px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 2.5rem 2.25rem;
  box-shadow: var(--shadow-neon);
}

/* ── Brand ── */
.logo-section {
  text-align: center;
  margin-bottom: 2.25rem;
}

.logo-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 58px;
  height: 58px;
  background: var(--bg-hover);
  border-radius: var(--radius-sm);
  color: var(--neon-cyan);
  margin-bottom: 1rem;
  border: 1px solid var(--border-color);
}

.brand-name {
  font-size: 1.6rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 0.35rem;
  letter-spacing: -0.02em;
}

.tagline {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

/* ── Form ── */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.1rem;
}

.field-group {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.field-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--text-secondary);
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 0.85rem;
  color: var(--text-secondary);
  pointer-events: none;
  flex-shrink: 0;
}

.text-input {
  width: 100%;
  padding: 0.7rem 0.9rem 0.7rem 2.6rem;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 0.95rem;
  font-family: inherit;
  transition: border-color 0.15s, box-shadow 0.15s;
  outline: none;
  box-sizing: border-box;
}

.text-input::placeholder {
  color: var(--text-dim);
}

.text-input:focus {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 0 3px rgba(26, 115, 232, 0.15);
}

.text-input:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

/* ── Error ── */
.error-message {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.65rem 0.9rem;
  background: rgba(241, 92, 109, 0.1);
  border: 1px solid rgba(241, 92, 109, 0.35);
  border-radius: var(--radius-sm);
  color: #f15c6d;
  font-size: 0.84rem;
  line-height: 1.4;
}

/* ── Button ── */
.login-btn {
  width: 100%;
  padding: 0.8rem 1rem;
  background: var(--neon-cyan);
  border: none;
  border-radius: var(--radius-sm);
  color: #ffffff;
  font-size: 0.95rem;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: background 0.15s, transform 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 0.25rem;
}

.login-btn:hover:not(:disabled) {
  background: var(--neon-purple);
}

.login-btn:active:not(:disabled) {
  transform: scale(0.99);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* ── Loading ── */
.btn-loading {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Info text ── */
.info-text {
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.8rem;
  margin: 0;
  line-height: 1.5;
}

/* ── Footer ── */
.login-footer {
  margin-top: 1.75rem;
  text-align: center;
  color: var(--text-dim);
  font-size: 0.75rem;
}

/* ── Responsive ── */
@media (max-width: 480px) {
  .login-card {
    padding: 2rem 1.5rem;
  }

  .brand-name {
    font-size: 1.4rem;
  }
}
</style>
