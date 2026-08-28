<template>
  <section class="group-create-neon">

    <header class="page-header-neon">
      <h1 class="title-neon">Create New Group</h1>
      <p class="subtitle-neon">Bring people together in one place</p>
    </header>

    <!-- Main Card -->
    <div class="card-container-neon">
      <div class="card-neon">
        <button class="close-btn-neon" @click="goHome" title="Back to Home">
          ×
        </button>
        <div v-if="errormsg" class="error-banner-neon">
          {{ errormsg }}
          <button @click="errormsg = null" class="error-close-neon">×</button>
        </div>

        <!-- Photo Selector -->
        <div class="photo-section-neon">
          <div class="photo-selector-neon" @click="$refs.photoInput?.click()">
            <input 
              ref="photoInput" 
              type="file" 
              accept="image/*" 
              @change="handlePhoto" 
              class="file-input-hidden" 
            />
            <div class="photo-wrap-neon">
              <div 
                class="photo-avatar-neon" 
                :style="{ backgroundImage: 'url(' + photoPreview + ')' }"
              >
                <div class="photo-overlay-neon">
                  <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"></path>
                    <circle cx="12" cy="13" r="4"></circle>
                  </svg>
                  <span>Click to upload</span>
                </div>
              </div>
            </div>
          </div>
          <p class="photo-hint-neon">Group photo (optional)</p>
        </div>

        <!-- Group Name -->
        <div class="field-neon">
          <label class="label-neon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
            Group Name
          </label>
          <div class="input-container-neon">
            <input
              v-model="groupName"
              type="text"
              class="input-neon"
              :class="{ invalid: nameError }"
              placeholder="e.g., Weekend Crew, Study Group..."
              maxlength="64"
            />
            <span class="char-count-neon">{{ groupName.length }}/64</span>
          </div>
          <p v-if="nameError" class="error-hint-neon">{{ nameError }}</p>
        </div>

        <!-- Add Members -->
        <div class="field-neon">
          <label class="label-neon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <line x1="19" y1="8" x2="19" y2="14"></line>
              <line x1="22" y1="11" x2="16" y2="11"></line>
            </svg>
            Add Members
          </label>
          
          <div class="search-container-neon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="search-icon-neon">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="m21 21-4.35-4.35"></path>
            </svg>
            <input
              v-model="userSearch"
              class="input-neon search"
              placeholder="Search users by username..."
              @input="searchUsers"
            />
          </div>

          <!-- Search Status -->
          <div v-if="searching || searchError || (userSearch && !searchResults.length)" class="search-status-neon">
            <div v-if="searching" class="status-loading-neon">
              <div class="spinner-small-neon"></div>
              <span>Searching...</span>
            </div>
            <div v-else-if="searchError" class="status-error-neon">{{ searchError }}</div>
            <div v-else-if="userSearch && !searchResults.length" class="status-empty-neon">
              No users found
            </div>
          </div>

          <!-- Search Results -->
          <div v-if="searchResults.length" class="search-results-neon">
            <div 
              v-for="user in searchResults" 
              :key="user.id" 
              class="user-item-neon"
            >
              <div class="user-avatar-neon-wrapper">
              <img 
                :src="getUserPhoto(user)" 
                :alt="user.username"
                class="user-avatar-img-neon"
                />
              </div>
              <div class="user-info-neon">
                <span class="username-neon">{{ user.username }}</span>
              </div>
              <button 
                type="button" 
                class="btn-add-neon"
                @click="addMember(user)"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <line x1="12" y1="5" x2="12" y2="19"></line>
                  <line x1="5" y1="12" x2="19" y2="12"></line>
                </svg>
                Add
              </button>
            </div>
          </div>

         <!-- Selected Members -->
<div v-if="members.length" class="members-section-neon">
  <div class="members-header-neon">
    <span>Selected Members ({{ members.length }})</span>
  </div>
  <div class="members-chips-neon">
    <div 
      v-for="member in members" 
      :key="member.id" 
      class="member-chip-neon"
    >
      <div class="chip-avatar-wrapper-neon">
        <img 
          :src="getUserPhoto(member)" 
          :alt="member.username"
          class="chip-avatar-img-neon"
        />
      </div>
      <span class="chip-name-neon">{{ member.username }}</span>
      <button 
        type="button" 
        class="chip-remove-neon"
        @click="removeMember(member)"
      >
        ×
      </button>
    </div>
  </div>
