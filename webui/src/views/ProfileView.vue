<template>
  <section class="profile centered">
    <div class="profile-card">
       <button class="btn btn-back" @click="goHome">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"></line>
          <polyline points="12 19 5 12 12 5"></polyline>
        </svg>
        Back to Home
      </button>
      <div class="photo-wrap" :class="{ clickable: editProfile }" @click="editProfile && triggerPhotoUpload()">
        <img :src="userPhotoSrc" alt="Profile Photo" class="pfp" width="140" height="140" />
        <div v-if="editProfile" class="photo-overlay">Click to change</div>
      </div>
      <h2 v-if="!editProfile" class="username">@{{ username }}</h2>

      <div class="controls">
        <button v-if="!editProfile" class="btn btn-edit" @click="openEdit">Edit Profile</button>
        <button v-if="!editProfile" class="btn btn-logout" @click="logOut">Log Out</button>
      </div>

      <div v-if="editProfile" class="editor">
        <div class="field-row username-edit">
          <span class="username-prefix">@</span>
          <input v-model="newUsername" class="input username-input" placeholder="Enter new username" maxlength="16" minlength="1" />
          <button class="btn btn-update" @click="updateUsername" :disabled="!canUpdateUsername">Update Username</button>
        </div>
        <button class="btn btn-cancel" @click="editProfile = false" style="justify-self: center; width: auto; padding: .6rem 1.5rem;">Done</button>
        <input ref="fileInput" type="file" accept="image/*" @change="handlePhotoUpload" style="display: none;" />
      </div>

      <div v-if="editProfile && successMsg" class="success-message">{{ successMsg }}</div>

      <ErrorMsg v-if="errormsg" :msg="errormsg" />
    </div>
  </section>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from '@/components/ErrorMsg.vue';

export default {
    name: 'ProfileView',
    components: { ErrorMsg },
    data() {
        return {
            editProfile: false,
            errormsg: null,
            successMsg: null,
            successTimer: null,
            user: {
                photo: localStorage.getItem('userPhoto') || '/nopfp.png',
            },
            username: localStorage.getItem('username') || '',
            newUsername: '',
            newPhoto: null,
        }
    },
    computed: {
        userPhotoSrc() {
            const p = this.user.photo || '';
            if (p && (p.startsWith('data:') || p.startsWith('http'))) return p;
            const name = this.username || '';
            return this.letterAvatar(name, 140);
        },
        canUpdateUsername() {
            const v = (this.newUsername || '').trim();
            return !!v && v !== this.username;
        },
    },
    methods: {
        letterAvatar(name, size = 140) {
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
        goHome() {
          this.$router.push('/home');
        },
        initFromLocal() {
            const name = localStorage.getItem('username') || '';
            let photo = localStorage.getItem('userPhoto') || '';
            if (!photo) {
                try { localStorage.setItem('userPhoto', '/nopfp.png'); } catch (e) {}
                photo = '/nopfp.png';
            }
            this.username = name;
            this.user.photo = photo;
        },
        triggerPhotoUpload() {
            if (this.$refs.fileInput) this.$refs.fileInput.click();
        },
        logOut() {
            try {
                const preserve = localStorage.getItem('leftConversations');
                localStorage.clear();
                if (preserve) localStorage.setItem('leftConversations', preserve);
            } catch (e) {}
            try { delete this.$axios.defaults.headers.common['Authorization'] } catch (e) {}
            this.$router.push('/login');
        },
        openEdit() {
            this.editProfile = true;
            this.newUsername = this.username;
        },
        showSuccess(msg) {
            this.successMsg = msg;
            if (this.successTimer) clearTimeout(this.successTimer);
            this.successTimer = setTimeout(() => {
                this.successMsg = null;
            }, 3000);
        },
        async updateUsername() {
            const next = (this.newUsername || '').trim();
            if (!next || next === this.username) { 
                this.errormsg = !next ? 'Username cannot be empty' : 'New username must be different'; 
                return; 
            }
            try {
                await axios.put('/user/username', { name: next });
                this.username = next;
                localStorage.setItem('username', this.username);
                this.newUsername = '';
                this.errormsg = null;
                this.showSuccess('Username updated successfully!');
            } catch (error) {
                console.error('Error updating username:', error);
                const data = error?.response?.data;
                this.errormsg = (typeof data === 'object' ? (data?.error || data?.message) : data) || 'Failed to update username.';
            }
        },
        handlePhotoUpload(e) {
            const file = e?.target?.files?.[0];
            if (!file) return;
            if (!file.type.startsWith('image/')) {
                this.errormsg = 'Please upload a valid image file.';
                if (e?.target) e.target.value = '';
                return;
            }
            const max = 10 * 1024 * 1024;
            if (file.size > max) {
                this.errormsg = 'Image too large (max 10MB)';
                if (e?.target) e.target.value = '';
                return;
            }
            const reader = new FileReader();
            reader.onload = async () => {
                this.newPhoto = reader.result;
                this.errormsg = null;
                await this.updatePhoto();
            };
            reader.onerror = () => {
                this.errormsg = 'Error reading file.';
                if (e?.target) e.target.value = '';
            };
            reader.readAsDataURL(file);
        },
        async updatePhoto() {
            if (!this.newPhoto) return;
            try {
                let b64 = this.newPhoto;
                if (b64.includes(',')) {
                    b64 = b64.split(',')[1];
                }

                await axios.put('/user/photo', { photo: b64 });
                this.user.photo = this.newPhoto;
                localStorage.setItem('userPhoto', this.newPhoto);
                this.newPhoto = null;
                this.errormsg = null;
                if (this.$refs.fileInput) this.$refs.fileInput.value = '';
                this.showSuccess('Photo updated successfully!');
            } catch (error) {
                console.error('Error updating photo:', error);
                this.errormsg = error.response?.data?.message || error.response?.data || 'Failed to update photo.';
                this.newPhoto = null;
                if (this.$refs.fileInput) this.$refs.fileInput.value = '';
            }
        },
    },
    mounted() {
        this.initFromLocal();
    },
}
</script>

<style scoped>
.profile.centered {
  min-height: 90vh;
  background: var(--bg);
  color: var(--text);
  display: grid;
  place-items: center;
  padding: 2rem 1rem;
}
.profile-card {
  background: var(--bg-alt);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 1.5rem;
  width: min(520px, 92vw);
  box-shadow: 0 10px 30px rgba(0,0,0,.35);
  display: grid;
  gap: 1rem;
  justify-items: center;
  position: relative;
}
.photo-wrap {
  display: grid;
  place-items: center;
  position: relative;
}
.photo-wrap.clickable {
  cursor: pointer;
}
.photo-wrap.clickable:hover .pfp {
  filter: brightness(0.7);
}
.photo-overlay {
  position: absolute;
  color: #fff;
  font-size: 0.9rem;
  font-weight: 600;
  pointer-events: none;
  text-shadow: 0 1px 3px rgba(0,0,0,.8);
}
.pfp {
  width: 140px; height: 140px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--accent-alt);
  box-shadow: 0 0 10px var(--accent-alt);
  background: var(--bg);
  transition: filter 0.2s;
}
.username {
  margin: 0;
  font-size: clamp(1.1rem, 1rem + .8vw, 1.6rem);
  font-weight: 800;
}
.user-description {
  margin: 0.5rem 0 0 0;
  font-size: 0.95rem;
  color: var(--text-dim);
  font-style: italic;
}
.controls { display: grid; gap: .5rem; width: 100%; justify-items: center; }

