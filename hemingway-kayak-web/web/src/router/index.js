import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue')
  },
  {
    path: '/tracking',
    name: 'Tracking',
    component: () => import('@/views/TrackingView.vue')
  },
  {
    path: '/tournaments',
    name: 'Tournaments',
    component: () => import('@/views/TournamentsView.vue')
  },
  {
    path: '/tournaments/:id',
    name: 'TournamentDetail',
    component: () => import('@/views/TournamentDetailView.vue')
  },
  {
    path: '/history',
    name: 'History',
    component: () => import('@/views/HistoryView.vue')
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile/edit',
    name: 'ProfileEdit',
    component: () => import('@/views/ProfileEditView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/RegisterView.vue')
  },
  {
    path: '/tournaments/create',
    name: 'TournamentCreate',
    component: () => import('@/views/TournamentCreateView.vue')
  },
  {
    path: '/descargar',
    name: 'Download',
    component: () => import('@/views/DownloadView.vue')
  },
  {
    path: '/trips/:id',
    name: 'TripDetail',
    component: () => import('@/views/TripDetailView.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFoundView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  // Hay token guardado pero el usuario aún no se cargó (recarga de página):
  // lo resuelve antes de decidir, si no la sesión no sobrevivía al F5.
  if (auth.token && !auth.user) {
    await auth.fetchUser()
  }
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.path === '/login' && auth.isAuthenticated) {
    return { path: '/' }
  }
  return true
})

export default router