</div>
          <p class="hint-neon">
            💡 You need at least one member to create a group
          </p>
        </div>

        <!-- Actions -->
        <div class="actions-neon">
          <button
            class="btn-neon primary"
            @click="createGroup"
            :disabled="!canCreate || isSubmitting"
          >
            <svg v-if="!isSubmitting" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <div v-else class="spinner-small-neon"></div>
            {{ isSubmitting ? 'Creating Group...' : 'Create Group' }}
          </button>
        </div>

        <!-- Success Message -->
        <div v-if="success" class="success-banner-neon">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          Group created successfully!
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: "GroupCreateView",
  
  data() {
    return {
      groupName: "",
      groupPhoto: "",
      photoName: "",
      userSearch: "",
      searchResults: [],
      searchError: null,
      searching: false,
      searchTimer: null,
      members: [],
      errormsg: null,
      success: false,
      isSubmitting: false,
      currentUserId: localStorage.getItem('userId') || '',
      nameError: '',
    };
  },
  
  computed: {
    canCreate() {
      return this.groupName.trim().length >= 1 && this.members.length > 0;
    },
    
    photoPreview() {
      if (this.groupPhoto) {
        return `data:image/png;base64,${this.groupPhoto}`;
      }
      return this.letterAvatar(this.groupName || 'G');
    }
  },
  
  methods: {
    goHome() {
    this.$router.push('/home')
    },
    getUserPhoto(user) {
    if (!user) return this.letterAvatar('?')
    
    const photoB64 = user.photo || user.Photo || user.profilePhoto || user.ProfilePhoto
    
    if (photoB64) {
      if (photoB64.startsWith('data:')) return photoB64
      return `data:image/png;base64,${photoB64}`
    }
    
    const username = user.username || user.Username || user.name || user.Name || 'User'
    return this.letterAvatar(username)
  },
  
    
    letterAvatar(name, size = 128) {
      try {
        const letter = (String(name || '').trim().charAt(0) || '?').toUpperCase();
        const bg = '#e8f0fe';
        const fg = 'var(--neon-cyan)';
        const svg = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns='http://www.w3.org/2000/svg' width='${size}' height='${size}' viewBox='0 0 ${size} ${size}'>
  <defs>
    <linearGradient id='grad-${name}' x1='0%' y1='0%' x2='100%' y2='100%'>
      <stop offset='0%' stop-color='var(--neon-cyan)'/>
      <stop offset='100%' stop-color='var(--neon-purple)'/>
    </linearGradient>
  </defs>
  <rect width='100%' height='100%' rx='${Math.floor(size/2)}' ry='${Math.floor(size/2)}' fill='url(#grad-${name})'/>
  <text x='50%' y='53%' dominant-baseline='middle' text-anchor='middle' fill='${bg}' font-family='Segoe UI, Roboto, sans-serif' font-weight='800' font-size='${Math.floor(size*0.48)}'>${letter}</text>
</svg>`;
        return 'data:image/svg+xml;utf8,' + encodeURIComponent(svg);
      } catch {
        return '/nopfp.png';
      }
    },
    
    getInitial(name) {
      if (!name || typeof name !== 'string') return '?';
      return name.trim().charAt(0).toUpperCase();
    },
    
    handlePhoto(e) {
      const file = e?.target?.files?.[0];
      if (!file) { 
        this.groupPhoto = "";
        this.photoName = "";
        return;
      }
      
      if (!file.type.startsWith('image/')) {
        this.errormsg = 'Please select an image file';
        return;
      }
      
      const max = 10 * 1024 * 1024;
      if (file.size > max) {
        this.errormsg = 'Image too large (max 10MB)';
        if (e?.target) e.target.value = '';
        return;
      }
      
      const reader = new FileReader();
      reader.onload = () => {
        const result = reader.result || '';
        const b64 = typeof result === 'string' && result.includes(',') 
          ? result.split(',')[1] 
          : result;
        this.groupPhoto = b64;
        this.photoName = file.name || '';
      };
      reader.onerror = () => {
        this.errormsg = 'Failed to read file';
      };
      reader.readAsDataURL(file);
    },
    
    async searchUsers() {
      const q = (this.userSearch || '').trim();
      
      if (this.searchTimer) clearTimeout(this.searchTimer);
      
      if (!q) {
        this.searchResults = [];
        this.searchError = null;
        this.searching = false;
        return;
      }
      
      this.searchError = null;
      this.searching = true;
      
      this.searchTimer = setTimeout(async () => {
        try {
          const token = localStorage.getItem('token');
          const res = await this.$axios.get(`/searchby?user=${encodeURIComponent(q)}`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          
          const all = res.data?.users || [];
          this.searchResults = all.filter(u => 
            u.id !== this.currentUserId && 
            !this.members.find(m => m.id === u.id)
          );
        } catch (e) {
          this.searchError = "Failed to search users";
        } finally {
          this.searching = false;
        }
      }, 300);
    },
    
    addMember(user) {
      if (user?.id === this.currentUserId) return;
      
      if (!this.members.find(m => m.id === user.id)) {
        this.members.push(user);
        // Remove from search results
        this.searchResults = this.searchResults.filter(u => u.id !== user.id);
      }
    },
    
    removeMember(user) {
      this.members = this.members.filter(m => m.id !== user.id);
      // Re-search to show removed user again
      if (this.userSearch) {
        this.searchUsers();
      }
    },
    
    async createGroup() {
      this.errormsg = null;
      this.success = false;
      this.isSubmitting = true;
      this.nameError = '';
      
      if (this.groupName.trim().length === 0) {
        this.nameError = 'Group name is required';
        this.isSubmitting = false;
        return;
      }
      
      if (this.members.length === 0) {
        this.errormsg = 'Please add at least one member';
        this.isSubmitting = false;
        return;
      }
      
      try {
        const token = localStorage.getItem('token');
        const payload = {
          groupName: this.groupName,
          members: this.members.map(m => m.id).filter(id => id && id !== this.currentUserId),
          groupPhoto: this.groupPhoto || undefined,
        };
        
        const res = await this.$axios.post('/groups', payload, {
          headers: { Authorization: `Bearer ${token}` }
        });
        
        this.success = true;
        
        // Reset form
        setTimeout(() => {
          const convId = res?.data?.id || res?.data?.conversationId;
          if (convId) {
            this.$router.push(`/home?conv=${convId}`);
          } else {
            this.$router.push('/home');
          }
        }, 1000);
        
      } catch (e) {
        console.error('Failed to create group:', e);
        this.errormsg = e.response?.data?.error || "Failed to create group";
      } finally {
        this.isSubmitting = false;
      }
    }
  }
};
</script>

<style scoped>
/* Container */
.group-create-neon {
  position: relative;
  min-height: 100vh;
  background: var(--bg-primary);
  color: var(--text-primary);
  padding: 2rem 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow-x: hidden;
}

/* Header */
.page-header-neon {
  position: relative;
  z-index: 10;
  text-align: center;
  margin-bottom: 2rem;
}

.title-neon {
  margin: 0 0 0.5rem 0;
  font-size: clamp(1.75rem, 2vw + 1rem, 2.5rem);
  font-weight: 800;
  background: linear-gradient(90deg, var(--neon-cyan) 0%, var(--neon-purple) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 0 40px rgba(26, 115, 232, 0.25);
}

.subtitle-neon {
  margin: 0;
  font-size: 1rem;
  color: var(--text-secondary);
  font-weight: 500;
}

/* Card Container */
.card-container-neon {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 600px;
}

.card-neon {
  background: #ffffff;
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-radius: 20px;
  padding: 1rem;
  box-shadow: 0 8px 32px rgba(26, 115, 232, 0.06), 0 0 60px rgba(26, 115, 232, 0.12);
  backdrop-filter: blur(12px);
  max-height: 77vh; 
  overflow-y: auto;
  overflow-x: hidden;
}
.card-neon::-webkit-scrollbar {
  width: 8px;
}

.card-neon::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.06);
  border-radius: 10px;
}

.card-neon::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, var(--neon-cyan), var(--neon-purple));
  border-radius: 10px;
  border: 2px solid rgba(0, 0, 0, 0.06);
}

.card-neon::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, var(--neon-purple), var(--neon-cyan));
}

/* Error Banner */
.error-banner-neon {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.875rem 1rem;
  margin-bottom: 1.5rem;
  background: rgba(255, 0, 100, 0.15);
  border: 1px solid rgba(255, 0, 100, 0.4);
  border-radius: 10px;
  color: #ff0064;
  font-size: 0.925rem;
}

.error-close-neon {
  background: transparent;
  border: none;
  color: #ff0064;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  transition: transform 0.2s ease;
}

.error-close-neon:hover {
  transform: rotate(90deg);
}

/* Photo Section */
.photo-section-neon {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 2rem;
}

.photo-selector-neon {
  cursor: pointer;
}

.photo-wrap-neon {
  position: relative;
}

.photo-avatar-neon {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  background-size: cover;
  background-position: center;
  border: 3px solid var(--neon-cyan);
  box-shadow: 0 0 30px rgba(26, 115, 232, 0.30);
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
}

.photo-avatar-neon:hover {
  transform: scale(1.05);
  box-shadow: 0 0 40px rgba(26, 115, 232, 0.45);
}

.photo-overlay-neon {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  background: rgba(26, 115, 232, 0.06);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.photo-avatar-neon:hover .photo-overlay-neon {
  opacity: 1;
}

.photo-overlay-neon svg {
  color: var(--neon-cyan);
}

.photo-overlay-neon span {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--neon-cyan);
}

.photo-hint-neon {
  margin-top: 0.75rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.file-input-hidden {
  display: none;
}

/* Fields */
.field-neon {
  margin-bottom: 1.75rem;
}

.label-neon {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--neon-cyan);
}

.label-neon svg {
  color: var(--neon-cyan);
}

.input-container-neon {
  position: relative;
  display: flex;
  align-items: center;
}

.input-neon {
  width: 100%;
  padding: 0.875rem 1rem;
  background: rgba(26, 115, 232, 0.06);
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-radius: 12px;
  color: var(--text-primary);
  font-size: 0.95rem;
  transition: all 0.3s ease;
  outline: none;
}

.input-neon::placeholder {
  color: var(--text-dim);
}

.input-neon:focus {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.25);
  background: #ffffff;
}

.input-neon.invalid {
  border-color: #ff0064;
  box-shadow: 0 0 15px rgba(255, 0, 100, 0.3);
}

.char-count-neon {
  position: absolute;
  right: 1rem;
  font-size: 0.8rem;
  color: var(--text-dim);
}

.error-hint-neon {
  margin-top: 0.5rem;
  font-size: 0.85rem;
  color: #ff0064;
}

/* Search Container */
.search-container-neon {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon-neon {
  position: absolute;
  left: 1rem;
  color: rgba(26, 115, 232, 0.45);
  pointer-events: none;
}

.input-neon.search {
  padding-left: 3rem;
}

/* Search Status */
.search-status-neon {
  margin-top: 0.75rem;
  text-align: center;
}

.status-loading-neon {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.status-error-neon {
  color: #ff0064;
  font-size: 0.9rem;
}

.status-empty-neon {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.spinner-small-neon {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-top-color: var(--neon-cyan);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Search Results */
.search-results-neon {
  margin-top: 1rem;
  max-height: 300px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.search-results-neon::-webkit-scrollbar {
  width: 6px;
}

.search-results-neon::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.06);
}

.search-results-neon::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, var(--neon-cyan), var(--neon-purple));
  border-radius: 10px;
}

.user-item-neon {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 0.75rem;
  background: rgba(26, 115, 232, 0.06);
  border: 1px solid rgba(26, 115, 232, 0.15);
  border-radius: 12px;
  transition: all 0.2s ease;
}
.user-avatar-neon-wrapper {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid rgba(26, 115, 232, 0.40);
  box-shadow: 0 0 10px rgba(26, 115, 232, 0.25);
  flex-shrink: 0;
}

.user-avatar-img-neon {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.user-item-neon:hover {
  background: rgba(26, 115, 232, 0.10);
  border-color: rgba(26, 115, 232, 0.30);
}

.user-info-neon {
  flex: 1;
  min-width: 0;
}

.username-neon {
  font-weight: 600;
  color: var(--text-primary);
}

.btn-add-neon {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.5rem 1rem;
  background: rgba(26, 115, 232, 0.12);
  border: 1px solid rgba(26, 115, 232, 0.30);
  border-radius: 8px;
  color: var(--neon-cyan);
  font-weight: 600;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add-neon:hover {
  background: rgba(26, 115, 232, 0.20);
  box-shadow: 0 0 15px rgba(26, 115, 232, 0.25);
}

/* Members Section */
.members-section-neon {
  margin-top: 1.25rem;
}

.members-section-neon {
  margin-top: 1.25rem;
  max-height: 250px; 
  overflow-y: auto;  
  overflow-x: hidden;
  padding-right: 0.25rem; 
}

.members-section-neon::-webkit-scrollbar {
  width: 6px;
}

.members-section-neon::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.06);
  border-radius: 10px;
}

.members-section-neon::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, var(--neon-cyan), var(--neon-purple));
  border-radius: 10px;
}

.members-section-neon::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, var(--neon-purple), var(--neon-cyan));
}

.members-header-neon {
  margin-bottom: 0.75rem;
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--text-secondary);
  position: sticky;
  top: 0;
  background: #ffffff;
  padding: 0.25rem 0;
  z-index: 1;
}
.chip-avatar-wrapper-neon {
  width: 34px; 
  height: 34px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid rgba(138, 43, 226, 0.5);
  box-shadow: 0 0 8px rgba(26, 115, 232, 0.20);
  flex-shrink: 0;
}

.chip-avatar-img-neon {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}


.members-chips-neon {
  display: flex;
  flex-wrap: wrap;
  gap: 0.65rem; 
}


.member-chip-neon {
  display: flex;
  align-items: center;
  gap: 0.65rem; 
  padding: 0.65rem 0.85rem; 
  background: rgba(26, 115, 232, 0.08);
  border: 1px solid rgba(26, 115, 232, 0.20);
  border-radius: 24px; 
  transition: all 0.2s ease;
}

.member-chip-neon:hover {
  background: rgba(26, 115, 232, 0.10);
  border-color: rgba(138, 43, 226, 0.5);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(26, 115, 232, 0.20);
}

.chip-name-neon {
  font-size: 0.925rem; 
  font-weight: 600;
  color: var(--text-primary);
}

.chip-remove-neon {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 1.4rem; 
  cursor: pointer;
  padding: 0;
  line-height: 1;
  transition: all 0.2s ease;
  margin-left: 0.15rem;
}

.chip-remove-neon:hover {
  color: #ff0064;
  transform: rotate(90deg) scale(1.1);
}

.hint-neon {
  margin-top: 0.75rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
  text-align: center;
}

/* Actions */
.actions-neon {
  margin-top: 2rem;
  display: flex;
  justify-content: center;
}

.btn-neon {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.875rem 2rem;
  border-radius: 12px;
  font-weight: 700;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.btn-neon.primary {
  background: var(--neon-cyan);
  color: var(--bg-primary);
  box-shadow: 0 0 25px rgba(26, 115, 232, 0.30);
}

.btn-neon.primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 0 35px rgba(26, 115, 232, 0.45);
}

.btn-neon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

/* Success Banner */
.success-banner-neon {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 1rem;
  padding: 0.875rem 1rem;
  background: rgba(26, 115, 232, 0.12);
  border: 1px solid rgba(26, 115, 232, 0.30);
  border-radius: 10px;
  color: var(--neon-cyan);
  font-weight: 600;
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
.close-btn-neon {
  position: sticky;
  top: 0;
  margin-left: auto;
  margin-bottom: 0.75rem;
  display: flex;
  background: transparent;
  border: none;

  color: #ff0064;
  font-size: 1.8rem;
  font-weight: 700;
  line-height: 1;

  cursor: pointer;
  z-index: 5;

  transition: transform 0.2s ease, text-shadow 0.2s ease;
}

.close-btn-neon:hover {
  transform: rotate(90deg) scale(1.2);
  text-shadow: 0 0 12px rgba(255, 0, 100, 0.8);
}

</style>
