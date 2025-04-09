import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/students',
    name: 'students',
    component: () => import('../views/StudentsView.vue')
  },
  {
    path: '/students/:id',
    name: 'student-detail',
    component: () => import('../views/StudentDetailView.vue')
  },
  {
    path: '/questions',
    name: 'questions',
    component: () => import('../views/QuestionsView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router 