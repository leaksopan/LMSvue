<template>
  <div class="student-detail">
    <div v-if="loading" class="loading">
      <p>Memuat data siswa...</p>
    </div>

    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <button @click="fetchStudent" class="btn">Coba Lagi</button>
      <router-link to="/students" class="btn btn-secondary"
        >Kembali</router-link
      >
    </div>

    <div v-else-if="student" class="student-card">
      <div class="actions">
        <router-link to="/students" class="btn btn-secondary"
          >← Kembali</router-link
        >
      </div>

      <h1>{{ student.name }}</h1>

      <div class="student-info">
        <div class="info-item">
          <strong>ID Siswa:</strong>
          <span>{{ student.id }}</span>
        </div>

        <div class="info-item">
          <strong>Kelas:</strong>
          <span>{{ student.class }}</span>
        </div>

        <div v-if="student.email" class="info-item">
          <strong>Email:</strong>
          <span>{{ student.email }}</span>
        </div>
      </div>

      <!-- Placeholder untuk data tambahan -->
      <div class="additional-info">
        <h2>Informasi Tambahan</h2>
        <p class="placeholder-text">
          Data nilai dan aktivitas siswa akan ditampilkan di sini.
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { studentService } from "@/services/api";

export default {
  name: "StudentDetailView",
  data() {
    return {
      student: null,
      loading: true,
      error: null,
    };
  },
  created() {
    this.fetchStudent();
  },
  methods: {
    async fetchStudent() {
      const studentId = this.$route.params.id;

      this.loading = true;
      this.error = null;

      try {
        const response = await studentService.getById(studentId);
        this.student = response.data.data;
      } catch (error) {
        console.error("Error fetching student details", error);
        this.error = "Gagal memuat detail siswa. Silakan coba lagi nanti.";
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.student-detail {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.loading,
.error {
  text-align: center;
  padding: 40px;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.error {
  color: #e53935;
}

.student-card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 30px;
  position: relative;
}

.actions {
  margin-bottom: 20px;
}

h1 {
  color: #42b983;
  margin-bottom: 20px;
  text-align: center;
}

h2 {
  color: #2c3e50;
  margin-top: 30px;
  margin-bottom: 15px;
}

.student-info {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 20px;
}

.info-item {
  display: flex;
  margin-bottom: 12px;
}

.info-item strong {
  width: 120px;
  color: #555;
}

.additional-info {
  margin-top: 30px;
}

.placeholder-text {
  color: #888;
  font-style: italic;
}

.btn {
  display: inline-block;
  padding: 8px 16px;
  background-color: #42b983;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  font-weight: bold;
  border: none;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: #3aa876;
}

.btn-secondary {
  background-color: #607d8b;
}

.btn-secondary:hover {
  background-color: #546e7a;
}
</style>
