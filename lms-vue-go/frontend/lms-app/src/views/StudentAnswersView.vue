<template>
  <div class="student-answers-view">
    <h1>{{ isAdmin ? 'Manajemen Jawaban Siswa' : 'Jawaban Saya' }}</h1>

    <div v-if="loading" class="loading">
      <p>Memuat data...</p>
    </div>

    <div v-else-if="error" class="error-message">
      <p>{{ error }}</p>
    </div>

    <div v-else>
      <div v-if="answers.length === 0" class="no-data">
        <p>{{ isAdmin ? 'Belum ada jawaban siswa.' : 'Belum ada jawaban yang dikirimkan.' }}</p>
      </div>

      <!-- Admin View -->
      <div v-else-if="isAdmin" class="admin-answers-list">
        <div class="filters">
          <div class="filter-row">
            <input
              type="text"
              v-model="searchQuery"
              placeholder="Cari berdasarkan nama siswa..."
              class="search-input"
            />
            <select v-model="selectedClass" class="class-filter">
              <option value="">Semua Kelas</option>
              <option v-for="classItem in uniqueClasses" :key="classItem" :value="classItem">
                {{ classItem }}
              </option>
            </select>
            <div class="export-buttons">
              <button
                class="export-btn"
                @click="exportToExcel"
                :disabled="filteredStudents.length === 0"
                title="Export daftar siswa"
              >
                <i class="fas fa-file-excel"></i> Export Daftar Siswa
              </button>
              <button
                class="export-btn export-btn-all"
                @click="exportAllAnswersToExcel"
                :disabled="answers.length === 0"
                title="Export semua jawaban"
              >
                <i class="fas fa-file-excel"></i> Export Semua Jawaban
              </button>
            </div>
          </div>
        </div>

        <!-- Student List View -->
        <div v-if="!selectedStudent" class="student-list">
          <h2>Daftar Siswa</h2>
          <table class="students-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Nama Siswa</th>
                <th>Kelas</th>
                <th>Jumlah Jawaban</th>
                <th>Rata-rata Nilai</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="student in filteredStudents" :key="student.id">
                <td>{{ student.id }}</td>
                <td>{{ student.name }}</td>
                <td>{{ student.class }}</td>
                <td>{{ student.answerCount }}</td>
                <td>{{ student.averageScore !== null ? student.averageScore.toFixed(1) : 'Belum dinilai' }}</td>
                <td>
                  <button class="view-btn" @click="viewStudentAnswers(student)">
                    Lihat Jawaban
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Student Answers Detail View -->
        <div v-else class="student-answers-detail">
          <div class="detail-header">
            <button class="back-btn" @click="backToStudentList">
              &larr; Kembali ke Daftar Siswa
            </button>
            <h2>Jawaban {{ selectedStudent.name }} ({{ selectedStudent.class }})</h2>
            <button
              class="export-btn"
              @click="exportStudentAnswersToExcel"
              :disabled="selectedStudentAnswers.length === 0"
            >
              <i class="fas fa-file-excel"></i> Export Jawaban
            </button>
          </div>

          <table class="answers-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Soal</th>
                <th>Jawaban</th>
                <th>Nilai</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="answer in selectedStudentAnswers" :key="answer.id">
                <td>{{ answer.id }}</td>
                <td>{{ answer.question_text }}</td>
                <td>{{ answer.answer }}</td>
                <td>
                  <span v-if="answer.score !== null">{{ answer.score }} / {{ answer.question_score }}</span>
                  <span v-else class="pending">Belum dinilai</span>
                </td>
                <td>
                  <button
                    v-if="answer.question_type === 'essay'"
                    class="grade-btn"
                    @click="showGradeModal(answer)"
                  >
                    {{ answer.score !== null ? 'Ubah Nilai' : 'Beri Nilai' }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Student View -->
      <div v-else class="answers-list">
        <div v-for="answer in answers" :key="answer.id" class="answer-card">
          <div class="answer-header">
            <h3>{{ answer.question_text }}</h3>
            <div class="answer-meta">
              <span class="question-type">{{ formatQuestionType(answer.question_type) }}</span>
              <span class="question-score">Nilai Maksimal: {{ answer.question_score }}</span>
            </div>
          </div>

          <div class="answer-content">
            <p><strong>Jawaban Anda:</strong> {{ answer.answer }}</p>
          </div>

          <div class="answer-footer">
            <div class="score" v-if="answer.score !== null">
              <span>Nilai: {{ answer.score }} / {{ answer.question_score }}</span>
            </div>
            <div class="score pending" v-else>
              <span>Belum dinilai</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Grade Modal -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal-content">
        <h3>Beri Nilai</h3>
        <p><strong>Siswa:</strong> {{ selectedAnswer.student_name }}</p>
        <p><strong>Kelas:</strong> {{ selectedAnswer.student_class }}</p>
        <p><strong>Soal:</strong> {{ selectedAnswer.question_text }}</p>
        <p><strong>Jawaban:</strong> {{ selectedAnswer.answer }}</p>

        <div class="form-group">
          <label for="score">Nilai (Maksimal {{ selectedAnswer.question_score }}):</label>
          <input
            type="number"
            id="score"
            v-model.number="gradeScore"
            :max="selectedAnswer.question_score"
            min="0"
            required
          />
        </div>

        <div class="modal-actions">
          <button class="cancel-btn" @click="closeModal">Batal</button>
          <button
            class="save-btn"
            @click="submitGrade"
            :disabled="gradeLoading"
          >
            {{ gradeLoading ? 'Menyimpan...' : 'Simpan Nilai' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { answerService } from '@/services/api';
import * as XLSX from 'xlsx';

export default {
  name: 'StudentAnswersView',
  data() {
    return {
      answers: [],
      loading: true,
      error: null,
      searchQuery: '',
      selectedClass: '',
      showModal: false,
      selectedAnswer: null,
      selectedStudent: null,
      gradeScore: 0,
      gradeLoading: false
    };
  },
  computed: {
    // Check if user is admin
    isAdmin() {
      const user = JSON.parse(localStorage.getItem('user') || '{}');
      return user && user.role === 'admin';
    },

    // Get unique classes from answers
    uniqueClasses() {
      const classes = this.answers.map(answer => answer.student_class);
      return [...new Set(classes)].sort();
    },

    // Group answers by student
    studentsList() {
      const studentsMap = new Map();

      this.answers.forEach(answer => {
        const studentId = answer.student_id;
        if (!studentsMap.has(studentId)) {
          studentsMap.set(studentId, {
            id: studentId,
            name: answer.student_name,
            class: answer.student_class,
            answers: [],
            answerCount: 0,
            totalScore: 0,
            scoredAnswersCount: 0
          });
        }

        const student = studentsMap.get(studentId);
        student.answers.push(answer);
        student.answerCount++;

        if (answer.score !== null) {
          student.totalScore += answer.score;
          student.scoredAnswersCount++;
        }
      });

      // Calculate average score
      const studentsList = Array.from(studentsMap.values());
      studentsList.forEach(student => {
        student.averageScore = student.scoredAnswersCount > 0
          ? student.totalScore / student.scoredAnswersCount
          : null;
      });

      return studentsList;
    },

    // Filter students based on search query and selected class
    filteredStudents() {
      return this.studentsList.filter(student => {
        const matchesSearch = !this.searchQuery ||
          student.name.toLowerCase().includes(this.searchQuery.toLowerCase());

        const matchesClass = !this.selectedClass ||
          student.class === this.selectedClass;

        return matchesSearch && matchesClass;
      });
    },

    // Get answers for selected student
    selectedStudentAnswers() {
      if (!this.selectedStudent) return [];

      return this.answers.filter(answer =>
        answer.student_id === this.selectedStudent.id
      );
    }
  },
  created() {
    this.fetchAnswers();
  },
  methods: {
    async fetchAnswers() {
      this.loading = true;
      this.error = null;

      try {
        let response;
        if (this.isAdmin) {
          // Admin gets all answers
          response = await answerService.getAllAnswers();
        } else {
          // Student gets only their answers
          response = await answerService.getMyAnswers();
        }

        this.answers = response.data.data;
      } catch (error) {
        console.error('Error fetching answers:', error);
        this.error = 'Gagal mengambil data jawaban. Silakan coba lagi nanti.';
      } finally {
        this.loading = false;
      }
    },

    formatQuestionType(type) {
      switch (type) {
        case 'multiple_choice':
          return 'Pilihan Ganda';
        case 'essay':
          return 'Essay';
        default:
          return type;
      }
    },

    // View a student's answers
    viewStudentAnswers(student) {
      this.selectedStudent = student;
    },

    // Go back to student list
    backToStudentList() {
      this.selectedStudent = null;
    },

    // Show modal for grading
    showGradeModal(answer) {
      this.selectedAnswer = answer;
      this.gradeScore = answer.score !== null ? answer.score : 0;
      this.showModal = true;
    },

    // Close modal
    closeModal() {
      this.showModal = false;
      this.selectedAnswer = null;
      this.gradeScore = 0;
    },

    // Submit grade
    async submitGrade() {
      if (this.gradeScore < 0 || this.gradeScore > this.selectedAnswer.question_score) {
        alert(`Nilai harus antara 0 dan ${this.selectedAnswer.question_score}`);
        return;
      }

      this.gradeLoading = true;

      try {
        await answerService.gradeAnswer(this.selectedAnswer.id, this.gradeScore);

        // Update the answer in the list
        const index = this.answers.findIndex(a => a.id === this.selectedAnswer.id);
        if (index !== -1) {
          this.answers[index].score = this.gradeScore;
        }

        // Recalculate student average
        if (this.selectedStudent) {
          const student = this.studentsList.find(s => s.id === this.selectedStudent.id);
          if (student) {
            // Recalculate student stats
            student.totalScore = 0;
            student.scoredAnswersCount = 0;

            student.answers.forEach(answer => {
              if (answer.score !== null) {
                student.totalScore += answer.score;
                student.scoredAnswersCount++;
              }
            });

            student.averageScore = student.scoredAnswersCount > 0
              ? student.totalScore / student.scoredAnswersCount
              : null;
          }
        }

        this.closeModal();
      } catch (error) {
        console.error('Error grading answer:', error);
        alert('Gagal menyimpan nilai. Silakan coba lagi.');
      } finally {
        this.gradeLoading = false;
      }
    },

    // Export data to Excel
    exportToExcel() {
      // Prepare data for export
      const data = this.filteredStudents.map(student => {
        return {
          'ID': student.id,
          'Nama Siswa': student.name,
          'Kelas': student.class,
          'Jumlah Jawaban': student.answerCount,
          'Rata-rata Nilai': student.averageScore !== null ? student.averageScore.toFixed(1) : 'Belum dinilai'
        };
      });

      // Create worksheet
      const worksheet = XLSX.utils.json_to_sheet(data);

      // Create workbook
      const workbook = XLSX.utils.book_new();

      // Add worksheet to workbook
      const title = this.selectedClass ? `Siswa Kelas ${this.selectedClass}` : 'Semua Siswa';
      XLSX.utils.book_append_sheet(workbook, worksheet, title);

      // Generate Excel file
      const excelBuffer = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });

      // Save file
      this.saveExcelFile(excelBuffer, title);
    },

    // Export student answers to Excel
    exportStudentAnswersToExcel() {
      if (!this.selectedStudent) return;

      // Prepare data for export
      const data = this.selectedStudentAnswers.map(answer => {
        return {
          'ID': answer.id,
          'Soal': answer.question_text,
          'Jawaban': answer.answer,
          'Jenis Soal': this.formatQuestionType(answer.question_type),
          'Nilai': answer.score !== null ? `${answer.score} / ${answer.question_score}` : 'Belum dinilai'
        };
      });

      // Create worksheet
      const worksheet = XLSX.utils.json_to_sheet(data);

      // Create workbook
      const workbook = XLSX.utils.book_new();

      // Add worksheet to workbook
      const title = `Jawaban ${this.selectedStudent.name}`;
      XLSX.utils.book_append_sheet(workbook, worksheet, title);

      // Generate Excel file
      const excelBuffer = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });

      // Save file
      this.saveExcelFile(excelBuffer, title);
    },

    // Export all answers to Excel
    exportAllAnswersToExcel() {
      // Create workbook
      const workbook = XLSX.utils.book_new();

      // Group answers by student
      const studentGroups = {};

      this.answers.forEach(answer => {
        const studentId = answer.student_id;
        const studentName = answer.student_name;

        if (!studentGroups[studentId]) {
          studentGroups[studentId] = {
            name: studentName,
            class: answer.student_class,
            answers: []
          };
        }

        studentGroups[studentId].answers.push(answer);
      });

      // Create a worksheet for each student
      Object.keys(studentGroups).forEach(studentId => {
        const student = studentGroups[studentId];

        // Prepare data for export
        const data = student.answers.map(answer => {
          return {
            'ID': answer.id,
            'Soal': answer.question_text,
            'Jawaban': answer.answer,
            'Jenis Soal': this.formatQuestionType(answer.question_type),
            'Nilai': answer.score !== null ? `${answer.score} / ${answer.question_score}` : 'Belum dinilai'
          };
        });

        // Create worksheet
        const worksheet = XLSX.utils.json_to_sheet(data);

        // Add worksheet to workbook
        const sheetName = `${student.name} (${student.class})`;
        XLSX.utils.book_append_sheet(workbook, worksheet, sheetName);
      });

      // Create a summary worksheet
      const summaryData = this.studentsList.map(student => {
        return {
          'ID': student.id,
          'Nama Siswa': student.name,
          'Kelas': student.class,
          'Jumlah Jawaban': student.answerCount,
          'Rata-rata Nilai': student.averageScore !== null ? student.averageScore.toFixed(1) : 'Belum dinilai'
        };
      });

      const summaryWorksheet = XLSX.utils.json_to_sheet(summaryData);
      XLSX.utils.book_append_sheet(workbook, summaryWorksheet, 'Ringkasan');

      // Generate Excel file
      const excelBuffer = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });

      // Save file
      const title = this.selectedClass ? `Semua Jawaban Kelas ${this.selectedClass}` : 'Semua Jawaban';
      this.saveExcelFile(excelBuffer, title);
    },

    // Save Excel file
    saveExcelFile(buffer, fileName) {
      const data = new Blob([buffer], { type: 'application/octet-stream' });
      const link = document.createElement('a');

      // Create a URL for the blob
      const url = window.URL.createObjectURL(data);
      link.href = url;

      // Set file name
      link.download = `${fileName}.xlsx`;

      // Trigger download
      link.click();

      // Clean up
      setTimeout(() => {
        window.URL.revokeObjectURL(url);
        link.remove();
      }, 100);
    }
  }
};
</script>

