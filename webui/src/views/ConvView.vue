<template>
  <section class="conversation-neon">
        <!-- Header -->
    <header class="conv-header-neon">
      <router-link v-if="!embedded" to="/home" class="back-btn-neon">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"></line>
          <polyline points="12 19 5 12 12 5"></polyline>
        </svg>
        Back
      </router-link>
      
      <div class="peer-info-neon">
        <img 
          v-if="conversationPhoto" 
          :src="conversationPhoto" 
          class="avatar-neon" 
          :alt="(conversation.name || 'Chat') + ' photo'"
          @click="openUserProfileModal"
        />
        <div class="meta-neon">
          <h2 class="name-neon" @click="openUserProfileModal">
            {{ conversationName }}
          </h2>
          <div class="subtitle-neon" v-if="conversation.membersIds?.length">
            {{ conversation.membersIds.length }} member(s)
          </div>
          <div v-if="isGroupChat" class="quick-actions-neon">
            <button type="button" class="link-neon" @click="openParticipants">Participants</button>
            <button type="button" class="link-neon" :disabled="!isMember" @click="editGroup">Edit Group</button>
          </div>
        </div>
      </div>
      
      <div class="header-actions-neon">
        <template v-if="isGroupChat">
          <button v-if="!hasLeftGroup" type="button" class="btn-neon danger" :disabled="!isMember" @click="leaveGroup">
            Leave Group
          </button>
          <button v-else type="button" class="btn-neon danger" @click="deleteGroup">
            Delete Group
          </button>
        </template>
      </div>
    </header>

    <!-- Messages Area -->
    <div ref="scrollArea" class="messages-neon">
      <div v-if="loading" class="loading-spinner-neon">
        <div class="spinner-neon"></div>
        <p>Loading messages...</p>
      </div>
      
      <div v-if="errorMessage && !loading" class="error-neon">
        {{ errorMessage }}
      </div>
      
      <div v-if="toast.show" class="toast-neon">
        {{ toast.msg }}
        <button v-if="toast.targetId" class="link-neon" @click="openToastTarget">Open chat</button>
      </div>

      <template v-if="!loading">
        <template v-for="(item, idx) in groupedMessages" :key="getItemKey(item, idx)">
        <!-- Date Separator -->
        <div v-if="item.type === 'date'" class="date-separator-neon">
          <div class="date-badge-neon">{{ formatDate(item.timestamp) }}</div>
        </div>
        
        <!-- System Messages -->
        <div v-else-if="item.type === 'creationNotice'" class="date-separator-neon">
          <div class="date-badge-neon">{{ creatorName }} created this group</div>
        </div>
        
        <div v-else-if="item.type === 'leaveNotice'" class="date-separator-neon">
          <div class="date-badge-neon">{{ item.username }} left the group chat</div>
        </div>
        
        <div v-else-if="item.type === 'memberAdded'" class="date-separator-neon">
          <div class="date-badge-neon">{{ item.message }}</div>
        </div>

        <!-- Regular Message -->
        <div v-else-if="item.type === 'message'" class="msg-row-neon" :class="{ own: isOwn(item.value) }">
          <div class="bubble-neon">
            <div class="sender-name-neon">
              {{ item.value.sender?.username || 'Unknown' }}
            </div>
            
            <div v-if="item.value.forwarded_from" class="fwd-neon">
              <svg class="fwd-icon" viewBox="0 0 16 16" fill="none">
                <path d="M1 8h10M8 4l4 4-4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              Forwarded
            </div>

            <div v-if="getReplyInfo(item.value)" class="reply-container-neon">
              <div class="reply-username-neon">{{ getReplyInfo(item.value).username }}</div>
              <div class="reply-message-neon">{{ getReplyInfo(item.value).preview }}</div>
            </div>

            <div v-if="isImageMessage(item.value)" class="attachment-neon">
              <img :src="imageSrc(item.value)" alt="Image" />
            </div>
            <div v-else-if="getActualMessageText(item.value)" class="text-neon">
              {{ getActualMessageText(item.value) }}
            </div>

            <div class="meta-line-neon">
              <span class="time-neon">{{ formatTime(item.value.timestamp) }}</span>
              <span 
                  v-if="isOwn(item.value) && (item.value.message_status || item.value.messageStatus)" 
                  class="status-neon" 
                  :class="statusClass(item.value.message_status || item.value.messageStatus)"
                  >
                {{ statusIcon(item.value.message_status || item.value.messageStatus) }}
              </span>
            </div>

            <div class="actions-neon">
              <button type="button" class="link-neon" @click.stop.prevent="openReply(item.value)">Reply</button>
              <button type="button" class="link-neon" @click.stop.prevent="openForward(item.value.id)">Forward</button>
              <button v-if="myReaction(item.value)" type="button" class="link-neon" @click.stop.prevent="removeMyReaction(item.value)">Remove Reaction</button>
              <button v-if="isOwn(item.value)" type="button" class="link-neon danger" @click.stop.prevent="del(item.value.id)">Delete</button>
            </div>

            <!-- Reaction Trigger -->
            <div
              class="reaction-trigger-neon"
              :class="{ left: isOwn(item.value), right: !isOwn(item.value), open: reactionPicker.openFor === item.value.id }"
              @click.stop.prevent="togglePicker(item.value.id)"
            >
              +
            </div>

            <!-- Reaction Picker -->
            <div 
              v-if="reactionPicker.openFor === item.value.id" 
              class="reaction-picker-neon" 
              :class="{ left: isOwn(item.value), right: !isOwn(item.value) }"
            >
              <button class="rx-neon" :class="{ active: myReaction(item.value) === 'like' }" @click.stop.prevent="react(item.value.id, 'like')">👍</button>
              <button class="rx-neon" :class="{ active: myReaction(item.value) === 'heart' }" @click.stop.prevent="react(item.value.id, 'heart')">❤️</button>
              <button class="rx-neon" :class="{ active: myReaction(item.value) === 'laugh' }" @click.stop.prevent="react(item.value.id, 'laugh')">😂</button>
              <button class="rx-neon" :class="{ active: myReaction(item.value) === 'sad_face' }" @click.stop.prevent="react(item.value.id, 'sad_face')">😢</button>
              <button class="rx-neon" :class="{ active: myReaction(item.value) === 'angry_face' }" @click.stop.prevent="react(item.value.id, 'angry_face')">😡</button>
            </div>

            <!-- Reactions Display -->
            <div v-if="getTotalReactionCount(item.value) > 0" class="reactions-container-neon">
              <div
                v-for="(group, gIdx) in getGroupedReactions(item.value)"
                :key="'g-' + gIdx"
                class="reaction-group-wrapper-neon"
              >
                <span class="reaction-group-neon" :class="{ mine: group.mine }">
                  {{ reactionEmoji(group.type) }}
                </span>
                <div class="reaction-tooltip-neon">
                  <div class="tooltip-label">Reacted:</div>
                  <div class="tooltip-usernames">{{ group.usernames.join(', ') }}</div>
                </div>
              </div>
              <span class="reaction-count-neon">{{ getTotalReactionCount(item.value) }}</span>
            </div>
          </div>
        </div>

        <!-- Combined Message (Reply + Image) -->
        <div v-else-if="item.type === 'messageCombined'" class="msg-row-neon" :class="{ own: isOwn(item.primary) }">
          <div class="bubble-neon">
            <div class="sender-name-neon">
              {{ item.primary?.sender?.username || 'Unknown' }}
            </div>

            <div v-if="item.primary?.forwarded_from" class="fwd-neon">
              <svg class="fwd-icon" viewBox="0 0 16 16" fill="none">
                <path d="M1 8h10M8 4l4 4-4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              Forwarded
            </div>

            <div class="reply-container-neon">
              <div class="reply-username-neon">{{ item.reply.username }}</div>
              <div class="reply-message-neon">{{ item.reply.preview }}</div>
            </div>

            <div v-if="isImageMessage(item.primary)" class="attachment-neon">
              <img :src="imageSrc(item.primary)" alt="Image" />
            </div>

            <div class="meta-line-neon">
              <span class="time-neon">{{ formatTime(item.primary.timestamp) }}</span>
              <span 
                v-if="isOwn(item.primary) && item.primary.message_status" 
                class="status-neon" 
                :class="statusClass(item.primary.message_status)"
              >
                {{ statusIcon(item.primary.message_status) }}
              </span>
            </div>

            <div class="actions-neon">
              <button type="button" class="link-neon" @click.stop.prevent="openReply(item.primary)">Reply</button>
              <button type="button" class="link-neon" @click.stop.prevent="openForward(item.primary.id)">Forward</button>
              <button v-if="myReaction(item.primary)" type="button" class="link-neon" @click.stop.prevent="removeMyReaction(item.primary)">Remove Reaction</button>
              <button v-if="isOwn(item.primary)" type="button" class="link-neon danger" @click.stop.prevent="del(item.primary.id)">Delete</button>
            </div>

            <div
              class="reaction-trigger-neon"
              :class="{ left: isOwn(item.primary), right: !isOwn(item.primary), open: reactionPicker.openFor === item.primary.id }"
              @click.stop.prevent="togglePicker(item.primary.id)"
            >
              +
            </div>

            <div 
              v-if="reactionPicker.openFor === item.primary.id" 
              class="reaction-picker-neon" 
              :class="{ left: isOwn(item.primary), right: !isOwn(item.primary) }"
            >
              <button class="rx-neon" :class="{ active: myReaction(item.primary) === 'like' }" @click.stop.prevent="react(item.primary.id, 'like')">👍</button>
              <button class="rx-neon" :class="{ active: myReaction(item.primary) === 'heart' }" @click.stop.prevent="react(item.primary.id, 'heart')">❤️</button>
              <button class="rx-neon" :class="{ active: myReaction(item.primary) === 'laugh' }" @click.stop.prevent="react(item.primary.id, 'laugh')">😂</button>
              <button class="rx-neon" :class="{ active: myReaction(item.primary) === 'sad_face' }" @click.stop.prevent="react(item.primary.id, 'sad_face')">😢</button>
              <button class="rx-neon" :class="{ active: myReaction(item.primary) === 'angry_face' }" @click.stop.prevent="react(item.primary.id, 'angry_face')">😡</button>
            </div>

            <div v-if="getTotalReactionCount(item.primary) > 0" class="reactions-container-neon">
              <div
                v-for="(group, gIdx) in getGroupedReactions(item.primary)"
                :key="'g-' + gIdx"
                class="reaction-group-wrapper-neon"
              >
                <span class="reaction-group-neon" :class="{ mine: group.mine }">
                  {{ reactionEmoji(group.type) }}
                </span>
                <div class="reaction-tooltip-neon">
                  <div class="tooltip-label">Reacted:</div>
                  <div class="tooltip-usernames">{{ group.usernames.join(', ') }}</div>
                </div>
              </div>
              <span class="reaction-count-neon">{{ getTotalReactionCount(item.primary) }}</span>
            </div>
          </div>
        </div>
      </template>
      </template>
    </div>

    <!-- Composer -->
    <footer class="composer-neon">
      <div v-if="reply.active" class="reply-banner-neon">
        <div class="reply-content-neon">
          <div class="reply-label-neon">Replying to:</div>
          <div class="reply-username-neon">{{ reply.username || 'User' }}</div>
          <div class="reply-message-neon">{{ reply.preview }}</div>
        </div>
        <button class="reply-close-neon" @click="cancelReply">×</button>
      </div>
      <div v-if="photoAttachB64 && !reply.active" class="photo-preview-banner-neon">
        <div class="photo-preview-content-neon">
          <div class="photo-preview-label-neon">Photo attached</div>
          <div class="photo-preview-thumbnail-neon">
            <img :src="'data:image/png;base64,' + photoAttachB64" alt="Photo preview" />
          </div>
        </div>
        <button class="photo-preview-close-neon" @click="removePhoto">×</button>
      </div>
      <div class="input-row-neon">
        <input
          v-model="newMessage"
          class="input-neon"
          type="text"
          placeholder="Type your message..."
          @keyup.enter="send"
          ref="messageInput"
          :disabled="!isMember && isGroupChat"
        />
        
        <input
          ref="fileInput"
          class="file-input-hidden"
          type="file"
          accept="image/*"
          @click="resetFileInput"
          @change="attachPhoto"
          :disabled="!isMember && isGroupChat"
        />
        
        <button
          type="button"
          class="attach-btn-neon"
          @click="openFilePicker"
          :disabled="!isMember && isGroupChat"
        >
          <svg viewBox="0 0 24 24" class="attach-icon-neon">
            <rect x="4" y="5" width="16" height="14" rx="3" ry="3" />
            <path d="M8.5 13.5 11 16l3.5-4.5 4 5" />
            <circle cx="10" cy="9" r="1.5" />
          </svg>
        </button>
        
        <button 
          class="btn-neon send" 
          :disabled="!canSend || sending" 
          @click="send"
        >
          {{ sending ? 'Sending...' : 'Send' }}
        </button>
      </div>
    </footer>

    <!-- Forward Modal -->
    <div v-if="forward.open" class="modal-overlay-neon" @click.self="closeForward">
      <div class="modal-card-neon">
        <h3 class="modal-title-neon">Forward message to</h3>
        
        <div class="forward-search-neon">
          <input 
            v-model="forward.newUsername" 
            class="input-neon" 
            placeholder="Search for users..." 
            @input="onForwardUserInput"
            ref="forwardSearchInput"
          />
          
          <div v-if="forward.suggestLoading" class="hint-neon">Searching...</div>
          <div v-else-if="forward.newUsername && !forward.suggestions.length" class="hint-neon">No users found</div>
          
          <div v-if="forward.suggestions && forward.suggestions.length" class="forward-suggestions-neon">
            <div 
              class="forward-item-neon" 
              v-for="u in forward.suggestions" 
              :key="u.id" 
              @click="selectForwardUser(u)"
            >
              {{ u.username }}
            </div>
          </div>
        </div>
        
        <div v-if="allConversations && allConversations.length" class="recent-chats-neon">
          <h4>Recent chats</h4>
          <div 
            class="forward-item-neon" 
            v-for="c in allConversations" 
            :key="c.conversationId" 
            @click="forwardToConversation(c.conversationId)"
          >
            {{ c.displayName }}
          </div>
        </div>

        <div class="modal-actions-neon">
          <button class="btn-neon secondary" @click="closeForward">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Delete Message Modal -->
    <div v-if="deleteDialog.open" class="modal-overlay-neon" @click.self="cancelDelete">
      <div class="modal-card-neon small">
        <h3 class="modal-title-neon">Delete Message</h3>
        <p class="modal-text-neon">Are you sure you want to delete this message? This action cannot be undone.</p>
        <div class="modal-actions-neon">
          <button class="btn-neon danger" @click="confirmDelete">Delete</button>
          <button class="btn-neon secondary" @click="cancelDelete">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Participants Modal -->
    <div v-if="participantsModal.open" class="modal-overlay-neon" @click.self="closeParticipants">
      <div class="modal-card-neon">
        <h3 class="modal-title-neon">Group participants</h3>
        
        <div v-if="participantsList.length === 0" class="hint-neon">No participants found</div>
        <div v-else class="participants-list-neon">
          <div class="participant-item-neon" v-for="participant in participantsList" :key="participant.id">
            <div class="participant-avatar-neon">
              <img :src="getParticipantPhoto(participant)" :alt="participant.username" />
            </div>
            <div class="participant-name-neon">{{ participant.username }}</div>
          </div>
        </div>
        
        <div class="modal-actions-neon">
          <button class="btn-neon secondary" @click="closeParticipants">Close</button>
        </div>
      </div>
    </div>

    <!-- Leave Group Modal -->
    <div v-if="leaveGroupDialog.open" class="modal-overlay-neon" @click.self="cancelLeaveGroup">
      <div class="modal-card-neon small">
        <h3 class="modal-title-neon">Leave Group</h3>
        <p class="modal-text-neon">Are you sure you want to leave this group? You won't be able to send or receive messages from this group anymore.</p>
        <div class="modal-actions-neon">
          <button class="btn-neon danger" @click="confirmLeaveGroup">Leave Group</button>
          <button class="btn-neon secondary" @click="cancelLeaveGroup">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Delete Group Modal -->
    <div v-if="deleteGroupDialog.open" class="modal-overlay-neon" @click.self="cancelDeleteGroup">
      <div class="modal-card-neon small">
        <h3 class="modal-title-neon">Delete Group</h3>
        <p class="modal-text-neon">Are you sure you want to delete this group from your chat list? This action cannot be undone.</p>
        <div class="modal-actions-neon">
          <button class="btn-neon danger" @click="confirmDeleteGroup">Delete Group</button>
          <button class="btn-neon secondary" @click="cancelDeleteGroup">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Delete Conversation Modal -->
    <div v-if="deleteConversationDialog.open" class="modal-overlay-neon" @click.self="cancelDeleteConversation">
      <div class="modal-card-neon small">
        <h3 class="modal-title-neon">Delete Chat</h3>
        <p class="modal-text-neon">Are you sure you want to delete this chat from your list? This action cannot be undone.</p>
        <div class="modal-actions-neon">
          <button class="btn-neon danger" @click="confirmDeleteConversation">Delete Chat</button>
          <button class="btn-neon secondary" @click="cancelDeleteConversation">Cancel</button>
        </div>
      </div>
    </div>

    <!-- User Profile Modal -->
    <div v-if="userProfileModal.open" class="modal-overlay-neon" @click.self="closeUserProfileModal">
      <div class="profile-modal-neon">
        <button class="close-btn-neon" @click="closeUserProfileModal">×</button>
        <div class="profile-content-neon">
          <img 
            v-if="userProfileModal.photo" 
            :src="userProfileModal.photo" 
            alt="Profile" 
            class="profile-photo-neon"
          />
          <div v-else class="profile-photo-default-neon">
            {{ userProfileModal.username?.charAt(0).toUpperCase() || '?' }}
          </div>
          <h3 class="profile-username-neon">@{{ userProfileModal.username }}</h3>
          <p class="profile-description-neon">{{ userProfileModal.description || 'Hey, I\'m using WASA!' }}</p>
        </div>
      </div>
    </div>

    <!-- Group Edit Modal -->
    <div v-if="groupEditOpen" class="modal-overlay-neon" @click.self="closeGroupEdit">
      <div class="modal-card-neon large">
        <GroupEditView 
          :group-id-prop="effectiveConversationId" 
          :embedded="true" 
          @close="closeGroupEdit" 
        />
      </div>
    </div>
  </section>