.editor { display: grid; gap: .6rem; width: 100%; }
.field-row { display: flex; gap: .5rem; flex-wrap: wrap; justify-content: center; }
.username-edit {
  display: flex;
  align-items: stretch;
  gap: .3rem;
}
.username-prefix {
  display: flex;
  align-items: center;
  padding: 0 .5rem;
  color: var(--text);
  font-weight: 600;
  font-size: 1rem;
}
.username-input {
  flex: 1 1 150px !important;
}
.input {
  flex: 1 1 240px;
  padding: .6rem .7rem;
  border-radius: var(--radius);
  background: var(--bg-hover);
  border: 1px solid var(--border);
  color: var(--text);
  transition: .15s ease;
}
.input::placeholder { color: var(--text-dim); }
.input:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 8px rgba(26, 115, 232, 0.2);
  background: var(--bg-hover);
}
.textarea {
  resize: vertical;
  min-height: 80px;
  flex: 1 1 100%;
}

.btn {
  appearance: none;
  background: color-mix(in oklab, var(--bg) 85%, var(--accent) 15%);
  border: 1px solid var(--border);
  color: var(--text);
  padding: .6rem .9rem;
  border-radius: var(--radius);
  cursor: pointer;
  transition: background .15s ease, transform .06s ease, box-shadow .15s ease;
}
.btn:hover { background: var(--bg-hover); transform: translateY(-1px); }
.btn:active { transform: translateY(0); }
.btn:disabled { opacity: .6; cursor: not-allowed; }

.btn-primary {
  background: linear-gradient(90deg, var(--accent), var(--accent-alt));
  color: #ffffff;
  border: 1px solid color-mix(in oklab, var(--accent) 60%, var(--accent-alt) 40%);
  box-shadow: 0 2px 8px rgba(26, 115, 232, 0.3);
}
.btn-primary:hover { filter: saturate(1.05) brightness(1.03); }

.btn-edit {
  background: var(--bg);
  border: 1.5px solid var(--neon-cyan);
  color: var(--neon-cyan);
}
.btn-edit:hover {
  background: var(--bg-hover);
}

.btn-logout {
  background: var(--bg);
  border: 1.5px solid #ef4444;
  color: #ef4444;
}
.btn-logout:hover {
  background: var(--bg-hover);
}

.btn-update {
  background: var(--bg);
  border: 1.5px solid var(--neon-cyan);
  color: var(--neon-cyan);
}
.btn-update:hover {
  background: var(--bg-hover);
}

.btn-cancel {
  background: var(--bg);
  border: 1.5px solid var(--border-color);
  color: var(--text-primary);
}
.btn-cancel:hover {
  background: var(--bg-hover);
}
.btn-back {
  position: absolute;
  top: 1rem;
  left: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: rgba(26, 115, 232, 0.08);
  border: 1.5px solid rgba(26, 115, 232, 0.25);
  color: var(--neon-cyan);
  font-size: 0.875rem;
  font-weight: 600;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.btn-back:hover {
  background: rgba(26, 115, 232, 0.15);
  border-color: var(--neon-cyan);
  box-shadow: 0 0 10px rgba(26, 115, 232, 0.2);
  transform: translateX(-3px);
}

.btn-back svg {
  stroke: var(--neon-cyan);
}

.success-message {
  position: absolute;
  bottom: -80px;
  left: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 0.75rem 1rem;
  background: rgba(34, 197, 94, 0.15);
  border: 1px solid #22c55e;
  color: #22c55e;
  border-radius: var(--radius);
  font-weight: 600;
  font-size: 0.95rem;
  text-align: center;
  width: calc(100% + 2rem);
  animation: slideIn 0.3s ease;
  transform: translate(-50%, 0);
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translate(-50%, -8px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}
</style>
