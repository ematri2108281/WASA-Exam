<template>
  <section class="group-edit-neon">
    <div class="particles-bg">
      <div class="particle-neon" v-for="n in 15" :key="n" :style="getParticleStyle(n)"></div>
    </div>

    <header class="edit-header-neon">
      <button v-if="!embedded" class="back-btn-neon" @click="goBack">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"></line>
          <polyline points="12 19 5 12 12 5"></polyline>
        </svg>
        Back
      </button>
      <h1 class="title-neon">Edit Group</h1>
      <button v-if="embedded" class="close-btn-neon" @click="closeModal">×</button>
    </header>

    <div v-if="errorMessage" class="error-banner-neon">
      {{ errorMessage }}
      <button @click="errorMessage = null" class="banner-close-neon">×</button>
    </div>
    
    <div v-if="successMessage" class="success-banner-neon">
      {{ successMessage }}
    </div>

    <div v-if="loading" class="loading-container-neon">
      <div class="spinner-neon"></div>
      <p>Loading group info...</p>
    </div>

    <div v-else class="card-neon">
        <div class="photo-section-neon">
          <div class="photo-selector-neon" @click="triggerPhotoInput">
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
                  <span>Change Photo</span>
                </div>
              </div>
            </div>
          </div>
          <button 
            v-if="photoChanged" 
            class="btn-save-photo-neon"
            @click="savePhoto"
            :disabled="savingPhoto"
          >
            <svg v-if="!savingPhoto" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <div v-else class="spinner-tiny-neon"></div>
            {{ savingPhoto ? 'Saving...' : 'Save Photo' }}
          </button>
        </div>

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
          <div class="input-row-neon">
            <input
              v-model="groupName"
              type="text"
              class="input-neon"
              placeholder="Group name..."
              maxlength="64"
            />
            <button 
              v-if="nameChanged" 
              class="btn-save-inline-neon"
              @click="saveName"
              :disabled="savingName || !groupName.trim()"
            >
              <svg v-if="!savingName" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              <div v-else class="spinner-tiny-neon"></div>
            </button>
          </div>
          <span class="char-count-neon">{{ groupName.length }}/64</span>
        </div>

        <div class="members-section-neon">
          <div class="section-header-neon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
            <span>Members ({{ members.length }})</span>
          </div>

          <div class="members-list-neon">
            <div 
              v-for="member in members" 
              :key="member.id"
              class="member-item-neon"
            >
              <div class="member-avatar-wrapper-neon">
                <img 
                  :src="getMemberPhoto(member)" 
                  :alt="member.username"
                  class="member-avatar-img-neon"
                />
              </div>
              <div class="member-info-neon">
                <span class="member-name-neon">{{ member.username }}</span>
                <span v-if="isCreator(member)" class="creator-badge-neon">Creator</span>
              </div>
              <button 
                v-if="!isMe(member) && !isCreator(member)"
                class="btn-remove-member-neon"
                @click="removeMemberConfirm(member)"
                :disabled="removingMember === member.id"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <div class="add-section-neon">
          <div class="section-header-neon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <line x1="19" y1="8" x2="19" y2="14"></line>
              <line x1="22" y1="11" x2="16" y2="11"></line>
            </svg>
            <span>Add Members</span>
          </div>

          <div class="search-container-neon">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="search-icon-neon">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="m21 21-4.35-4.35"></path>
            </svg>
            <input
              v-model="searchQuery"
              class="input-neon search"
              placeholder="Search users..."
              @input="searchUsers"
            />
          </div>

          <div v-if="searching || searchError || (searchQuery && !searchResults.length)" class="search-status-neon">
            <div v-if="searching" class="status-loading-neon">
              <div class="spinner-small-neon"></div>
              <span>Searching...</span>
            </div>
            <div v-else-if="searchError" class="status-error-neon">{{ searchError }}</div>
            <div v-else-if="searchQuery && !searchResults.length" class="status-empty-neon">
              No users found
            </div>
          </div>

          <div v-if="searchResults.length" class="search-results-neon">
            <div 
              v-for="user in searchResults" 
              :key="user.id" 
              class="user-item-neon"
            >
              <div class="user-avatar-wrapper-neon">
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
                :disabled="addingMember === user.id"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <line x1="12" y1="5" x2="12" y2="19"></line>
                  <line x1="5" y1="12" x2="19" y2="12"></line>
                </svg>
                Add
              </button>
            </div>
          </div>
        </div>
      </div>
    <div v-if="removeMemberModal.open" class="modal-overlay-neon" @click.self="closeRemoveModal">
      <div class="modal-card-neon">
        <h3 class="modal-title-neon">Remove Member</h3>
        <p class="modal-text-neon">
          Are you sure you want to remove <strong>{{ removeMemberModal.username }}</strong> from this group?
        </p>
        <div class="modal-actions-neon">
          <button class="btn-neon danger" @click="confirmRemoveMember">Remove</button>
          <button class="btn-neon secondary" @click="closeRemoveModal">Cancel</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'GroupEditView',
  props: {
    groupIdProp: { type: String, default: null },
    embedded: { type: Boolean, default: false }
  },
  data: function() {
    return {
      loading: true,
      errorMessage: null,
      successMessage: null,
      groupId: null,
      groupName: '',
      originalName: '',
      groupPhoto: '',
      originalPhoto: '',
      members: [],
      creatorId: null,
      photoFile: null,
      photoChanged: false,
      savingPhoto: false,
      savingName: false,
      searchQuery: '',
      searchResults: [],
      searching: false,
      searchError: null,
      searchTimer: null,
      addingMember: null,
      removingMember: null,
      removeMemberModal: { open: false, userId: null, username: '' },
      userId: localStorage.getItem('userId') || ''
    }
  },
  computed: {
    nameChanged: function() {
      return this.groupName.trim() !== this.originalName && this.groupName.trim().length > 0
    },
    photoPreview: function() {
      if (this.groupPhoto && this.groupPhoto.length > 0) {
        if (this.groupPhoto.indexOf('data:') === 0) {
          return this.groupPhoto
        }
        return 'data:image/png;base64,' + this.groupPhoto
      }
      return this.letterAvatar(this.groupName || 'Group')
    }
  },
  methods: {
    triggerPhotoInput: function() {
      if (this.$refs.photoInput) {
        this.$refs.photoInput.click()
      }
    },
    getParticleStyle: function(n) {
      var x = Math.random() * 100
      var delay = Math.random() * 5
      var duration = 20 + Math.random() * 30
      var size = 1 + Math.random() * 2
      return {
        left: x + '%',
        animationDelay: delay + 's',
        animationDuration: duration + 's',
        width: size + 'px',
        height: size + 'px',
        opacity: Math.random() * 0.3 + 0.1
      }
    },
    letterAvatar: function(name, size) {
      if (!size) size = 128
      try {
        var letter = (String(name || '').trim().charAt(0) || '?').toUpperCase()
        var svg = '<?xml version="1.0" encoding="UTF-8"?><svg xmlns="http://www.w3.org/2000/svg" width="' + size + '" height="' + size + '" viewBox="0 0 ' + size + ' ' + size + '"><defs><linearGradient id="grad-' + name + '" x1="0%" y1="0%" x2="100%" y2="100%"><stop offset="0%" stop-color="var(--neon-cyan)"/><stop offset="100%" stop-color="var(--neon-purple)"/></linearGradient></defs><rect width="100%" height="100%" rx="' + Math.floor(size/2) + '" ry="' + Math.floor(size/2) + '" fill="url(#grad-' + name + ')"/><text x="50%" y="53%" dominant-baseline="middle" text-anchor="middle" fill="#ffffff" font-family="Segoe UI, Roboto, sans-serif" font-weight="800" font-size="' + Math.floor(size*0.48) + '">' + letter + '</text></svg>'
        return 'data:image/svg+xml;utf8,' + encodeURIComponent(svg)
      } catch (e) {
        return '/nopfp.png'
      }
    },
    getUserPhoto: function(user) {
      if (!user) return this.letterAvatar('?')
      var photoB64 = user.photo || user.Photo || user.profilePhoto || user.ProfilePhoto
      if (photoB64) {
        if (photoB64.indexOf('data:') === 0) return photoB64
        return 'data:image/png;base64,' + photoB64
      }
      return this.letterAvatar(user.username || 'User')
    },
    getMemberPhoto: function(member) {
      return this.getUserPhoto(member)
    },
    isMe: function(member) {
      return String(member.id || member.userId) === String(this.userId)
    },
    isCreator: function(member) {
      return String(member.id || member.userId) === String(this.creatorId)
    },
    load: function() {
      var self = this
      self.loading = true
      self.errorMessage = null
      var token = localStorage.getItem('token')
      var config = token ? { headers: { Authorization: 'Bearer ' + token } } : {}

      Promise.all([
        self.$axios.get('/conversations/' + self.groupId, config),
        self.$axios.get('/conversations/' + self.groupId + '/members', config)
      ])
        .then(function(results) {
          var conv = results[0].data || {}
          var members = results[1].data || []

          self.groupName = conv.name || ''
          self.originalName = conv.name || ''
          self.groupPhoto = conv.groupPhoto || conv.profilePhoto || conv.photo || ''
          self.originalPhoto = self.groupPhoto
          self.members = members
          self.creatorId = conv.createdBy || conv.creatorId || null
          self.loading = false
        })
        .catch(function(error) {
          console.error('Failed to load group:', error)
          self.errorMessage = 'Failed to load group information'
          self.loading = false
        })
    },
    handlePhoto: function(e) {
      var file = e && e.target && e.target.files && e.target.files[0]
      if (!file) {
        this.photoFile = null
        this.photoChanged = false
        return
      }
      if (file.type.indexOf('image/') !== 0) {
        this.errorMessage = 'Please select an image file'
        return
      }
      var max = 10 * 1024 * 1024
      if (file.size > max) {
        this.errorMessage = 'Image too large (max 10MB)'
        if (e && e.target) e.target.value = ''
        return
      }
      var self = this
      var reader = new FileReader()
      reader.onload = function() {
        var result = reader.result || ''
        var b64 = typeof result === 'string' && result.indexOf(',') > -1 ? result.split(',')[1] : result
        self.groupPhoto = b64
        self.photoFile = b64
        self.photoChanged = true
      }
      reader.onerror = function() {
        self.errorMessage = 'Failed to read file'
      }
      reader.readAsDataURL(file)
    },
    savePhoto: function() {
      if (!this.photoFile) return
      var self = this
      self.savingPhoto = true
      self.errorMessage = null
      self.successMessage = null
      var token = localStorage.getItem('token')
      self.$axios.put('/groups/' + self.groupId + '/photo', { photo: self.photoFile }, token ? { headers: { Authorization: 'Bearer ' + token } } : {})
        .then(function() {
          self.originalPhoto = self.groupPhoto
          self.photoChanged = false
          self.successMessage = 'Photo updated successfully!'
          setTimeout(function() { self.successMessage = null }, 3000)
          self.$emit('updated', { convId: self.groupId })
          self.savingPhoto = false
        })
        .catch(function(error) {
          console.error('Failed to update photo:', error)
          self.errorMessage = error.response && error.response.data && error.response.data.error || 'Failed to update photo'
          self.savingPhoto = false
        })
    },
    saveName: function() {
      if (!this.groupName.trim()) return
      var self = this
      self.savingName = true
      self.errorMessage = null
      self.successMessage = null
      var token = localStorage.getItem('token')
      self.$axios.put('/groups/' + self.groupId + '/name', { name: self.groupName.trim() }, token ? { headers: { Authorization: 'Bearer ' + token } } : {})
        .then(function() {
          self.originalName = self.groupName.trim()
          self.successMessage = 'Name updated successfully!'
          setTimeout(function() { self.successMessage = null }, 3000)
          self.$emit('updated', { convId: self.groupId })
          self.savingName = false
        })
        .catch(function(error) {
          console.error('Failed to update name:', error)
          self.errorMessage = error.response && error.response.data && error.response.data.error || 'Failed to update name'
          self.savingName = false
        })
    },
    searchUsers: function() {
      var self = this
      var q = (self.searchQuery || '').trim()
      if (self.searchTimer) clearTimeout(self.searchTimer)
      if (!q) {
        self.searchResults = []
        self.searchError = null
        self.searching = false
        return
      }
      self.searchError = null
      self.searching = true
      self.searchTimer = setTimeout(function() {
        var token = localStorage.getItem('token')
        self.$axios.get('/searchby?user=' + encodeURIComponent(q), { headers: { Authorization: 'Bearer ' + token } })
          .then(function(res) {
            var all = res.data && res.data.users || []
            self.searchResults = all.filter(function(u) {
              return u.id !== self.userId && !self.members.find(function(m) { return (m.id || m.userId) === u.id })
            })
            self.searching = false
          })
          .catch(function() {
            self.searchError = 'Failed to search users'
            self.searching = false
          })
      }, 300)
    },
    addMember: function(user) {
      if (!user || !user.id) return
      var self = this
      self.addingMember = user.id
      self.errorMessage = null
      self.successMessage = null
      var token = localStorage.getItem('token')
      self.$axios.post('/groups/' + self.groupId, { conversationId: self.groupId, userId: user.id }, token ? { headers: { Authorization: 'Bearer ' + token } } : {})
        .then(function() {
          self.members.push(user)
          self.searchResults = self.searchResults.filter(function(u) { return u.id !== user.id })
          self.successMessage = user.username + ' added to the group!'
          setTimeout(function() { self.successMessage = null }, 3000)
          self.$emit('updated', { convId: self.groupId })
          self.addingMember = null
        })
        .catch(function(error) {
          console.error('Failed to add member:', error)
          self.errorMessage = error.response && error.response.data && error.response.data.error || 'Failed to add member'
          self.addingMember = null
        })
    },
    removeMemberConfirm: function(member) {
      this.removeMemberModal = {
        open: true,
        userId: member.id || member.userId,
        username: member.username || 'User'
      }
    },
    closeRemoveModal: function() {
      this.removeMemberModal = { open: false, userId: null, username: '' }
    },
    confirmRemoveMember: function() {
      var userId = this.removeMemberModal.userId
      if (!userId) return
      var self = this
      self.removingMember = userId
      self.errorMessage = null
      self.successMessage = null
      self.closeRemoveModal()
      var token = localStorage.getItem('token')
      self.$axios.delete('/groups/' + self.groupId + '?userId=' + userId, token ? { headers: { Authorization: 'Bearer ' + token } } : {})
        .then(function() {
          self.members = self.members.filter(function(m) { return String(m.id || m.userId) !== String(userId) })
          self.successMessage = 'Member removed from the group'
          setTimeout(function() { self.successMessage = null }, 3000)
          self.$emit('updated', { convId: self.groupId })
          self.removingMember = null
        })
        .catch(function(error) {
          console.error('Failed to remove member:', error)
          self.errorMessage = error.response && error.response.data && error.response.data.error || 'Failed to remove member'
          self.removingMember = null
        })
    },
    goBack: function() {
      this.$router.back()
    },
    closeModal: function() {
      this.$emit('close', { updated: false })
    }
  },
  mounted: function() {
    this.groupId = this.groupIdProp || this.$route.params.groupId
    if (!this.groupId) {
      this.errorMessage = 'No group ID provided'
      this.loading = false
      return
    }
    this.load()
  }
}
</script>

