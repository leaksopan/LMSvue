<template>
  <div class="home">
    <div class="hero">
      <h1>Selamat Datang di LMS Sekolah</h1>
      <p>Sistem Pembelajaran Online untuk Siswa dan Guru</p>

      <div class="status-card" v-if="serverStatus">
        <p><strong>Status Server:</strong> {{ serverStatus.status }}</p>
        <p>
          <small>{{ serverStatus.message }}</small>
        </p>
        <p>
          <small>{{ serverStatus.time }}</small>
        </p>
      </div>
    </div>

    <div class="feature-section">
      <div class="feature-card">
        <h2>Siswa</h2>
        <p>Lihat daftar siswa dan data mereka</p>
        <router-link to="/students" class="btn">Lihat Siswa</router-link>
      </div>

      <div class="feature-card">
        <h2>Soal</h2>
        <p>Akses soal-soal dan jawab kuis</p>
        <router-link to="/questions" class="btn">Lihat Soal</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { serverStatus } from "@/services/api";

export default {
  name: "HomeView",
  data() {
    return {
      serverStatus: null,
    };
  },
  created() {
    this.checkServerStatus();
  },
  methods: {
    async checkServerStatus() {
      try {
        const response = await serverStatus.check();
        this.serverStatus = response.data;
      } catch (error) {
        console.error("Error connecting to server", error);
      }
    },
  },
};
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.hero {
  text-align: center;
  margin-bottom: 40px;
  padding: 40px 20px;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.hero h1 {
  font-size: 2.5em;
  margin-bottom: 10px;
  color: #2c3e50;
}

.hero p {
  font-size: 1.2em;
  color: #666;
}

.status-card {
  margin-top: 20px;
  padding: 15px;
  background-color: #e8f5e9;
  border-radius: 8px;
  display: inline-block;
}

.feature-section {
  display: flex;
  gap: 20px;
  justify-content: center;
  flex-wrap: wrap;
}

.feature-card {
  flex: 1;
  min-width: 300px;
  padding: 30px;
  border-radius: 8px;
  background-color: #fff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.feature-card h2 {
  color: #42b983;
  margin-bottom: 15px;
}

.btn {
  display: inline-block;
  margin-top: 15px;
  padding: 10px 20px;
  background-color: #42b983;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  font-weight: bold;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: #3aa876;
}
</style>
