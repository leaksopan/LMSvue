<template>
  <div class="question-form">
    <h2>{{ isEdit ? 'Edit Soal' : 'Tambah Soal Baru' }}</h2>
    
    <form @submit.prevent="submitForm">
      <div class="form-group">
        <label for="question-type">Tipe Soal</label>
        <select id="question-type" v-model="form.type" required>
          <option value="multiple_choice">Pilihan Ganda</option>
          <option value="essay">Essay</option>
        </select>
      </div>
      
      <div class="form-group">
        <label for="question-text">Pertanyaan</label>
        <textarea 
          id="question-text" 
          v-model="form.question" 
          rows="4" 
          placeholder="Masukkan pertanyaan" 
          required
        ></textarea>
      </div>
      
      <div class="form-group" v-if="form.type === 'multiple_choice'">
        <label>Pilihan Jawaban</label>
        <div class="options-container">
          <div v-for="(option, index) in form.options" :key="index" class="option-item">
            <input 
              type="text" 
              v-model="form.options[index]" 
              :placeholder="`Pilihan ${String.fromCharCode(65 + index)}`"
              required
            />
            <button 
              type="button" 
              class="remove-option" 
              @click="removeOption(index)" 
              v-if="form.options.length > 2"
            >
              &times;
            </button>
          </div>
          
          <button 
            type="button" 
            class="add-option" 
            @click="addOption" 
            v-if="form.options.length < 5"
          >
            + Tambah Pilihan
          </button>
        </div>
      </div>
      
      <div class="form-group" v-if="form.type === 'multiple_choice'">
        <label for="question-answer">Jawaban Benar</label>
        <select id="question-answer" v-model="form.answer" required>
          <option v-for="(option, index) in form.options" :key="index" :value="String.fromCharCode(65 + index)">
            {{ String.fromCharCode(65 + index) }}
          </option>
        </select>
      </div>
      
      <div class="form-group">
        <label for="question-score">Nilai (Poin)</label>
        <input 
          type="number" 
          id="question-score" 
          v-model.number="form.score" 
          min="1" 
          max="100" 
          required
        />
      </div>
      
      <div class="form-group">
        <label for="question-image">URL Gambar (opsional)</label>
        <input 
          type="url" 
          id="question-image" 
          v-model="form.image_url" 
          placeholder="https://example.com/image.jpg"
        />
      </div>
      
      <div class="form-actions">
        <button type="button" class="cancel-btn" @click="$emit('cancel')">Batal</button>
        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? 'Menyimpan...' : (isEdit ? 'Simpan Perubahan' : 'Tambah Soal') }}
        </button>
      </div>
      
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </form>
  </div>
</template>

<script>
import { questionService } from '@/services/api';

export default {
  name: 'QuestionForm',
  props: {
    question: {
      type: Object,
      default: null
    },
    isEdit: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      form: {
        type: 'multiple_choice',
        question: '',
        options: ['', ''],
        answer: 'A',
        score: 1,
        image_url: ''
      },
      loading: false,
      error: null
    };
  },
  created() {
    // If editing, populate form with question data
    if (this.isEdit && this.question) {
      this.form.type = this.question.type;
      this.form.question = this.question.question;
      this.form.options = [...this.question.options];
      this.form.answer = this.question.answer;
      this.form.score = this.question.score;
      this.form.image_url = this.question.image_url || '';
      
      // Ensure we have at least 2 options for multiple choice
      if (this.form.type === 'multiple_choice' && (!this.form.options || this.form.options.length < 2)) {
        this.form.options = ['', ''];
      }
    }
  },
  methods: {
    addOption() {
      if (this.form.options.length < 5) {
        this.form.options.push('');
      }
    },
    removeOption(index) {
      if (this.form.options.length > 2) {
        this.form.options.splice(index, 1);
        
        // Update answer if the removed option was the correct answer
        const answerIndex = this.form.answer.charCodeAt(0) - 65;
        if (answerIndex === index) {
          this.form.answer = 'A'; // Default to first option
        } else if (answerIndex > index) {
          // Adjust answer if it was after the removed option
          this.form.answer = String.fromCharCode(this.form.answer.charCodeAt(0) - 1);
        }
      }
    },
    async submitForm() {
      this.loading = true;
      this.error = null;
      
      try {
        const questionData = {
          type: this.form.type,
          question: this.form.question,
          options: this.form.type === 'multiple_choice' ? this.form.options : [],
          answer: this.form.type === 'multiple_choice' ? this.form.answer : '',
          score: this.form.score,
          image_url: this.form.image_url || ''
        };
        
        let response;
        if (this.isEdit) {
          // Update existing question
          response = await questionService.update(this.question.id, questionData);
        } else {
          // Create new question
          response = await questionService.create(questionData);
        }
        
        this.$emit('saved', response.data.data);
      } catch (error) {
        console.error('Error saving question:', error);
        this.error = 'Gagal menyimpan soal. Silakan coba lagi.';
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.question-form {
  background-color: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  max-width: 800px;
  margin: 0 auto;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

input[type="text"],
input[type="number"],
input[type="url"],
textarea,
select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

textarea {
  resize: vertical;
}

.options-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.remove-option {
  background-color: #f44336;
  color: white;
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.add-option {
  background-color: #e0e0e0;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  margin-top: 10px;
  cursor: pointer;
  font-size: 14px;
}

.add-option:hover {
  background-color: #d0d0d0;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 24px;
}

.cancel-btn,
.submit-btn {
  padding: 10px 20px;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  border: none;
}

.cancel-btn {
  background-color: #e0e0e0;
}

.submit-btn {
  background-color: #4caf50;
  color: white;
}

.submit-btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}

.error-message {
  margin-top: 16px;
  color: #f44336;
  text-align: center;
}
</style>
