<template>
  <div id="app">
    <nav class="navbar">
      <div class="navbar-brand">
        <router-link to="/" class="brand-logo">LMS Sekolah</router-link>
      </div>
      <div class="navbar-menu">
        <template v-if="isLoggedIn">
          <router-link to="/" class="nav-link">Beranda</router-link>
          <router-link to="/students" class="nav-link">Siswa</router-link>
          <router-link to="/questions" class="nav-link">Soal</router-link>
          <router-link to="/my-answers" class="nav-link">Jawaban Saya</router-link>
          <router-link v-if="isAdmin" to="/admin/questions" class="nav-link admin-link">Manajemen Soal</router-link>
        </template>
      </div>
      <div class="navbar-auth">
        <router-link to="/test-api" class="nav-link test-btn">Test API</router-link>
        <template v-if="isLoggedIn">
          <span class="user-info">{{ currentUser ? currentUser.username : '' }}</span>
          <button @click="logout" class="logout-btn">Logout</button>
        </template>
        <template v-else>
          <router-link to="/login" class="nav-link login-btn">Login</router-link>
        </template>
      </div>
    </nav>

    <router-view/>

    <footer class="footer">
      <p>&copy; 2025 LMS Sekolah - Sistem Pembelajaran Online</p>
    </footer>
  </div>
</template>

<script>
import { authService } from '@/services/api';

export default {
  name: 'App',
  data() {
    return {
      isLoggedIn: false,
      currentUser: null
    };
  },
  created() {
    // Check login status
    this.checkAuth();

    // Listen for route changes to update auth status
    this.$router.beforeEach((to, from, next) => {
      this.checkAuth();
      next();
    });
  },
  computed: {
    isAdmin() {
      return this.currentUser && this.currentUser.role === 'admin';
    }
  },
  methods: {
    checkAuth() {
      this.isLoggedIn = authService.isLoggedIn();
      this.currentUser = authService.getUser();
    },
    logout() {
      authService.logout();
      this.isLoggedIn = false;
      this.currentUser = null;
      this.$router.push('/login');
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f5f5;
  color: #333;
  line-height: 1.6;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background-color: #42b983;
  color: white;
  padding: 15px 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.navbar-brand {
  font-size: 1.5rem;
  font-weight: bold;
}

.brand-logo {
  color: white;
  text-decoration: none;
}

.navbar-menu {
  display: flex;
  gap: 20px;
  flex-grow: 1;
}

.navbar-auth {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-info {
  color: white;
  font-weight: 500;
}

.logout-btn {
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.logout-btn:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

.login-btn {
  background-color: white;
  color: #42b983 !important;
  padding: 8px 16px;
  border-radius: 4px;
  font-weight: bold;
}

.test-btn {
  background-color: #f0ad4e;
  color: white !important;
  padding: 8px 16px;
  border-radius: 4px;
  font-weight: bold;
  margin-right: 10px;
}

.nav-link {
  color: white;
  text-decoration: none;
  padding: 8px 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-link:hover {
  background-color: rgba(255, 255, 255, 0.15);
}

.router-link-active {
  background-color: rgba(255, 255, 255, 0.2);
  font-weight: bold;
}

.admin-link {
  background-color: #ff9800;
  color: white !important;
  margin-left: 10px;
  font-weight: 500;
}

.footer {
  background-color: #333;
  color: #fff;
  text-align: center;
  padding: 20px;
  margin-top: auto;
}
</style>