<style scoped>
/* Container */
.group-edit-neon {
  position: relative;
  min-height: 100vh;
  background: var(--bg-primary);
  color: var(--text-primary);
  padding: 2rem 1rem 4rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow-x: hidden;
  overflow-y: auto; /* Allow page to scroll naturally */
}

/* Particles Background */
.particles-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}

.particle-neon {
  position: absolute;
  background: linear-gradient(135deg, var(--neon-cyan), var(--neon-purple));
  border-radius: 50%;
  animation: float-particle linear infinite;
}

@keyframes float-particle {
  0% {
    transform: translateY(100vh) rotate(0deg);
    opacity: 0;
  }
  10% {
    opacity: 0.6;
  }
  90% {
    opacity: 0.6;
  }
  100% {
    transform: translateY(-100px) rotate(360deg);
    opacity: 0;
  }
}

/* Header */
.edit-header-neon {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 600px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.back-btn-neon {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  background: rgba(26, 115, 232, 0.08);
  border: 1px solid rgba(26, 115, 232, 0.25);
  border-radius: 999px;
  color: var(--neon-cyan);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.back-btn-neon:hover {
  background: rgba(26, 115, 232, 0.15);
  box-shadow: 0 0 15px rgba(26, 115, 232, 0.25);
}

.title-neon {
  margin: 0;
  font-size: clamp(1.5rem, 2vw + 1rem, 2rem);
  font-weight: 800;
  background: linear-gradient(90deg, var(--neon-cyan) 0%, var(--neon-purple) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.close-btn-neon {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 2rem;
  cursor: pointer;
  line-height: 1;
  transition: all 0.2s ease;
}

.close-btn-neon:hover {
  color: #ff0064;
  transform: rotate(90deg);
}

/* Card */
.card-neon {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 600px;
  background: #ffffff;
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(26, 115, 232, 0.06), 0 0 60px rgba(26, 115, 232, 0.12);
  backdrop-filter: blur(12px);
  /* Removed max-height and overflow - single card now */
}

/* Banners */
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

.success-banner-neon {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.875rem 1rem;
  margin-bottom: 1.5rem;
  background: rgba(26, 115, 232, 0.12);
  border: 1px solid rgba(26, 115, 232, 0.30);
  border-radius: 10px;
  color: var(--neon-cyan);
  font-weight: 600;
}

.banner-close-neon {
  background: transparent;
  border: none;
  color: currentColor;
  font-size: 1.5rem;
  cursor: pointer;
  line-height: 1;
}

/* Loading */
.loading-container-neon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 3rem 0;
}

.spinner-neon {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(26, 115, 232, 0.15);
  border-top-color: var(--neon-cyan);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.spinner-tiny-neon {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-top-color: var(--neon-cyan);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
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

/* Photo Section */
.photo-section-neon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid rgba(26, 115, 232, 0.15);
}

.photo-selector-neon {
  cursor: pointer;
}

.photo-avatar-neon {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background-size: cover;
  background-position: center;
  border: 3px solid var(--neon-cyan);
  box-shadow: 0 0 25px rgba(26, 115, 232, 0.30);
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
}

.photo-avatar-neon:hover {
  transform: scale(1.05);
  box-shadow: 0 0 35px rgba(26, 115, 232, 0.45);
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

.file-input-hidden {
  display: none;
}

.btn-save-photo-neon {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  background: var(--neon-cyan);
  border: none;
  border-radius: 999px;
  color: var(--bg-primary);
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.25);
}

.btn-save-photo-neon:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 0 30px rgba(26, 115, 232, 0.40);
}

.btn-save-photo-neon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Fields */
.field-neon {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid rgba(26, 115, 232, 0.15);
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

.input-row-neon {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.input-neon {
  flex: 1;
  padding: 0.875rem 1rem;
  background: rgba(26, 115, 232, 0.06);
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-radius: 12px;
  color: var(--text-primary);
  font-size: 0.95rem;
  transition: all 0.3s ease;
  outline: none;
}

.input-neon:focus {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.25);
  background: #ffffff;
}

.input-neon.search {
  padding-left: 3rem;
}

.btn-save-inline-neon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  padding: 0;
  background: var(--neon-cyan);
  border: none;
  border-radius: 50%;
  color: var(--bg-primary);
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.btn-save-inline-neon:hover:not(:disabled) {
  transform: scale(1.1);
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.40);
}

.btn-save-inline-neon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.char-count-neon {
  font-size: 0.8rem;
  color: var(--text-dim);
  margin-top: 0.5rem;
  display: block;
}

/* Sections */
.members-section-neon,
.add-section-neon {
  margin-bottom: 2rem;
}

.section-header-neon {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  font-size: 1rem;
  font-weight: 700;
  color: var(--neon-cyan);
}

/* Members List */
.members-list-neon {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  max-height: 300px;
  overflow-y: auto;
  padding-right: 0.5rem;
}

.members-list-neon::-webkit-scrollbar {
  width: 6px;
}

.members-list-neon::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.06);
  border-radius: 10px;
}

.members-list-neon::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, var(--neon-cyan), var(--neon-purple));
  border-radius: 10px;
}

