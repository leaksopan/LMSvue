<template>
  <div class="questions">
    <h1>Daftar Soal</h1>

    <div v-if="loading" class="loading">
      <p>Memuat daftar soal...</p>
    </div>

    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <button @click="fetchQuestions" class="btn">Coba Lagi</button>
    </div>

    <div v-else>
      <div class="question-cards">
        <div
          v-for="question in questions"
          :key="question.id"
          class="question-card"
        >
          <div class="question-type">
            {{ formatQuestionType(question.type) }} <span class="question-score">({{ question.score }} poin)</span>
          </div>

          <h3>{{ question.question }}</h3>

          <div v-if="question.image_url" class="question-image">
            <img :src="question.image_url" alt="Gambar soal" />
          </div>

          <!-- SECTION: Tambah Form Jawaban & Submit -->
          <div class="answer-section" v-if="!question.isSubmitted">
            <!-- Form untuk soal pilihan ganda -->
            <div v-if="question.type === 'multiple_choice'" class="multiple-choice options">
              <div v-for="(option, index) in question.options" :key="index" class="option-item">
                <label :for="`option-${question.id}-${index}`" class="option-label">
                  <input
                    type="radio"
                    :id="`option-${question.id}-${index}`"
                    :value="String.fromCharCode(65 + index)"
                    v-model="question.userAnswer"
                    :name="`answer-option-${question.id}`"
                  >
                  <span class="option-text">
                    <strong>{{ String.fromCharCode(65 + index) }}.</strong> {{ option }}
                  </span>
                </label>
              </div>
            </div>

            <!-- Form untuk soal essay -->
            <div v-else-if="question.type === 'essay'" class="essay">
              <textarea
                v-model="question.userAnswer"
                placeholder="Tulis jawaban Anda di sini..."
                rows="4"
                class="essay-input"
              ></textarea>
            </div>

            <div class="submit-section">
              <button
                @click="submitAnswer(question)"
                class="btn btn-primary"
                :disabled="!question.userAnswer || question.isSubmitting"
              >
                {{ question.isSubmitting ? 'Mengirim...' : 'Kirim Jawaban' }}
              </button>
              <p v-if="question.submitError" class="submit-error">{{ question.submitError }}</p>
            </div>
          </div>
          <!-- END SECTION: Tambah Form Jawaban & Submit -->

          <!-- SECTION: Tampilkan feedback setelah submit -->
          <div v-if="question.isSubmitted" class="feedback success-message">
             <p>Jawaban Anda:
               <strong v-if="question.type === 'multiple_choice'">Opsi {{ question.userAnswer }}</strong>
               <strong v-else>{{ question.userAnswer }}</strong>
             </p>
             <p> Jawaban telah dikirim!</p>
          </div>
           <!-- END SECTION: Tampilkan feedback setelah submit -->
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { questionService, answerService } from "@/services/api";

export default {
  name: "QuestionsView",
  data() {
    return {
      questions: [],
      loading: true,
      error: null,
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
        this.questions = response.data.data.map(q => ({
          ...q,
          userAnswer: '',
          isSubmitted: false,
          submitError: null,
          isSubmitting: false
        }));
      } catch (error) {
        console.error("Error fetching questions", error);
        this.error = "Gagal memuat daftar soal. Silakan coba lagi nanti.";
      } finally {
        this.loading = false;
      }
    },

    formatQuestionType(type) {
      const typeMap = {
        multiple_choice: "Pilihan Ganda",
        essay: "Essay/Uraian",
      };

      return typeMap[type] || type;
    },

    async submitAnswer(question) {
      if (!question.userAnswer) {
        question.submitError = 'Jawaban tidak boleh kosong.';
        return;
      }

      question.isSubmitting = true;
      question.submitError = null;

      console.log('Jawaban dikirim:', {
        questionId: question.id,
        answer: question.userAnswer
      });

      try {
        // Submit answer to the server
        const response = await answerService.submitAnswer(question.id, question.userAnswer);
        console.log('Answer submitted successfully:', response.data);

        question.isSubmitted = true;

      } catch (error) {
        console.error('Error submitting answer for question', question.id, error);
        question.submitError = 'Gagal mengirim jawaban. Coba lagi.';
      } finally {
        question.isSubmitting = false;
      }
    }
  },
};
</script>

<style scoped>
.questions {
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

.question-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.question-card {
  padding: 20px;
  border-radius: 8px;
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: relative;
}

.question-type {
  position: absolute;
  top: 10px;
  right: 10px;
  background-color: #2196f3;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8rem;
}

.question-card h3 {
  margin-top: 10px;
  margin-bottom: 20px;
  color: #2c3e50;
}

.question-image {
  margin-bottom: 15px;
}

.question-image img {
  max-width: 100%;
  border-radius: 4px;
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

/* Tambahkan style untuk section jawaban */
.answer-section {
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #eee;
}

.multiple-choice .option-item {
  margin-bottom: 10px;
}

.multiple-choice .option-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.multiple-choice .option-label:hover {
  background-color: #f0f0f0;
}

.multiple-choice input[type="radio"] {
  margin-right: 10px;
}

.essay .essay-input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
  line-height: 1.4;
  box-sizing: border-box; /* Agar padding tidak menambah lebar */
}

.submit-section {
  margin-top: 15px;
  text-align: right; /* Posisikan tombol ke kanan */
}

.btn-primary {
  background-color: #42b983;
  color: white;
}

.btn-primary:hover {
  background-color: #3aa876;
}

.btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.submit-error {
  color: #e53935;
  font-size: 0.9rem;
  margin-top: 5px;
  text-align: right;
}

.feedback {
  margin-top: 15px;
  padding: 10px 15px;
  border-radius: 4px;
}

.success-message {
  background-color: #e8f5e9; /* Warna hijau muda */
  color: #2e7d32; /* Warna hijau tua */
  border: 1px solid #a5d6a7; /* Border hijau */
  text-align: center;
}

.question-score {
  font-size: 0.9rem;
  color: #666;
  font-weight: normal;
  margin-left: 5px;
}

/* END Tambahan style */
</style>
