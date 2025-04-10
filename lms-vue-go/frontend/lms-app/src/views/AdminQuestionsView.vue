<template>
  <div class="admin-questions-view">
    <div class="header">
      <h1>Manajemen Soal</h1>
      <button class="add-btn" @click="showAddForm = true" v-if="!showAddForm && !showEditForm">
        + Tambah Soal Baru
      </button>
    </div>
    
    <div v-if="loading" class="loading">
      <p>Memuat data...</p>
    </div>
    
    <div v-else-if="error" class="error-message">
      <p>{{ error }}</p>
      <button @click="fetchQuestions" class="retry-btn">Coba Lagi</button>
    </div>
    
    <div v-else>
      <QuestionForm 
        v-if="showAddForm" 
        @saved="onQuestionSaved" 
        @cancel="showAddForm = false"
      />
      
      <QuestionForm 
        v-else-if="showEditForm" 
        :question="selectedQuestion" 
        :isEdit="true"
        @saved="onQuestionSaved" 
        @cancel="showEditForm = false"
      />
      
      <div v-else>
        <div v-if="questions.length === 0" class="no-data">
          <p>Belum ada soal. Klik tombol "Tambah Soal Baru" untuk membuat soal pertama.</p>
        </div>
        
        <div v-else class="questions-list">
          <div v-for="question in questions" :key="question.id" class="question-card">
            <div class="question-header">
              <div class="question-type">
                {{ formatQuestionType(question.type) }}
                <span class="question-score">({{ question.score }} poin)</span>
              </div>
              <div class="question-actions">
                <button class="edit-btn" @click="editQuestion(question)">Edit</button>
                <button class="delete-btn" @click="confirmDelete(question)">Hapus</button>
              </div>
            </div>
            
            <h3>{{ question.question }}</h3>
            
            <div v-if="question.type === 'multiple_choice'" class="options-list">
              <div 
                v-for="(option, index) in question.options" 
                :key="index"
                class="option-item"
                :class="{ 'correct-answer': String.fromCharCode(65 + index) === question.answer }"
              >
                <span class="option-letter">{{ String.fromCharCode(65 + index) }}.</span>
                <span class="option-text">{{ option }}</span>
              </div>
            </div>
            
            <div v-if="question.image_url" class="question-image">
              <img :src="question.image_url" alt="Question image" />
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal-overlay">
      <div class="modal-content">
        <h3>Konfirmasi Hapus</h3>
        <p>Apakah Anda yakin ingin menghapus soal ini?</p>
        <p class="warning">Tindakan ini tidak dapat dibatalkan.</p>
        
        <div class="modal-actions">
          <button class="cancel-btn" @click="showDeleteModal = false">Batal</button>
          <button 
            class="delete-confirm-btn" 
            @click="deleteQuestion" 
            :disabled="deleteLoading"
          >
            {{ deleteLoading ? 'Menghapus...' : 'Hapus' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { questionService } from '@/services/api';
import QuestionForm from '@/components/QuestionForm.vue';

export default {
  name: 'AdminQuestionsView',
  components: {
    QuestionForm
  },
  data() {
    return {
      questions: [],
      loading: true,
      error: null,
      showAddForm: false,
      showEditForm: false,
      selectedQuestion: null,
      showDeleteModal: false,
      deleteLoading: false
    };
  },
  created() {
    this.fetchQuestions();
  },
  methods: {
    async fetchQuestions() {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await questionService.getAll();
        this.questions = response.data.data;
      } catch (error) {
        console.error('Error fetching questions:', error);
        this.error = 'Gagal mengambil data soal. Silakan coba lagi nanti.';
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
    
    editQuestion(question) {
      this.selectedQuestion = question;
      this.showEditForm = true;
    },
    
    confirmDelete(question) {
      this.selectedQuestion = question;
      this.showDeleteModal = true;
    },
    
    async deleteQuestion() {
      this.deleteLoading = true;
      
      try {
        await questionService.delete(this.selectedQuestion.id);
        this.questions = this.questions.filter(q => q.id !== this.selectedQuestion.id);
        this.showDeleteModal = false;
        this.selectedQuestion = null;
      } catch (error) {
        console.error('Error deleting question:', error);
        alert('Gagal menghapus soal. Silakan coba lagi.');
      } finally {
        this.deleteLoading = false;
      }
    },
    
    onQuestionSaved(question) {
      if (this.showEditForm) {
        // Update existing question in the list
        const index = this.questions.findIndex(q => q.id === question.id);
        if (index !== -1) {
          this.questions.splice(index, 1, question);
        }
        this.showEditForm = false;
      } else {
        // Add new question to the list
        this.questions.unshift(question);
        this.showAddForm = false;
      }
      
      this.selectedQuestion = null;
    }
  }
};
</script>

<style scoped>
.admin-questions-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.add-btn {
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 16px;
  font-size: 16px;
  cursor: pointer;
}

.loading, .error-message, .no-data {
  text-align: center;
  margin: 40px 0;
}

.error-message {
  color: #e53935;
}

.retry-btn {
  background-color: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 8px 16px;
  margin-top: 16px;
  cursor: pointer;
}

.questions-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.question-card {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  background-color: #fff;
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.question-type {
  background-color: #e3f2fd;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 14px;
  display: flex;
  align-items: center;
}

.question-score {
  margin-left: 8px;
  font-size: 14px;
  color: #666;
}

.question-actions {
  display: flex;
  gap: 8px;
}

.edit-btn, .delete-btn {
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  border: none;
}

.edit-btn {
  background-color: #f5f5f5;
  color: #333;
}

.delete-btn {
  background-color: #ffebee;
  color: #e53935;
}

.options-list {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.option-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.option-item.correct-answer {
  background-color: #e8f5e9;
  border-left: 4px solid #4caf50;
}

.option-letter {
  font-weight: bold;
  margin-right: 8px;
}

.question-image {
  margin-top: 16px;
  text-align: center;
}

.question-image img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 4px;
}

/* Modal styles */
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
  width: 400px;
  max-width: 90%;
}

.modal-content h3 {
  margin-top: 0;
}

.warning {
  color: #e53935;
  font-weight: 500;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 24px;
}

.cancel-btn, .delete-confirm-btn {
  padding: 10px 16px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  border: none;
}

.cancel-btn {
  background-color: #f5f5f5;
}

.delete-confirm-btn {
  background-color: #e53935;
  color: white;
}

.delete-confirm-btn:disabled {
  background-color: #ffcdd2;
  cursor: not-allowed;
}
</style>