<style scoped>
.student-answers-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.loading, .error-message, .no-data {
  text-align: center;
  margin: 40px 0;
}

.error-message {
  color: #e53935;
}

/* Student View Styles */
.answers-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.answer-card {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  background-color: #fff;
}

.answer-header {
  margin-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 12px;
}

.answer-header h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
}

.answer-meta {
  display: flex;
  gap: 12px;
  font-size: 14px;
  color: #666;
}

.question-type {
  background-color: #e3f2fd;
  padding: 4px 8px;
  border-radius: 4px;
}

.question-score {
  background-color: #f1f8e9;
  padding: 4px 8px;
  border-radius: 4px;
}

.answer-content {
  margin-bottom: 16px;
}

.answer-footer {
  display: flex;
  justify-content: flex-end;
}

.score {
  background-color: #e8f5e9;
  color: #2e7d32;
  padding: 6px 12px;
  border-radius: 4px;
  font-weight: 500;
}

.score.pending, .pending {
  background-color: #fff8e1;
  color: #ff8f00;
}

/* Admin View Styles */
.admin-answers-list {
  width: 100%;
}

.filters {
  margin-bottom: 20px;
}

.filter-row {
  display: flex;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;
}

.export-buttons {
  display: flex;
  gap: 8px;
  margin-left: auto;
}