</template>

<script>
import GroupEditView from './GroupEditView.vue'

export default {
  name: 'ConvView',
  components: { GroupEditView },
  
  props: {
    conversationId: {
      type: String,
      default: null
    },
    embedded: {
      type: Boolean,
      default: false
    }
  },
  
  data() {
    return {
      loading: true,
      sending: false,
      forwarding: false,
      errorMessage: null,
      newMessage: '',
      photoAttachB64: '',
      userScrolled: false, 
      isLoadingInProgress: false, 
    lastMessageCount: 0,
      toast: { show: false, msg: '', targetId: '' },
      reply: { active: false, preview: '', username: '' },
      deleteDialog: { open: false, messageId: null, pairedId: null },
      reactionPicker: { openFor: '' },
      tempMyReaction: {},
      conversation: {
        id: null,
        displayName: '',
        photo: null,
        membersIds: [],
        messages: []
      },
      participantsModal: { open: false },
      leaveGroupDialog: { open: false },
      deleteGroupDialog: { open: false },
      deleteConversationDialog: { open: false },
      hasLeftGroup: false,
      allConversations: [],
      forward: { 
        open: false, 
        messageId: null, 
        target: '', 
        newUsername: '', 
        suggestions: [], 
        suggestTimer: null, 
        suggestLoading: false, 
        selectedUserId: '' 
      },
      pollId: null,   
      userId: localStorage.getItem('userId') || null,
      reactBusy: {},
      groupEditOpen: false,
      userProfileModal: { open: false, username: '', photo: null, description: '' }
    }
  },

  computed: {
    isMember() {
      if (!this.isGroupChat) return true
      const me = String(this.userId || '')
      const ids = (this.conversation?.membersIds || []).map(x => String(x || ''))
      const partIds = (this.conversation?.participants || []).map(p => String((p && (p.id || p.userId)) || ''))
      const arraysPopulated = (ids.length > 0) || (partIds.length > 0)
      const inGroup = arraysPopulated ? (ids.includes(me) || partIds.includes(me)) : true
      return inGroup && !this.hasLeftGroup
    },
    
    effectiveConversationId() {
      return this.conversationId || this.$route.params.conversationId
    },
    
    conversationPhoto() {
  // Per gruppi
  if (this.isGroupChat) {
    const b64 = this.conversation?.groupPhoto || this.conversation?.group_photo || this.conversation?.GroupPhoto
    if (b64) return 'data:image/png;base64,' + b64
    const name = this.conversation?.name || this.conversation?.displayName || 'Group'
    return this.makeLetterAvatar(name)
  }
  
  // Per chat dirette: prendi la foto dell'altro utente
  try {
    const myUserId = String(this.userId || '')
    const participants = this.conversation?.participants || []
    
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
      return this.makeLetterAvatar(peerName)
    }
    
    const name = this.conversation?.name || this.conversation?.displayName || 'Chat'
    return this.makeLetterAvatar(name)
  } catch (error) {
    console.error('Error getting conversation photo:', error)
    return this.makeLetterAvatar('?')
  }
},
    conversationName() {
    // Per i gruppi
    if (this.isGroupChat) {
      return this.conversation?.name || this.conversation?.displayName || 'Group Chat'
    }
    
    // Per chat dirette: trova l'altro utente
    try {
      const myUserId = String(this.userId || '')
      const participants = this.conversation?.participants || []
      
      // Cerca l'altro utente (non te stesso)
      const peer = participants.find(p => {
        const pId = String(p?.id || p?.userId || '')
        return pId && pId !== myUserId
      })
      
      if (peer?.username) {
        return peer.username
      }
      
      // Fallback: prova con membersIds
      const members = this.conversation?.members || []
      const peerMember = members.find(m => {
        const mId = String(m?.id || m?.userId || '')
        return mId && mId !== myUserId
      })
      
      if (peerMember?.username) {
        return peerMember.username
      }
      
      // Ultimo fallback
      return this.conversation?.name || this.conversation?.displayName || 'Chat'
    } catch (error) {
      console.error('Error getting conversation name:', error)
      return 'Chat'
    }
  },
  
    
    canSend() {
      if (this.isGroupChat && !this.isMember) return false
      return (this.newMessage && this.newMessage.trim().length > 0) || !!this.photoAttachB64
    },
    
    groupedMessages() {
      const messages = this.conversation?.messages || []
      const groups = []
      let lastDate = null
      const creationTs = this.getCreationTimestamp()
      const creationDay = creationTs ? this.getMessageDate(creationTs) : null
      let creationInserted = false

      if (messages.length === 0 && this.isGroupChat && creationDay) {
        groups.push({ type: 'date', value: creationDay, timestamp: creationTs })
        groups.push({ type: 'creationNotice' })
        return groups
      }

      for (let i = 0; i < messages.length; i++) {
        const msg = messages[i]
        const msgDate = this.getMessageDate(msg.timestamp)
        
        if (msgDate !== lastDate) {
          groups.push({ type: 'date', value: msgDate, timestamp: msg.timestamp })
          if (this.isGroupChat && creationDay && msgDate === creationDay && !creationInserted) {
            groups.push({ type: 'creationNotice' })
            creationInserted = true
          }
          lastDate = msgDate
        }

        if (msg.type === 'system' && msg.systemType === 'leave') {
          groups.push({ type: 'leaveNotice', username: msg.username })
          continue
        }
        
        const actualText = this.getActualMessageText(msg) || ''
        if (this.isGroupChat && typeof actualText === 'string' && actualText.endsWith('left the group chat')) {
          const uname = actualText.replace(' left the group chat', '').trim()
          if (uname) {
            groups.push({ type: 'leaveNotice', username: uname })
            continue
          }
        }

        const contentType = msg?.content?.type || ''
        const contentValue = msg?.content?.value || ''
        if (contentType === 'system' && contentValue.startsWith('system:member_added:')) {
          const parts = contentValue.split(':')
          if (parts.length >= 6) {
            const adderID = parts[2]
            const adderUsername = parts[3]
            const addedUserID = parts[4]
            const addedUsername = parts[5]
            const message = this.formatMemberAddedMessage(adderID, adderUsername, addedUserID, addedUsername)
            groups.push({ type: 'memberAdded', message })
            continue
          }
        }

        const replyInfo = this.getReplyInfo(msg)
        const isReplyOnly = !!replyInfo && (!replyInfo.actualMessage || replyInfo.actualMessage.length === 0) && (String(msg?.content?.type || 'text').toLowerCase() === 'text')
        const next = messages[i + 1]
        
        if (isReplyOnly && next && (next?.sender?.id || next?.senderId) === (msg?.sender?.id || msg?.senderId) && this.isImageMessage(next) && this.isNearSameDay(msg.timestamp, next.timestamp)) {
          groups.push({ type: 'messageCombined', reply: replyInfo, primary: next })
          i += 1
          continue
        }

        groups.push({ type: 'message', value: msg })
      }
      
      return groups
    },
    
    creatorName() {
      const cid = this.creatorId()
      if (!cid) return 'User'
      const name = this.getUsernameById(cid)
      return name || 'User'
    },
    
    isGroupChat() {
      return this.conversation?.isGroup === true || (this.conversation?.membersIds?.length || 0) > 2
    },
    
    participantsList() {
      const participants = this.conversation?.participants || []
      
      if (this.hasLeftGroup) {
        return participants.filter(p => {
          const pId = String(p?.id || p?.ID || '')
          const uId = String(this.userId || '')
          return pId !== uId
        })
      }
      
      return participants
    }
  },

  methods: {
    
    getItemKey(item, idx) {
      if (item.type === 'date') return 'date-' + idx
      if (item.type === 'creationNotice') return 'crt-' + idx
      if (item.type === 'leaveNotice') return 'lv-' + idx
      if (item.type === 'memberAdded') return 'ma-' + idx
      if (item.type === 'messageCombined') return 'cmb-' + (item.primary?.id || idx)
      return item.value?.id || idx
    },
     creatorId() {
          try {
            const cid = this.conversation?.createdBy || this.conversation?.creatorId || null;
            return cid;
          } catch { return null; }
        },
        formatMemberAddedMessage(adderID, adderUsername, addedUserID, addedUsername) {
          const currentUserId = String(this.userId || '');
          
          if (String(adderID) === currentUserId) {
            return `You added ${addedUsername} to the group chat`;
          }
          
          if (String(addedUserID) === currentUserId) {
            return `You have been added in this group by ${adderUsername}`;
          }
          
          return `${adderUsername} added ${addedUsername} to the group chat`;
        },
        openParticipants() {
          this.participantsModal.open = true;
        },
        editGroup() {
          if (!this.isGroupChat) return;
          if (!this.isMember) { this.showToast('You are no longer in this group'); return; }
          const id = this.effectiveConversationId;
          if (!id) return;
          this.groupEditOpen = true;
        },
        closeParticipants() {
          this.participantsModal.open = false;
        },
        getParticipantPhoto(participant) {
          if (!participant) return this.makeLetterAvatar('');
          const photo = participant.photo || participant.Photo;
          if (typeof photo === 'string' && photo.length > 0) {
            if (photo.startsWith('data:')) return photo;
            if (photo.trim()) return `data:image/png;base64,${photo}`;
          }
          const name = participant.username || participant.name || '';
          return this.makeLetterAvatar(name);
        },
        async leaveGroup() {
          this.leaveGroupDialog.open = true;
        },
        cancelLeaveGroup() {
          this.leaveGroupDialog.open = false;
        },
        deleteGroup() {
          this.deleteGroupDialog.open = true;
        },
        deleteConversation() {
          this.deleteConversationDialog.open = true;
        },
        cancelDeleteGroup() {
          this.deleteGroupDialog.open = false;
        },
        cancelDeleteConversation() {
          this.deleteConversationDialog.open = false;
        },
        async confirmDeleteGroup() {
          this.deleteGroupDialog.open = false;
          
          try {

            const raw = localStorage.getItem('leftConversations') || '{}';
            let leftMap = {};
            try { leftMap = JSON.parse(raw) || {}; } catch { leftMap = {}; }
            
            const deletedGroupId = this.effectiveConversationId;
            
            delete leftMap[deletedGroupId];
            localStorage.setItem('leftConversations', JSON.stringify(leftMap));
            
            this.showToast('Group deleted from your chats');
            
            await this.loadConversationsList();
            
            const deletedIdStr = String(deletedGroupId || '');
            const availableConversations = (this.allConversations || []).filter(conv => {
              const cid = String(conv?.id || conv?.conversationId || conv?.conversation_id || conv?.ID || '');
              return cid && cid !== deletedIdStr;
            });

            const pickNextConversationId = (list) => {
              if (!Array.isArray(list) || list.length === 0) return null;
              const tsOf = (c) => {
                try {
                  const lm = c?.lastMessage || c?.LastMessage || {};
                  const t1 = lm?.timestamp ? new Date(lm.timestamp).getTime() : NaN;
                  const t2 = c?.createdAt ? new Date(c.createdAt).getTime() : NaN;
                  const t3 = c?.CreatedAt ? new Date(c.CreatedAt).getTime() : NaN;
                  const t4 = c?._snapshotLastTs || NaN;
                  return [t1, t2, t3, t4].find(v => Number.isFinite(v)) || 0;
                } catch { return 0; }
              };
              const sorted = [...list].sort((a, b) => tsOf(b) - tsOf(a));
              const top = sorted[0] || {};
              return top.id || top.conversationId || top.conversation_id || top.ID || null;
            };

            const nextConversationId = pickNextConversationId(availableConversations);

            if (this.embedded) {
              this.$emit('conversation-deleted', { deletedId: deletedGroupId, nextId: nextConversationId });
              if (nextConversationId) {
                this.$router.push({ path: '/home', query: { conv: nextConversationId } });
              } else {
                this.$router.push('/home');
              }
              return;
            }

            if (nextConversationId) {
              this.$router.push({ path: '/home', query: { conv: nextConversationId } });
            } else {
              this.$router.push('/home');
            }
          } catch (e) {
            console.error('Failed to delete group:', e);
            this.errorMessage = 'Failed to delete group';
          }
        },
        async confirmDeleteConversation() {
          this.deleteConversationDialog.open = false;
          const deletedId = this.effectiveConversationId;

          try {
            try {
              const token = localStorage.getItem('token');
              const cfg = token ? { headers: { Authorization: `Bearer ${token}` } } : {};
              await this.$axios.delete(`/conversations/${deletedId}`, cfg);
            } catch (apiErr) {
              console.warn('Delete chat API not available or failed, proceeding locally', apiErr);
            }

            await this.loadConversationsList();

            const deletedIdStr = String(deletedId || '');
            const availableConversations = (this.allConversations || []).filter(conv => {
              const cid = String(conv?.id || conv?.conversationId || conv?.conversation_id || conv?.ID || '');
              return cid && cid !== deletedIdStr;
            });

            const pickNextConversationId = (list) => {
              if (!Array.isArray(list) || list.length === 0) return null;
              const tsOf = (c) => {
                try {
                  const lm = c?.lastMessage || c?.LastMessage || {};
                  const t1 = lm?.timestamp ? new Date(lm.timestamp).getTime() : NaN;
                  const t2 = c?.createdAt ? new Date(c.createdAt).getTime() : NaN;
                  const t3 = c?.CreatedAt ? new Date(c.CreatedAt).getTime() : NaN;
                  const t4 = c?._snapshotLastTs || NaN;
                  return [t1, t2, t3, t4].find(v => Number.isFinite(v)) || 0;
                } catch { return 0; }
              };
              const sorted = [...list].sort((a, b) => tsOf(b) - tsOf(a));
              const top = sorted[0] || {};
              return top.id || top.conversationId || top.conversation_id || top.ID || null;
            };

            const nextConversationId = pickNextConversationId(availableConversations);

            if (this.embedded) {
              this.$emit('conversation-deleted', { deletedId, nextId: nextConversationId });
              if (nextConversationId) {
                this.$router.push({ path: '/home', query: { conv: nextConversationId } });
              } else {
                this.$router.push('/home');
              }
              return;
            }

            if (nextConversationId) {
              this.$router.push({ path: '/home', query: { conv: nextConversationId } });
            } else {
              this.$router.push('/home');
            }
          } catch (e) {
            console.error('Failed to delete conversation:', e);
            this.errorMessage = 'Failed to delete conversation';
          }
        },
        removePhoto() {
      this.photoAttachB64 = '';
      if (this.$refs.fileInput) this.$refs.fileInput.value = '';
    },
       async confirmLeaveGroup() {
  this.leaveGroupDialog.open = false;
  
  try {
    const token = localStorage.getItem('token');
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const axiosCfg = token ? { headers: { Authorization: `Bearer ${token}` } } : {};
    const myUsername = this.getUsernameById(this.userId) || 'A user';
    const serverSystemMessage = `${myUsername} left the group chat`;
    
    try {
      
      await this.$axios.post(
        `/conversations/${this.effectiveConversationId}/messages`,
        { content: { type: 'text', value: this.toBase64(serverSystemMessage) } },
        axiosCfg
      );
      
      const localSystemMessage = 'You left the group chat';
      const nowIso = new Date().toISOString();
      this.conversation.messages = [
        ...(this.conversation?.messages || []),
        {
          id: `temp-leave-${Date.now()}`,
          timestamp: nowIso,
          sender: { id: this.userId, username: myUsername },
          content: { type: 'text', value: this.toBase64(localSystemMessage) }
        }
      ];
      
      try {
        const raw = localStorage.getItem('leftConversations') || '{}';
        let leftMap = {};
        try { leftMap = JSON.parse(raw) || {}; } catch { leftMap = {}; }
        leftMap[this.effectiveConversationId] = {
          userId: this.userId,
          leftAt: nowIso,
          conversation: {
            ...(this.conversation || {}),
            messages: (this.conversation?.messages || []).filter(m => {
              try { return new Date(m.timestamp).getTime() <= new Date(nowIso).getTime(); } catch { return true; }
            })
          }
        };
        localStorage.setItem('leftConversations', JSON.stringify(leftMap));
      } catch {}
      
      this.$nextTick(() => this.scrollToBottom());
    } catch (msgErr) {
      console.error('Could not send leave message:', msgErr);
    }
 
    await this.$axios.delete(`/groups/${this.effectiveConversationId}?userId=${this.userId}`, {
      headers
    });
    
    this.conversation.membersIds = (this.conversation?.membersIds || []).filter(id => String(id) !== String(this.userId));
    this.conversation.participants = (this.conversation?.participants || []).filter(p => {
      const pId = String(p?.id || p?.userId || '')
      return pId !== String(this.userId)
    });
    
    this.hasLeftGroup = true;
    this.showToast('You left the group');
    
  } catch (e) {
    console.error('Failed to leave group:', e);
    console.error('Error response:', e.response?.data);
    console.error('Error status:', e.response?.status);
    this.errorMessage = 'Failed to leave group: ' + (e.response?.data || e.message);
    this.showToast('Failed to leave group');
  }
},
    onKeyDown(e) {
      if (e?.key === 'Escape' && this.forward.open) this.closeForward();
    },
    myReaction(message) {
      try {
        const mid = message?.id;
        const local = mid ? (this.tempMyReaction[mid] || '') : '';
        if (local) return String(local).toLowerCase();
        const uid = this.userId;
        const arr = message?.Reactions || message?.reactions || [];
        const mine = arr.find(r => (r.authorId || r.AuthorID || r.userId || r.user_id) === uid);
        return (mine?.type || mine?.Type || mine?.emoji || '').toLowerCase();
      } catch { return ''; }
    },
    getAllReactions(message) {
      try {
        const arr = message?.Reactions || message?.reactions || [];
        return arr.map(r => ({
          type: (r?.type || r?.Type || '').toLowerCase(),
          userId: r?.authorId || r?.AuthorID || r?.userId || '',
          username: this.getUsernameById(r?.authorId || r?.AuthorID || r?.userId || '')
        })).filter(r => r.type);
      } catch {
        return [];
      }
    },
    getUsernameById(userId) {
      if (!userId) return 'Unknown';
      const members = this.conversation?.members || [];
      const user = members.find(m => m.userId === userId || m.id === userId);
      if (user?.username) return user.username;
      const participants = this.conversation?.participants || [];
      const participant = participants.find(p => p.id === userId || p.userId === userId);
      if (participant?.username) return participant.username;
      return userId;
    },
    groupReactions(message) {
      try {
        const arr = message?.Reactions || message?.reactions || [];
        const map = new Map();
        for (const r of arr) {
          const type = (r?.type || r?.Type || '').toLowerCase();
          const uid = r?.authorId || r?.AuthorID || r?.userId || '';
          if (!type) continue;
          if (!map.has(type)) map.set(type, { type, users: [], count: 0, mine: false });
          const g = map.get(type);
          g.users.push(uid);
          g.count += 1;
          if (uid === this.userId) g.mine = true;
        }
        return Array.from(map.values());
      } catch {
        return [];
      }
    },
    getGroupedReactions(message) {
      try {
        const arr = message?.Reactions || message?.reactions || [];
        const map = new Map();
        for (const r of arr) {
          const type = (r?.type || r?.Type || '').toLowerCase();
          const uid = r?.authorId || r?.AuthorID || r?.userId || '';
          if (!type) continue;
          if (!map.has(type)) {
            map.set(type, { type, userIds: [], usernames: [], mine: false });
          }
          const g = map.get(type);
          g.userIds.push(uid);
          g.usernames.push(this.getUsernameById(uid));
          if (String(uid) === String(this.userId)) g.mine = true;
        }
        return Array.from(map.values());
      } catch {
        return [];
      }
    },
    getTotalReactionCount(message) {
      try {
        const arr = message?.Reactions || message?.reactions || [];
        return arr.length;
      } catch {
        return 0;
      }
    },
     async load() {
      if (this.isLoadingInProgress) {
      console.log('⏳ Load already in progress, skipping...')
      return
    }
    this.isLoadingInProgress = true
    this.errorMessage = null;
    const shouldScrollToBottom = !this.userScrolled; 
    const previousMessageCount = this.lastMessageCount;
    
    try {
      const token = localStorage.getItem('token');
      const response = await this.$axios.get(
        `/conversations/${this.effectiveConversationId}`, 
        token ? { headers: { Authorization: `Bearer ${token}` } } : {}
      );
      const serverConv = response.data || {};

      const raw = localStorage.getItem('leftConversations') || '{}';
      let leftMap = {};
      try { leftMap = JSON.parse(raw) || {}; } catch { leftMap = {}; }
      const snap = leftMap[this.effectiveConversationId];
      
      if (snap && String(snap.userId || '') === String(this.userId || '')) {
        this.hasLeftGroup = true;
        this.conversation = snap.conversation || serverConv;
      } else {
        this.conversation = serverConv;
      }

      // Aggiorna il conteggio messaggi
      this.lastMessageCount = (this.conversation?.messages || []).length;
      
      // Controlla se ci sono nuovi messaggi
      const hasNewMessages = this.lastMessageCount > previousMessageCount;

      this.$nextTick(() => {
        // Scrolla al bottom solo se:
        // 1. L'utente NON ha scrollato manualmente, oppure
        // 2. Ci sono nuovi messaggi (significa che qualcuno ha inviato)
        if (shouldScrollToBottom || hasNewMessages) {
          this.scrollToBottom();
        }
        this.markRead();
        this.conversation = serverConv;
      });
    } catch (error) {
      console.error('Error loading conversation:', error);
      const raw = localStorage.getItem('leftConversations') || '{}';
      let leftMap = {};
      try { leftMap = JSON.parse(raw) || {}; } catch { leftMap = {}; }
      const snap = leftMap[this.effectiveConversationId];
      
      if (snap && String(snap.userId || '') === String(this.userId || '')) {
        this.conversation = snap.conversation || {};
        this.hasLeftGroup = true;
        this.errorMessage = null;
        if (shouldScrollToBottom) {
          this.$nextTick(() => this.scrollToBottom());
        }
      } 
    } finally {
      this.loading = false;
      this.isLoadingInProgress = false
    }
  },

    async loadConversationsList() {
        try {
            const token = localStorage.getItem('token');
            const res = await this.$axios.get('/conversations', token ? { headers: { Authorization: `Bearer ${token}` } } : {});
            this.allConversations = res.data || [];
      } catch (e) {
        console.error('Failed to load conversations list', e);
      }
    },
    showToast(msg) {
      this.toast = { show: true, msg, targetId: this.toast?.targetId || '' };
      setTimeout(() => { this.toast.show = false; }, 2000);
    },
    openToastTarget() {
      const id = this.toast?.targetId;
      if (id) this.$router.push(`/conversations/${id}`);
      this.toast = { show: false, msg: "", targetId: '' };
    },
    async markRead() {
  try {
    const token = localStorage.getItem('token');
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const msgs = (this.conversation?.messages || []).filter(m => m && m.id && m.senderId !== this.userId);
    
    let anyMarked = false;
    for (const m of msgs) {
      try {
        await this.$axios.post(`/conversations/${this.effectiveConversationId}/messages/${m.id}/status`, { status: 'read' }, { headers });
        anyMarked = true;
      } catch (e) {
      }
    }
  } catch (e) {
    console.error('Error marking messages as read:', e);
  }
},
     async send() {
  if (!this.canSend || this.sending) return
  this.sending = true
  this.errorMessage = null
  
  try {
    const token = localStorage.getItem('token')
    const axiosCfg = token ? { headers: { Authorization: `Bearer ${token}` } } : {}

    let textToSend = this.newMessage || ''
    if (this.reply.active) {
      const quoted = `> ${this.reply.username ? '@' + this.reply.username + ': ' : ''}${this.reply.preview}`
      textToSend = quoted + (textToSend ? `\n\n${textToSend}` : '')
    }

    const trimmedText = (textToSend || '').trim()
    const hasPhoto = !!this.photoAttachB64

    // Send text message
    if (trimmedText) {
      const textPayload = { 
        content: { 
          type: 'text', 
          value: this.toBase64(trimmedText) 
        } 
      }
      
      try {
        const response = await this.$axios.post(
          `/conversations/${this.effectiveConversationId}/messages`,
          textPayload,
          axiosCfg
        )
      } catch (err) {
        throw err
      }
    }

    // Send photo
    if (hasPhoto) {
      
      const photoPayload = { 
        content: { 
          type: 'image', 
          value: this.photoAttachB64 
        } 
      }
      
      try {
        const response = await this.$axios.post(
          `/conversations/${this.effectiveConversationId}/messages`,
          photoPayload,
          axiosCfg
        )
      } catch (err) {
        
        throw new Error('Failed to send photo: ' + (err.response?.data?.error || err.message))
      }
    }

   
    this.newMessage = ''
    this.photoAttachB64 = ''
    if (this.$refs.fileInput) this.$refs.fileInput.value = ''
    this.reply = { active: false, preview: '', username: '' }
    
    this.userScrolled = false
    
  
    await this.load()
    
    
  } catch (error) {
    console.error('💥 SEND ERROR:', error)
    this.errorMessage = error.message || 'Failed to send message'
    this.showToast(this.errorMessage)
  } finally {
    this.sending = false
  }
  
  this.$nextTick(() => {
    this.$refs.messageInput?.focus()
  })
},
    attachPhoto(e) {
  if (this.isGroupChat && !this.isMember) {
    this.photoAttachB64 = ''
    if (e?.target) e.target.value = ''
    this.showToast('You are no longer in this group')
    return
  }
  
  const file = e?.target?.files?.[0]
  if (!file) { 
    this.photoAttachB64 = '' 
    return 
  }
  
  if (!file.type.startsWith('image/')) { 
    this.errorMessage = 'Please select an image' 
    return 
  }
  
  const max = 10 * 1024 * 1024
  if (file.size > max) { 
    console.error('File too large:', file.size, 'max:', max)
    this.errorMessage = 'Image too large (max 10MB)' 
    if (e?.target) e.target.value = ''
    return 
  }
  
  const reader = new FileReader()
  reader.onload = () => {
    const result = reader.result || ''
    
    if (result.includes(',')) {
      this.photoAttachB64 = result.split(',')[1]
    } else {
      this.photoAttachB64 = result
      
    }
  }
  
  reader.onerror = (err) => { 
    this.errorMessage = 'Failed to read file' 
  }
  
  reader.readAsDataURL(file)
},
    statusIcon(status) {
    const s = String(status || '').toLowerCase()
    if (s === 'read') return '✓✓'
    if (s === 'sent') return '✓'
    return '✓' 
  },
  
    statusClass(status) {
    const s = String(status || '').toLowerCase()
    if (s === 'read') return 'read'
    return 'sent'
  },
    resetFileInput(e) {
      if (e?.target) e.target.value = '';
    },
    messageText(message) {
      const type = (message?.content?.type || 'text').toLowerCase();
      if (type !== 'text') return '';
      return this.decodeText(message?.content?.value);
    },
    getReplyInfo(message) {
      const text = this.messageText(message);
      if (!text || !text.startsWith('>')) return null;
      
      const lines = text.split('\n');
      const quoteLine = lines[0];
      const match = quoteLine.match(/^>\s*@?(\S+):\s*(.*)$/);
      
      if (!match) return null;
      
      let preview = match[2];
      const actualText = lines.slice(2).join('\n').trim();
      
      return {
        username: match[1],
        preview: preview,
        actualMessage: actualText
      };
    },
    getActualMessageText(message) {
      const replyInfo = this.getReplyInfo(message);
      return replyInfo ? replyInfo.actualMessage : this.messageText(message);
    },
    isImageMessage(message) {
      const type = (message?.content?.type || '').toLowerCase();
      return type === 'image' && !!(message?.content?.value);
    },
    imageSrc(message) {
      const val = message?.content?.value || '';
      if (!val) return '';
      return val.startsWith('data:') ? val : `data:image/png;base64,${val}`;
    },
    makeLetterAvatar(name, size = 96) {
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
    reactionEmoji(type) {
      const map = { like: '👍', heart: '❤️', laugh: '😂', sad_face: '😢', angry_face: '😡' };
      return map[type] || type;
    },
   async react(messageId, type) {
  if (this.reactBusy[messageId]) return;
  this.$set ? this.$set(this.reactBusy, messageId, true) : (this.reactBusy[messageId] = true);
  try {
    const token = localStorage.getItem('token');
    const url = `/conversations/${this.effectiveConversationId}/messages/${messageId}/reaction`;
    const current = this.myReaction((this.conversation?.messages || []).find(m => m.id === messageId)) || '';
    if (current === type) {
      await this.unreact(messageId, type);
      return;
    }
    if (current) {
      await this.$axios.delete(url + `?type=${current}`, token ? { headers: { Authorization: `Bearer ${token}` } } : {});
    }
    await this.$axios.put(url, { type }, token ? { headers: { Authorization: `Bearer ${token}` } } : {});
    this.tempMyReaction[messageId] = type;
    this.showToast('Reacted');
    this.closePicker();
    await this.load(); 
  } catch (e) {
    console.error('Failed reaction', e);
    this.errorMessage = 'Failed to react';
  } finally {
    this.reactBusy[messageId] = false;
  }
},

async unreact(messageId, type) {
  if (this.reactBusy[messageId]) return;
  this.$set ? this.$set(this.reactBusy, messageId, true) : (this.reactBusy[messageId] = true);
  try {
    const token = localStorage.getItem('token');
    const url = `/conversations/${this.effectiveConversationId}/messages/${messageId}/reaction?type=${type}`;
    await this.$axios.delete(url, token ? { headers: { Authorization: `Bearer ${token}` } } : {});
    delete this.tempMyReaction[messageId];
    this.showToast('Removed reaction');
    this.closePicker();
    await this.load(); // ← CAMBIATO: tolto setTimeout, messo await
  } catch (e) {
    console.error('Failed to remove reaction', e);
    this.errorMessage = 'Failed to remove reaction';
  } finally {
    this.reactBusy[messageId] = false;
  }
},
   async del(messageId) {
    if (!messageId) return;
    const pairedId = this.findPairedReplyStub(messageId) || null;
    this.deleteDialog = { open: true, messageId, pairedId };
    },
    async confirmDelete() {
    const messageId = this.deleteDialog.messageId;
    const pairedId = this.deleteDialog.pairedId;
    this.deleteDialog = { open: false, messageId: null, pairedId: null };
    if (!messageId) return;
    const token = localStorage.getItem('token');
    const cfg = token ? { headers: { Authorization: `Bearer ${token}` } } : {};
    let ok = true;
    try {
      await this.$axios.delete(`/conversations/${this.effectiveConversationId}/messages/${messageId}`, cfg);
    } catch (e) {
      ok = false;
      console.error('Failed to delete message', e);
      this.errorMessage = 'Failed to delete message';
    }
    if (pairedId) {
      try {
        await this.$axios.delete(`/conversations/${this.effectiveConversationId}/messages/${pairedId}`, cfg);
      } catch (e) {
        ok = false;
        console.error('Failed to delete paired message', e);
      }
    }
    this.showToast(pairedId ? (ok ? 'Message and reply deleted.' : 'Delete partially failed.') : (ok ? 'Message deleted.' : 'Delete failed.'));
    await this.load();
    },
    cancelDelete() {
      this.deleteDialog = { open: false, messageId: null, pairedId: null };
    },

    openForward(messageId) {
      try {
        this.$router.push({ path: '/home', query: { forward: messageId, from: this.effectiveConversationId } });
      } catch {}
    },
    openReply(message) {
      const actualText = (this.getActualMessageText(message) || '').trim();
      const hasImage = this.isImageMessage(message);

        let preview = '';
        if (hasImage && actualText) {
          const suffix = ' • Photo';
          const max = 80 - suffix.length;
          const truncated = actualText.length > max ? (actualText.slice(0, Math.max(0, max - 3)) + '...') : actualText;
          preview = truncated + suffix;
        } else if (hasImage) {
          preview = 'Photo';
        } else {
          preview = actualText;
        }

        const username = message?.sender?.username || '';
        this.reply = { active: true, preview, username };
        this.$nextTick(() => this.$refs.messageInput?.focus());
    },
    cancelReply() {
        this.reply = { active: false, preview: '', username: '' };
    },
    removeMyReaction(message) {
      try {
        const type = this.myReaction(message);
        if (!type) return;
        this.unreact(message.id, type);
      } catch {}
    },
    closeForward() {
        this.forward = { open: false, messageId: null, target: '', newUsername: '', suggestions: [], suggestTimer: null, suggestLoading: false, selectedUserId: '' };
    },
    openUserProfileModal() {
  if (this.isGroupChat) return;

  try {
    const myUserId = String(this.userId || '');
    const participants = this.conversation?.participants || [];
    const peer = participants.find(p => {
      const pid = String(p?.id || p?.userId || '');
      return pid && pid !== myUserId;
    });

    if (!peer) {
      console.warn('Peer user not found');
      return;
    }

    const username =
      peer.username ||
      peer.Username ||
      peer.name ||
      'User';

    const photoB64 =
      peer.photo ||
      peer.Photo ||
      peer.profilePhoto ||
      peer.ProfilePhoto ||
      null;

    this.userProfileModal = {
      open: true,
      username,
      photo: photoB64
        ? (photoB64.startsWith('data:')
            ? photoB64
            : `data:image/png;base64,${photoB64}`)
        : this.makeLetterAvatar(username),
      description: 'Hey, this is my profile!'
    };
  } catch (e) {
    console.error('Failed to open user profile modal:', e);
  }
},

    closeUserProfileModal() {
        this.userProfileModal = { open: false, username: '', photo: null, description: '' };
    },
    async doForwardNew() {
        if (!this.forward.messageId || !this.canForwardNew) return;
        this.forwarding = true;
        try {
            const username = (this.forward.newUsername || '').trim();
            const token = localStorage.getItem('token');
            const headers = token ? { Authorization: `Bearer ${token}` } : {};

            let userId = this.forward.selectedUserId || '';
            if (!userId) {
              const sr = await this.$axios.get(`/searchby?user=${encodeURIComponent(username)}`, { headers });
              const users = (sr?.data?.users || []).filter(u => u && u.id && (u.username || '').toLowerCase() === username.toLowerCase());
              userId = users?.[0]?.id || '';
            }
            if (!userId) {
                this.errorMessage = 'User not found';
                return;
            }

            const cd = await this.$axios.post('/direct-conversations', { peerUserId: userId }, { headers });
            const targetId = cd?.data?.conversationId;
            if (!targetId) {
                this.errorMessage = 'Failed to create conversation';
                return;
            }

            await this.$axios.post(
              `/conversations/${this.effectiveConversationId}/messages/${this.forward.messageId}/forward`,
              { targetConversationId: targetId },
              { headers }
            );
            this.toast = { show: true, msg: 'Message forwarded.', targetId };
            setTimeout(() => { this.toast.show = false; }, 2000);
            this.closeForward();
        } catch (e) {
            console.error('Failed to forward to new user', e);
            this.errorMessage = 'Failed to forward';
        } finally {
            this.forwarding = false;
            this.$nextTick(() => { this.$refs.messageInput?.focus(); });
        }
    },
   async selectForwardUser(user) {
  if (!user || !this.forward.messageId) return
  
  this.forwarding = true
  
  try {
    const token = localStorage.getItem('token')
    const headers = token ? { Authorization: `Bearer ${token}` } : {}
    
    // 🆕 STEP 1: Ricarica la lista conversazioni PRIMA di cercare
    console.log('📋 Reloading conversations list before search...')
    await this.loadConversationsList()
    
    console.log('📊 Total conversations:', this.allConversations.length)
    console.log('🔍 Looking for existing chat with user:', user.id, user.username)
    
    // 🆕 STEP 2: Debug - mostra tutte le conversazioni
    this.allConversations.forEach((conv, idx) => {
      const participants = conv.participants || []
      const participantIds = participants.map(p => String(p?.id || p?.userId || ''))
      console.log(`  Conv ${idx}:`, {
        id: conv.id || conv.conversationId,
        isGroup: conv.isGroup,
        participantIds: participantIds
      })
    })
    
    // 🔍 CERCA CONVERSAZIONE ESISTENTE
    const myUserId = String(this.userId || '')
    const targetUserId = String(user.id)
    
    const existingConv = this.allConversations.find(conv => {
      // Ignora gruppi
      if (conv.isGroup === true) {
        console.log('  ⏭️  Skipping group:', conv.id)
        return false
      }
      
      const participants = conv.participants || []
      
      // Debug participants
      console.log('  🔎 Checking conv:', conv.id || conv.conversationId, 'participants:', participants.length)
      
      // Controlla se questa conversazione contiene ESATTAMENTE 2 persone: me e target
      const participantIds = participants.map(p => String(p?.id || p?.userId || ''))
      
      const hasMe = participantIds.includes(myUserId)
      const hasTarget = participantIds.includes(targetUserId)
      const isTwoPersonChat = participants.length === 2
      
      console.log('    hasMe:', hasMe, 'hasTarget:', hasTarget, 'isTwoPersonChat:', isTwoPersonChat)
      
      return hasMe && hasTarget && isTwoPersonChat
    })
    
    let targetId
    
    if (existingConv) {
      // ✅ USA LA CONVERSAZIONE ESISTENTE
      targetId = existingConv.id || existingConv.conversationId
      console.log('✅ FOUND existing conversation:', targetId)
    } else {
      // ❌ CREA NUOVA CONVERSAZIONE
      console.log('❌ NO existing conversation found, creating new...')
      const response = await this.$axios.post('/direct-conversations', { peerUserId: user.id }, { headers })
      const conv = response?.data || {}
      targetId = conv.id || conv.ID || conv.conversationId || conv.conversation_id
      console.log('🆕 Created new conversation:', targetId)
    }
    
    if (!targetId) {
      console.error('❌ Failed to get targetId')
      this.errorMessage = 'Failed to get conversation'
      return
    }
    
    // 🔄 FORWARD IL MESSAGGIO
    console.log('📤 Forwarding message:')
    console.log('  From:', this.effectiveConversationId)
    console.log('  To:', targetId)
    console.log('  MessageId:', this.forward.messageId)
    
    await this.$axios.post(
      `/conversations/${this.effectiveConversationId}/messages/${this.forward.messageId}/forward`,
      { targetConversationId: targetId },
      { headers }
    )
    
    console.log('✅ Forward successful!')
    
    this.closeForward()
    
    // Redirect alla conversazione target
    this.$router.push(`/conversations/${targetId}`)
    
  } catch (error) {
    console.error('❌ Failed to forward to user', error)
    console.error('Error response:', error.response?.data)
    this.errorMessage = error.response?.data?.error || error.response?.data?.message || 'Failed to forward message'
  } finally {
    this.forwarding = false
  }
},
    async forwardToConversation(conversationId) {
        if (!this.forward.messageId || !conversationId) return;
        this.forwarding = true;
        try {
            const token = localStorage.getItem('token');
            await this.$axios.post(
              `/conversations/${this.effectiveConversationId}/messages/${this.forward.messageId}/forward`,
              { targetConversationId: conversationId },
              token ? { headers: { Authorization: `Bearer ${token}` } } : {}
            );
            this.closeForward();
            if (this.$route.path !== `/conversations/${conversationId}`) {
              this.$router.push(`/conversations/${conversationId}`);
            } else {
              await this.load();
            }
        } catch (e) {
            console.error('Failed to forward message', e);
            this.errorMessage = 'Failed to forward message';
        } finally {
            this.forwarding = false;
        }
    },
    onForwardUserInput() {
      this.forward.selectedUserId = '';
      const q = (this.forward.newUsername || '').trim();
      if (this.forward.suggestTimer) clearTimeout(this.forward.suggestTimer);
      if (q.length < 1) { this.forward.suggestions = []; this.forward.suggestLoading = false; return; }
      this.forward.suggestLoading = true;
      this.forward.suggestTimer = setTimeout(async () => {
        try {
          const token = localStorage.getItem('token');
          const headers = token ? { Authorization: `Bearer ${token}` } : {};
          const res = await this.$axios.get(`/searchby?user=${encodeURIComponent(q)}`, { headers });
          const all = res?.data?.users || [];
          this.forward.suggestions = all.filter(u => u && u.id && u.id !== this.userId).slice(0, 10);
        } catch { this.forward.suggestions = []; }
        finally { this.forward.suggestLoading = false; }
      }, 250);
    },
    isOwn(message) {
      try {
        const sid = message?.sender?.id || message?.sender?.ID || message?.senderId || '';
        return sid === this.userId;
      } catch { return false; }
    },
    formatTime(timestamp) {
        if (!timestamp) return '';
      try {
        const date = new Date(timestamp);
        return date.toLocaleTimeString('it-IT', { hour: '2-digit', minute: '2-digit' });
      } catch { return timestamp; }
    },
    getMessageDate(timestamp) {
      if (!timestamp) return '';
      try {
        const date = new Date(timestamp);
        return date.toDateString();
      } catch { return ''; }
    },
    getCreationTimestamp() {
      try {
        const ts = this.conversation?.createdAt || this.conversation?.CreatedAt || (this.conversation?.messages?.[0]?.timestamp || null);
        return ts || null;
      } catch { return null; }
    },
    formatDate(timestamp) {
      if (!timestamp) return '';
      try {
        const date = new Date(timestamp);
        const today = new Date();
        const yesterday = new Date(today);
        yesterday.setDate(yesterday.getDate() - 1);

        const isSameDay = (d1, d2) => 
          d1.getDate() === d2.getDate() && 
          d1.getMonth() === d2.getMonth() && 
          d1.getFullYear() === d2.getFullYear();

        if (isSameDay(date, today)) {
          return 'Oggi';
        }

        if (isSameDay(date, yesterday)) {
          return 'Ieri';
        }

        const monthNames = ['Gennaio', 'Febbraio', 'Marzo', 'Aprile', 'Maggio', 'Giugno',
                           'Luglio', 'Agosto', 'Settembre', 'Ottobre', 'Novembre', 'Dicembre'];
        
        const day = date.getDate();
        const month = monthNames[date.getMonth()];
        const year = date.getFullYear();

        if (year === today.getFullYear()) {
          return `${day} ${month}`;
        } else {
          return `${day} ${month} ${year}`;
        }
      } catch { return timestamp; }
    },
    toBase64(str) {
        const bytes = new TextEncoder().encode(str);
      let binary = '';
      bytes.forEach(b => binary += String.fromCharCode(b));
      return btoa(binary);
    },
    decodeText(b64) {
      if (!b64) return '';
      try {
        const bin = atob(b64);
        const bytes = Uint8Array.from(bin, c => c.charCodeAt(0));
        return new TextDecoder().decode(bytes);
      } catch { return ''; }
    },
    findPairedReplyStub(imageMessageId) {
      try {
        const arr = this.conversation?.messages || [];
        const idx = arr.findIndex(m => m && m.id === imageMessageId);
        if (idx <= 0) return '';
        const imageMsg = arr[idx];
        const prev = arr[idx - 1];
        if (!prev) return '';
        const sameSender = (prev?.sender?.id || prev?.senderId) === (imageMsg?.sender?.id || imageMsg?.senderId);
        const isText = String(prev?.content?.type || 'text').toLowerCase() === 'text';
        const ri = this.getReplyInfo(prev);
        const replyOnly = !!ri && (!ri.actualMessage || ri.actualMessage.length === 0);
        const near = this.isNearSameDay(prev?.timestamp, imageMsg?.timestamp);
        if (sameSender && isText && replyOnly && near) return prev.id;
        return '';
      } catch { return ''; }
    },
    isNearSameDay(ts1, ts2) {
      try {
        const d1 = new Date(ts1);
        const d2 = new Date(ts2);
        const sameDay = d1.getFullYear() === d2.getFullYear() && d1.getMonth() === d2.getMonth() && d1.getDate() === d2.getDate();
        const diff = Math.abs(d2.getTime() - d1.getTime());
        return sameDay && diff <= 2 * 60 * 1000;
      } catch { return false; }
    },
    onScroll() {
    const el = this.$refs.scrollArea;
    if (!el) return;
    const isAtBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 100;
    this.userScrolled = !isAtBottom;
    },
    scrollToBottom() {
    const el = this.$refs.scrollArea;
    if (el) {
      el.scrollTop = el.scrollHeight;
      this.userScrolled = false;
    }
  },
    openFilePicker() {
      if (this.isGroupChat && !this.isMember) {
        this.showToast('You are no longer in this group');
        return;
      }
      this.$refs.fileInput?.click();
    },
    togglePicker(messageId) {
      const cur = this.reactionPicker.openFor || '';
      this.reactionPicker.openFor = cur === messageId ? '' : messageId;
    },
    closePicker() {
      this.reactionPicker.openFor = '';
    },
    async closeGroupEdit(payload) {
      this.groupEditOpen = false;
      if (payload?.updated) {
        this.loading = true;
        await this.load();
        this.$emit('group-updated', { convId: this.effectiveConversationId });
      }
    },
    onOutsideClick(e) {
      try {
        if (!this.reactionPicker.openFor) return;
        const el = this.$el;
        const t = e?.target;
        if (!el || !t) { this.closePicker(); return; }
        const insideTrigger = t.closest && (t.closest('.reaction-trigger') || t.closest('.reaction-picker'));
        if (insideTrigger && el.contains(insideTrigger)) return;
        this.closePicker();
      } catch { this.closePicker(); }
    },
  },

  watch: {
    effectiveConversationId(newId, oldId) {
      if (newId && newId !== oldId) {
        this.load()
      }
    }
  },

  mounted() {
  this.load();
  this.loadConversationsList(); 
  // Aggiungi listener per lo scroll
  this.$nextTick(() => {
    const scrollArea = this.$refs.scrollArea;
    if (scrollArea) {
      scrollArea.addEventListener('scroll', this.onScroll);
    }
  });
  
  try { window.addEventListener('keydown', this.onKeyDown); } catch {}
  try { window.addEventListener('click', this.onOutsideClick, { capture: true }); } catch {}
  this.pollId = setInterval(() => this.load(), 2000); 
  
  this.$nextTick(() => {
    this.$refs.messageInput?.focus();
  });
},


  unmounted() {
  if (this.pollId) clearInterval(this.pollId);
  
  // Rimuovi listener per lo scroll
  const scrollArea = this.$refs.scrollArea;
  if (scrollArea) {
    scrollArea.removeEventListener('scroll', this.onScroll);
  }
  
  try { window.removeEventListener('keydown', this.onKeyDown); } catch {}
  try { window.removeEventListener('click', this.onOutsideClick, { capture: true }); } catch {}
}
}
</script>
<style scoped>


