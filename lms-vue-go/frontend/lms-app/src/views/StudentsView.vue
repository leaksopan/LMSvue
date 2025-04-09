<template>
  <div class="students">
    <h1>Daftar Siswa</h1>

    <div v-if="loading" class="loading">
      <p>Memuat data siswa...</p>
    </div>

    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <button @click="fetchStudents" class="btn">Coba Lagi</button>
    </div>

    <div v-else>
      <div class="students-grid">
        <div v-for="student in students" :key="student.id" class="student-card">
          <h2>{{ student.name }}</h2>
          <p><strong>Kelas:</strong> {{ student.class }}</p>
          <p v-if="student.email">
            <strong>Email:</strong> {{ student.email }}
          </p>
          <router-link :to="`/students/${student.id}`" class="btn"
            >Lihat Detail</router-link
          >
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { studentService } from "@/services/api";

export default {
  name: "StudentsView",
  data() {
    return {
      students: [],
      loading: true,
      error: null,
    };
  },
  created() {
    this.fetchStudents();
  },
  methods: {
    async fetchStudents() {
      this.loading = true;
      this.error = null;

      try {
        const response = await studentService.getAll();
        this.students = response.data.data;
      } catch (error) {
        console.error("Error fetching students", error);
        this.error = "Gagal memuat data siswa. Silakan coba lagi nanti.";
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.students {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

h1 {
  margin-bottom: 30px;
  color: #2c3e50;
  text-align: center;
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

.students-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.student-card {
  padding: 20px;
  border-radius: 8px;
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.student-card h2 {
  color: #42b983;
  margin-bottom: 10px;
}

.btn {
  display: inline-block;
  margin-top: 15px;
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
</style>