.search-input {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  flex: 1;
  max-width: 400px;
}

.class-filter {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  min-width: 150px;
}

/* Student List Styles */
.student-list {
  margin-top: 20px;
}

.student-list h2 {
  margin-bottom: 16px;
  font-size: 20px;
  color: #333;
}

.students-table {
  width: 100%;
  border-collapse: collapse;
  background-color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  overflow: hidden;
}

.students-table th, .students-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.students-table th {
  background-color: #f5f5f5;
  font-weight: 600;
}

.students-table tr:hover {
  background-color: #f9f9f9;
}

.view-btn {
  background-color: #2196f3;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 14px;
}

.view-btn:hover {
  background-color: #1976d2;
}

.export-btn {
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.export-btn:hover {
  background-color: #388e3c;
}

.export-btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}

.export-btn-all {
  background-color: #2196f3;
}

.export-btn-all:hover {
  background-color: #1976d2;
}

.export-btn-all:disabled {
  background-color: #90caf9;
  cursor: not-allowed;
}

/* Student Answers Detail Styles */
.student-answers-detail {
  margin-top: 20px;
}

.detail-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;
}

.detail-header h2 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.back-btn {
  background-color: #f5f5f5;
  color: #333;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
}

.back-btn:hover {
  background-color: #e0e0e0;
}

.answers-table {
  width: 100%;
  border-collapse: collapse;
  background-color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  overflow: hidden;
}

.answers-table th, .answers-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.answers-table th {
  background-color: #f5f5f5;
  font-weight: 600;
}

.answers-table tr:hover {
  background-color: #f9f9f9;
}

.grade-btn {
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 14px;
}

.grade-btn:hover {
  background-color: #388e3c;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  padding: 24px;
  width: 500px;
  max-width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 16px;
}

.form-group {
  margin-bottom: 20px;
  margin-top: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.cancel-btn, .save-btn {
  padding: 10px 16px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  border: none;
}

.cancel-btn {
  background-color: #f5f5f5;
  color: #333;
}

.save-btn {
  background-color: #4caf50;
  color: white;
}

.save-btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}
</style>
