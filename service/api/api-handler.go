package api

import "net/http"

// Handler returns the HTTP handler used by the server
func (rt *_router) Handler() http.Handler {

	// --- AUTH ---
	rt.router.POST("/login", rt.wrap(rt.doLogin))

	// --- USERS ---
	rt.router.GET("/users", rt.wrap(rt.getAllUsers))
	rt.router.GET("/searchby", rt.wrap(rt.searchBy))
	rt.router.PUT("/user/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/user/photo", rt.wrap(rt.setMyPhoto))

	// --- CONVERSATIONS ---
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	rt.router.GET("/conversations/:conversationId", rt.wrap(rt.getConversation))
	rt.router.GET("/conversations/:conversationId/members", rt.wrap(rt.getConversationMembers))
	rt.router.DELETE("/conversations/:conversationId", rt.wrap(rt.deleteConversation))

	// --- MESSAGES ---
	rt.router.POST("/conversations/:conversationId/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/conversations/:conversationId/messages/:messageId/forward", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId", rt.wrap(rt.deleteMessage))
	rt.router.POST("/conversations/:conversationId/messages/:messageId/status", rt.wrap(rt.setMessageStatus))

	// --- REACTIONS ---
	rt.router.PUT("/conversations/:conversationId/messages/:messageId/reaction", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId/reaction", rt.wrap(rt.uncommentMessage))

	// --- DIRECT CHAT ---
	rt.router.POST("/direct-conversations", rt.wrap(rt.createDirectConversation))

	// --- GROUPS ---
	rt.router.POST("/groups", rt.wrap(rt.createGroup))
	rt.router.POST("/groups/:groupId", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/groups/:groupId", rt.wrap(rt.leaveGroup))
	rt.router.PUT("/groups/:groupId/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/groups/:groupId/photo", rt.wrap(rt.setGroupPhoto))

	// --- HEALTH CHECK ---
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
