import { createRouter, createWebHistory } from 'vue-router'

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
    component: () => import('@/views/ProfileView.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