.conversation-neon {
  position: relative;
  display: grid;
  grid-template-rows: auto 1fr auto; 
  height: 100vh; 
  width: 100%;
  background: var(--bg-primary);
  color: var(--text-primary);
  overflow: hidden;
}

.conv-header-neon {
  position: sticky;
  top: 70px; 
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem 1rem;
  min-height: 70px;
  max-height: 70px;
  background: #ffffff;
  border-bottom: 2px solid rgba(26, 115, 232, 0.25);
  backdrop-filter: blur(10px);
  flex-shrink: 0;
}

.back-btn-neon {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.4rem 0.8rem; 
  background: rgba(26, 115, 232, 0.08);
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-radius: 999px;
  color: var(--neon-cyan);
  text-decoration: none;
  font-weight: 600;
  font-size: 0.875rem; 
  transition: all 0.3s ease;
}

.back-btn-neon:hover {
  background: rgba(26, 115, 232, 0.15);
  border-color: var(--neon-cyan);
  box-shadow: 0 0 15px rgba(26, 115, 232, 0.30);
  transform: translateX(-3px);
}

.peer-info-neon {
  display: flex;
  align-items: center;
  gap: 0.75rem; 
  flex: 1;
  min-width: 0;
}

.avatar-neon {
  width: 42px; 
  height: 42px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--neon-cyan); 
  box-shadow: 0 0 12px rgba(26, 115, 232, 0.30);
  cursor: pointer;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.avatar-neon:hover {
  transform: scale(1.05); 
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.45);
}