.member-item-neon {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 0.75rem;
  background: rgba(26, 115, 232, 0.06);
  border: 1px solid rgba(26, 115, 232, 0.15);
  border-radius: 12px;
  transition: all 0.2s ease;
}

.member-item-neon:hover {
  background: rgba(26, 115, 232, 0.10);
  border-color: rgba(26, 115, 232, 0.30);
}

.member-avatar-wrapper-neon {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid rgba(26, 115, 232, 0.40);
  box-shadow: 0 0 10px rgba(26, 115, 232, 0.25);
  flex-shrink: 0;
}

.member-avatar-img-neon {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.member-info-neon {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.member-name-neon {
  font-weight: 600;
  color: var(--text-primary);
}

.creator-badge-neon {
  font-size: 0.75rem;
  color: var(--neon-purple);
  font-weight: 600;
}

.btn-remove-member-neon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  background: rgba(255, 0, 100, 0.15);
  border: 1px solid rgba(255, 0, 100, 0.3);
  border-radius: 50%;
  color: #ff0064;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-remove-member-neon:hover:not(:disabled) {
  background: rgba(255, 0, 100, 0.25);
  box-shadow: 0 0 15px rgba(255, 0, 100, 0.3);
}

.btn-remove-member-neon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Search */
.search-container-neon {
  position: relative;
  margin-bottom: 1rem;
}

.search-icon-neon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(26, 115, 232, 0.45);
  pointer-events: none;
}

