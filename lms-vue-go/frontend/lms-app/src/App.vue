<template>
  <div id="app">
    <nav class="navbar">
      <div class="navbar-brand">
        <router-link to="/" class="brand-logo">LMS Sekolah</router-link>
      </div>
      <div class="navbar-menu">
        <router-link to="/" class="nav-link">Beranda</router-link>
        <router-link to="/students" class="nav-link">Siswa</router-link>
        <router-link to="/questions" class="nav-link">Soal</router-link>
      </div>
    </nav>

    <router-view/>

    <footer class="footer">
      <p>&copy; 2025 LMS Sekolah - Sistem Pembelajaran Online</p>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'App'
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

.footer {
  background-color: #333;
  color: #fff;
  text-align: center;
  padding: 20px;
  margin-top: auto;
}
</style>