.meta-neon {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.15rem; 
}

.name-neon {
  margin: 0;
  font-size: 1rem; 
  font-weight: 700;
  background: linear-gradient(90deg, var(--neon-cyan) 0%, var(--neon-purple) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.name-neon:hover {
  filter: drop-shadow(0 0 8px rgba(26, 115, 232, 0.40));
}

.subtitle-neon {
  color: var(--text-secondary);
  font-size: 0.75rem;
}

.quick-actions-neon {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.25rem;
}

.header-actions-neon {
  display: flex;
  gap: 0.5rem; 
  flex-shrink: 0;
}

.btn-neon {
  padding: 0.4rem 0.8rem; 
  font-size: 0.875rem;
  border-radius: 999px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.btn-neon.danger {
  background: rgba(255, 0, 100, 0.15);
  border-color: rgba(255, 0, 100, 0.4);
  color: #ff0064;
}

.btn-neon.danger:hover {
  background: rgba(255, 0, 100, 0.25);
  box-shadow: 0 0 15px rgba(255, 0, 100, 0.4);
}

.btn-neon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Messages Area*/
.messages-neon {
  position: absolute;
  top: 142px; 
  bottom: 80px; 
  left: 0;
  right: 0;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  background: transparent;
  z-index: 10;
}


.messages-neon::-webkit-scrollbar {
  width: 8px;
}

.messages-neon::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.06);
  border-radius: 10px;
}

.messages-neon::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, var(--neon-cyan), var(--neon-purple));
  border-radius: 10px;
  border: 2px solid rgba(0, 0, 0, 0.06);
}

