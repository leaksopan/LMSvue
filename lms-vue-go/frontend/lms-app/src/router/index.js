import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import { authService } from "../services/api";

const routes = [
  {
    path: "/",
    name: "home",
    component: HomeView,
    meta: { requiresAuth: true },
  },
  {
    path: "/login",
    name: "login",
    component: () => import("../views/LoginView.vue"),
    meta: { guest: true },
  },
  {
    path: "/test-api",
    name: "test-api",
    component: () => import("../views/TestApiView.vue"),
  },
  {
    path: "/students",
    name: "students",
    component: () => import("../views/StudentsView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/students/:id",
    name: "student-detail",
    component: () => import("../views/StudentDetailView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/questions",
    name: "questions",
    component: () => import("../views/QuestionsView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/my-answers",
    name: "my-answers",
    component: () => import("../views/StudentAnswersView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/admin/questions",
    name: "admin-questions",
    component: () => import("../views/AdminQuestionsView.vue"),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: "/:pathMatch(.*)*",
    name: "not-found",
    component: () => import("../views/MaintenanceView.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

// Navigation guards
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const requiresAdmin = to.matched.some((record) => record.meta.requiresAdmin);
  const isGuestOnly = to.matched.some((record) => record.meta.guest);
  const isLoggedIn = authService.isLoggedIn();
  const user = authService.getUser();
  const isAdmin = user && user.role === "admin";

  // Always allow access to test-api page
  if (to.name === "test-api") {
    next();
    return;
  }

  if (requiresAuth && !isLoggedIn) {
    // Jika halaman memerlukan auth tapi user belum login
    next({ name: "login" });
  } else if (requiresAdmin && !isAdmin) {
    // Jika halaman memerlukan role admin tapi user bukan admin
    next({ name: "home" });
  } else if (isGuestOnly && isLoggedIn) {
    // Jika halaman hanya untuk guest tapi user sudah login
    next({ name: "home" });
  } else {
    // Lanjutkan navigasi
    next();
  }
});

export default router;