.search-status-neon {
  margin-bottom: 1rem;
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

/* Search Results */
.search-results-neon {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-height: 250px;
  overflow-y: auto;
  padding-right: 0.5rem;
}

.search-results-neon::-webkit-scrollbar {
  width: 6px;
}

.search-results-neon::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.06);
  border-radius: 10px;
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
  background: rgba(138, 43, 226, 0.08);
  border: 1px solid rgba(26, 115, 232, 0.10);
  border-radius: 12px;
  transition: all 0.2s ease;
}

.user-item-neon:hover {
  background: rgba(138, 43, 226, 0.12);
  border-color: rgba(138, 43, 226, 0.4);
}

.user-avatar-wrapper-neon {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid rgba(138, 43, 226, 0.5);
  box-shadow: 0 0 10px rgba(26, 115, 232, 0.20);
  flex-shrink: 0;
}

.user-avatar-img-neon {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-info-neon {
  flex: 1;
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
  background: rgba(26, 115, 232, 0.08);
  border: 1px solid rgba(138, 43, 226, 0.4);
  border-radius: 999px;
  color: var(--neon-purple);
  font-weight: 600;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add-neon:hover:not(:disabled) {
  background: rgba(138, 43, 226, 0.25);
  box-shadow: 0 0 15px rgba(26, 115, 232, 0.20);
}

.btn-add-neon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Modal */
.modal-overlay-neon {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-card-neon {
  background: #ffffff;
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-radius: 20px;
  padding: 2rem;
  max-width: 400px;
  width: 100%;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.20), 0 0 60px rgba(26, 115, 232, 0.15);
}

.modal-title-neon {
  margin: 0 0 1rem 0;
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--neon-cyan);
}

.modal-text-neon {
  margin: 0 0 1.5rem 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.modal-actions-neon {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
}

.btn-neon {
  padding: 0.75rem 1.5rem;
  border-radius: 999px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.btn-neon.danger {
  background: #ff0064;
  color: #ffffff;
}

.btn-neon.danger:hover {
  box-shadow: 0 0 20px rgba(255, 0, 100, 0.5);
  transform: translateY(-2px);
}

.btn-neon.secondary {
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
}

.btn-neon.secondary:hover {
  background: var(--bg-hover);
  border-color: var(--neon-cyan);
}
</style>
