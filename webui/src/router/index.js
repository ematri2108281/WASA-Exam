import { createRouter, createWebHashHistory } from 'vue-router'

// Lazy-load views for better performance
const HomeView        = () => import('../views/HomeView.vue')
const LoginView       = () => import('../views/LoginPage/LoginView.vue')
const ProfileView     = () => import('../views/ProfileView.vue')
const SearchView      = () => import('../views/SearchView.vue')
const ConvView        = () => import('../views/ConvView.vue')
const GroupCreateView = () => import('../views/GroupCreateView.vue')
const GroupEditView   = () => import('../views/GroupEditView.vue')
const UsersView       = () => import('../views/UsersView.vue')

const routes = [
	{ path: '/',                             redirect: { name: 'login' } },
	{ path: '/login',                        name: 'login',        component: LoginView },
	{ path: '/home',                         name: 'home',         component: HomeView,        meta: { requiresAuth: true } },
	{ path: '/profile',                      name: 'profile',      component: ProfileView,     meta: { requiresAuth: true } },
	{ path: '/users',                        name: 'users',        component: UsersView,       meta: { requiresAuth: true } },
	{ path: '/search',                       name: 'search',       component: SearchView,      meta: { requiresAuth: true } },
	{ path: '/conversations/:conversationId',name: 'conversation', component: ConvView,        meta: { requiresAuth: true } },
	{ path: '/groups/create',               name: 'group-create', component: GroupCreateView, meta: { requiresAuth: true } },
	{ path: '/groups/:groupId/edit',        name: 'group-edit',   component: GroupEditView,   meta: { requiresAuth: true } },
]

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes,
})

function isLoggedIn() {
	return Boolean(localStorage.getItem('token'))
}

router.beforeEach((to, _from, next) => {
	const authenticated = isLoggedIn()

	if (to.meta.requiresAuth && !authenticated) {
		next({ name: 'login' })
	} else if (to.name === 'login' && authenticated) {
		next({ name: 'home' })
	} else {
		next()
	}
})

export default router