.messages-neon::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, var(--neon-purple), var(--neon-cyan));
}

.loading-spinner-neon {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 2rem;
  color: var(--neon-cyan);
}

.spinner-neon {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(26, 115, 232, 0.15);
  border-top-color: var(--neon-cyan);
  border-radius: 50%;
  animation: spin-neon 1s linear infinite;
}
.photo-preview-banner-neon {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  padding: 0.6rem 0.75rem;
  background: rgba(138, 43, 226, 0.08);
  border-left: 3px solid var(--neon-purple);
  border-radius: 10px;
  color: var(--text-primary);
}

.photo-preview-content-neon {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 1;
  min-width: 0;
}

.photo-preview-label-neon {
  font-size: 0.75rem;
  color: var(--neon-purple);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.photo-preview-thumbnail-neon {
  width: 50px;
  height: 50px;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid rgba(26, 115, 232, 0.20);
  flex-shrink: 0;
}

.photo-preview-thumbnail-neon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.photo-preview-close-neon {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 1.5rem;
  padding: 0;
  line-height: 1;
  flex-shrink: 0;
  transition: all 0.3s ease;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.photo-preview-close-neon:hover {
  color: #ff0064;
  background: rgba(255, 0, 100, 0.1);
  transform: rotate(90deg);
}

@keyframes spin-neon {
  to { transform: rotate(360deg); }
}

.error-neon {
  padding: 0.875rem 1rem;
  background: rgba(255, 0, 100, 0.15);
  border: 1px solid rgba(255, 0, 100, 0.4);
  border-radius: 10px;
  color: #ff0064;
  text-align: center;
  font-size: 0.875rem;
}

.toast-neon {
  position: fixed;
  bottom: 5rem;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  padding: 0.75rem 1.25rem;
  background: #ffffff;
  border: 2px solid var(--neon-cyan);
  border-radius: 10px;
  color: var(--neon-cyan);
  font-weight: 600;
  font-size: 0.875rem;
  box-shadow: 0 0 25px rgba(26, 115, 232, 0.40);
  backdrop-filter: blur(10px);
  animation: toast-slide-up 0.3s ease;
}

@keyframes toast-slide-up {
  from {
    transform: translateX(-50%) translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateX(-50%) translateY(0);
    opacity: 1;
  }
}

/* Date Separators */
.date-separator-neon {
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 0.75rem 0;
}

.date-badge-neon {
  padding: 0.4rem 1rem;
  background: #f8faff;
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-radius: 16px;
  color: var(--neon-cyan);
  font-size: 0.8rem;
  font-weight: 600;
  box-shadow: 0 0 15px rgba(26, 115, 232, 0.15);
  backdrop-filter: blur(8px);
}

/* Message Row */
.msg-row-neon {
  display: flex;
  justify-content: flex-start;
  margin-bottom: 1rem;
  position: relative;
  animation: fadeInUp 0.3s ease;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.msg-row-neon.own {
  justify-content: flex-end;
}

/* Message Bubble */
.bubble-neon {
  max-width: 70ch;
  padding: 0.75rem 0.875rem;
  border-radius: 14px;
  background: #f8faff;
  border: 2px solid rgba(26, 115, 232, 0.20);
  box-shadow: 0 3px 15px rgba(26, 115, 232, 0.06);
  backdrop-filter: blur(10px);
  position: relative;
  z-index: 2;
  transition: all 0.3s ease;
}

.msg-row-neon.own .bubble-neon {
  background: rgba(26, 115, 232, 0.10);
  border-color: rgba(26, 115, 232, 0.30);
  box-shadow: 0 3px 15px rgba(26, 115, 232, 0.12);
}

.bubble-neon:hover {
  transform: translateY(-1px);
  box-shadow: 0 5px 20px rgba(26, 115, 232, 0.06);
}

.msg-row-neon.own .bubble-neon:hover {
  box-shadow: 0 5px 20px rgba(26, 115, 232, 0.20);
}

/* Sender Name */
.sender-name-neon {
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--neon-cyan);
  margin-bottom: 0.4rem;
  display: block;
}

.msg-row-neon.own .sender-name-neon {
  color: var(--neon-purple);
}

/* Forwarded Badge */
.fwd-neon {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.3rem 0.6rem;
  background: rgba(26, 115, 232, 0.08);
  border-left: 3px solid var(--neon-purple);
  border-radius: 6px;
  font-size: 0.75rem;
  color: var(--neon-purple);
  margin-bottom: 0.4rem;
}

.fwd-icon {
  width: 12px;
  height: 12px;
  color: var(--neon-purple);
}

/* Reply Container */
.reply-container-neon {
  padding: 0.6rem 0.8rem;
  margin-bottom: 0.6rem;
  background: rgba(26, 115, 232, 0.06);
  border-left: 3px solid var(--neon-cyan);
  border-radius: 8px;
  backdrop-filter: blur(5px);
}

.reply-username-neon {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--neon-cyan);
  margin-bottom: 0.3rem;
}

.reply-message-neon {
  font-size: 0.8rem;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Text & Attachments */
.text-neon {
  white-space: pre-wrap;
  word-break: break-word;
  color: var(--text-primary);
  line-height: 1.4;
  font-size: 0.925rem;
}

.attachment-neon img {
  max-width: 380px;
  border-radius: 10px;
  display: block;
  border: 2px solid rgba(26, 115, 232, 0.20);
}

/* Meta Line */
.meta-line-neon {
  margin-top: 0.4rem;
  font-size: 0.7rem;
  color: var(--text-secondary);
  display: flex;
  justify-content: flex-end;
  gap: 0.4rem;
  align-items: center;
}

.time-neon {
  color: var(--text-secondary);
}

.status-neon {
  margin-left: 0.2rem;
}

.status-neon.sent {
  color: var(--text-secondary);
}

.status-neon.read {
  color: var(--neon-cyan);
  text-shadow: 0 0 6px rgba(26, 115, 232, 0.40);
}

/* Actions */
.actions-neon {
  margin-top: 0.6rem;
  padding-top: 0.4rem;
  border-top: 1px solid rgba(26, 115, 232, 0.12);
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.link-neon {
  background: transparent;
  border: none;
  color: var(--neon-cyan);
  cursor: pointer;
  padding: 0;
  font-size: 0.8rem;
  font-weight: 600;
  transition: all 0.3s ease;
}

.link-neon:hover {
  color: var(--neon-purple);
  text-shadow: 0 0 8px rgba(138, 43, 226, 0.5);
}

.link-neon.danger {
  color: #ff0064;
}

.link-neon.danger:hover {
  text-shadow: 0 0 8px rgba(255, 0, 100, 0.5);
}

.link-neon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Reaction Trigger */
.reaction-trigger-neon {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  border: 2px solid rgba(26, 115, 232, 0.30);
  background: #ffffff;
  color: var(--neon-cyan);
  font-weight: 800;
  font-size: 1.15rem;
  cursor: pointer;
  opacity: 0;
  transition: all 0.25s ease;
  z-index: 3;
  backdrop-filter: blur(5px);
}

.msg-row-neon:hover .reaction-trigger-neon {
  opacity: 1;
}

.reaction-trigger-neon.open {
  opacity: 1;
  background: rgba(26, 115, 232, 0.15);
  box-shadow: 0 0 12px rgba(26, 115, 232, 0.30);
}

.reaction-trigger-neon.left {
  left: -36px;
}

.reaction-trigger-neon.right {
  right: -36px;
}

.reaction-trigger-neon:hover {
  background: rgba(26, 115, 232, 0.20);
  box-shadow: 0 0 16px rgba(26, 115, 232, 0.40);
  transform: translateY(-50%) scale(1.08);
}

/* Reaction Picker */
.reaction-picker-neon {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  background: #ffffff;
  border: 2px solid rgba(26, 115, 232, 0.30);
  border-radius: 10px;
  padding: 0.6rem 0.4rem;
  box-shadow: 0 6px 25px rgba(26, 115, 232, 0.08);
  z-index: 4;
  backdrop-filter: blur(10px);
}

.reaction-picker-neon.left {
  left: -36px;
  transform: translate(-100%, -50%);
}

.reaction-picker-neon.right {
  right: -36px;
  transform: translate(100%, -50%);
}

.rx-neon {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 1.4rem;
  line-height: 1;
  padding: 0.4rem;
  border-radius: 8px;
  transition: all 0.25s ease;
}

.rx-neon:hover {
  background: rgba(26, 115, 232, 0.15);
  transform: scale(1.15);
}

.rx-neon.active {
  background: rgba(26, 115, 232, 0.22);
  box-shadow: 0 0 12px rgba(26, 115, 232, 0.30);
}

/* Reactions Container */
.reactions-container-neon {
  position: absolute;
  bottom: -28px;
  right: 10px;
  background: #ffffff;
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-radius: 18px;
  padding: 0.3rem 0.65rem;
  display: flex;
  align-items: center;
  gap: 0.4rem;
  box-shadow: 0 3px 12px rgba(26, 115, 232, 0.06);
  z-index: 2;
  backdrop-filter: blur(10px);
}

.reaction-group-wrapper-neon {
  position: relative;
  display: inline-flex;
  align-items: center;
}

.reaction-group-neon {
  font-size: 0.95rem;
  line-height: 1;
  cursor: pointer;
  position: relative;
  transition: transform 0.2s ease;
}

.reaction-group-neon:hover {
  transform: scale(1.25);
}

.reaction-group-neon.mine {
  filter: drop-shadow(0 0 6px rgba(26, 115, 232, 0.40));
}

.reaction-tooltip-neon {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  background: #ffffff;
  border: 2px solid rgba(26, 115, 232, 0.30);
  border-radius: 8px;
  padding: 0.55rem 0.75rem;
  margin-bottom: 0.6rem;
  white-space: nowrap;
  font-size: 0.75rem;
  color: #ffffff;
  z-index: 10;
  box-shadow: 0 5px 18px rgba(26, 115, 232, 0.06);
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s ease;
  backdrop-filter: blur(10px);
}

.reaction-group-wrapper-neon:hover .reaction-tooltip-neon {
  opacity: 1;
  pointer-events: auto;
}

.tooltip-label {
  font-weight: 700;
  color: var(--neon-cyan);
  margin-bottom: 0.3rem;
}

.tooltip-usernames {
  color: var(--text-secondary);
}

.reaction-count-neon {
  color: var(--text-secondary);
  font-size: 0.8rem;
  font-weight: 600;
  padding-left: 0.2rem;
}

.composer-neon {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 50;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: #ffffff;
  border-top: 2px solid rgba(26, 115, 232, 0.25);
  backdrop-filter: blur(10px);
}
.reply-banner-neon {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 0.75rem;
  padding: 0.6rem 0.75rem;
  background: rgba(26, 115, 232, 0.06);
  border-left: 3px solid var(--neon-cyan);
  border-radius: 10px;
  color: var(--text-primary);
}

.reply-content-neon {
  flex: 1;
  min-width: 0;
}

.reply-label-neon {
  font-size: 0.7rem;
  color: var(--text-secondary);
  margin-bottom: 0.2rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.reply-close-neon {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 1.15rem;
  padding: 0;
  line-height: 1;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.reply-close-neon:hover {
  color: #ff0064;
  transform: rotate(90deg);
}

.input-row-neon {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.input-neon {
  flex: 1;
  padding: 0.75rem 1rem;
  background: rgba(26, 115, 232, 0.06);
  border: 2px solid rgba(26, 115, 232, 0.25);
  border-radius: 12px;
  color: var(--text-primary);
  font-size: 0.925rem;
  transition: all 0.3s ease;
  outline: none;
}

.input-neon::placeholder {
  color: var(--text-dim);
}

.input-neon:focus {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 15px rgba(26, 115, 232, 0.25);
  background: #ffffff;
}

.input-neon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.file-input-hidden {
  display: none;
}

.attach-btn-neon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  border: 2px solid rgba(26, 115, 232, 0.25);
  background: rgba(26, 115, 232, 0.06);
  color: var(--neon-cyan);
  cursor: pointer;
  transition: all 0.3s ease;
}

.attach-btn-neon:hover {
  background: rgba(26, 115, 232, 0.15);
  box-shadow: 0 0 12px rgba(26, 115, 232, 0.25);
}

.attach-btn-neon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.attach-icon-neon {
  width: 20px;
  height: 20px;
  stroke: var(--neon-cyan);
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
  fill: none;
}

.btn-neon.send {
  padding: 0.75rem 1.25rem;
  background: var(--neon-cyan);
  border: none;
  border-radius: 10px;
  color: var(--bg-primary);
  font-weight: 700;
  font-size: 0.875rem;
  cursor: pointer;
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.30);
  transition: all 0.3s ease;
}

.btn-neon.send:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 0 28px rgba(26, 115, 232, 0.45);
}

.btn-neon.send:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

/* Modals */
.modal-overlay-neon {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.30);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
  backdrop-filter: blur(6px);
}

.modal-card-neon {
  width: 90%;
  max-width: 480px;
  max-height: 80vh;
  overflow-y: auto;
  background: #ffffff;
  border: 2px solid rgba(26, 115, 232, 0.30);
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.20);
  backdrop-filter: blur(10px);
}

.modal-card-neon.small {
  max-width: 380px;
}

.modal-card-neon.large {
  max-width: 680px;
}

.modal-title-neon {
  color: var(--neon-cyan);
  font-size: 1.25rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.modal-text-neon {
  color: var(--text-secondary);
  font-size: 0.925rem;
  line-height: 1.5;
  margin-bottom: 1.5rem;
}

.modal-actions-neon {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
}

.btn-neon.secondary {
  background: transparent;
  border: 2px solid rgba(26, 115, 232, 0.30);
  color: var(--neon-cyan);
}

.btn-neon.secondary:hover {
  background: rgba(26, 115, 232, 0.12);
  box-shadow: 0 0 15px rgba(26, 115, 232, 0.25);
}


.forward-search-neon {
  margin-bottom: 1rem;
}

.forward-suggestions-neon {
  margin-top: 0.5rem;
  max-height: 220px;
  overflow-y: auto;
  background: rgba(26, 115, 232, 0.06);
  border: 1px solid rgba(26, 115, 232, 0.20);
  border-radius: 10px;
}

.forward-item-neon {
  padding: 0.75rem 1rem;
  cursor: pointer;
  border-bottom: 1px solid rgba(26, 115, 232, 0.12);
  color: var(--text-primary);
  transition: background 0.2s ease;
  font-size: 0.925rem;
}

.forward-item-neon:last-child {
  border-bottom: none;
}

.forward-item-neon:hover {
  background: rgba(26, 115, 232, 0.12);
}

.hint-neon {
  color: var(--text-secondary);
  font-size: 0.8rem;
  margin-top: 0.5rem;
}

.recent-chats-neon h4 {
  color: var(--text-secondary);
  font-size: 0.85rem;
  margin: 1rem 0 0.5rem 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* Participants List */
.participants-list-neon {
  display: grid;
  gap: 0.75rem;
  max-height: 400px;
  overflow-y: auto;
}

.participant-item-neon {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 0.75rem;
  background: rgba(26, 115, 232, 0.06);
  border: 1px solid rgba(26, 115, 232, 0.15);
  border-radius: 10px;
  transition: background 0.2s ease;
}

.participant-item-neon:hover {
  background: rgba(26, 115, 232, 0.12);
}

.participant-avatar-neon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  background: rgba(26, 115, 232, 0.06);
  flex-shrink: 0;
  border: 2px solid rgba(26, 115, 232, 0.25);
}

.participant-avatar-neon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.participant-name-neon {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 0.925rem;
}

/* Profile Modal */
.profile-modal-neon {
  position: relative;
  background: #ffffff;
  border: 2px solid rgba(26, 115, 232, 0.30);
  border-radius: 16px;
  padding: 2rem 1.5rem;
  width: min(320px, 90vw);
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.20);
  backdrop-filter: blur(10px);
}

.close-btn-neon {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 1.75rem;
  cursor: pointer;
  padding: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.close-btn-neon:hover {
  background: rgba(255, 0, 100, 0.2);
  color: #ff0064;
  transform: rotate(90deg);
}

.profile-content-neon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  text-align: center;
}

.profile-photo-neon {
  width: 110px;
  height: 110px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid var(--neon-cyan);
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.40);
}

.profile-photo-default-neon {
  width: 110px;
  height: 110px;
  border-radius: 50%;
  background: rgba(26, 115, 232, 0.12);
  border: 3px solid var(--neon-cyan);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.75rem;
  font-weight: 700;
  color: var(--neon-cyan);
  box-shadow: 0 0 20px rgba(26, 115, 232, 0.40);
}

.profile-username-neon {
  margin: 0.5rem 0 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--neon-cyan);
}

.profile-description-neon {
  margin: 0.5rem 0 0;
  font-size: 0.925rem;
  color: var(--text-secondary);
  line-height: 1.4;
  font-style: italic;
}
</style>
