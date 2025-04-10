<template>
  <div class="login-page">
    <div class="auth-container">
      <h1>{{ isRegister ? 'Daftar Akun' : 'Login' }}</h1>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>

      <form @submit.prevent="handleSubmit" class="auth-form">
        <div class="form-group">
          <label for="username">Username</label>
          <input
            type="text"
            id="username"
            v-model="form.username"
            required
            autocomplete="username"
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="form.password"
            required
            autocomplete="current-password"
          />
        </div>

        <div v-if="isRegister" class="form-group">
          <label for="email">Email</label>
          <input
            type="email"
            id="email"
            v-model="form.email"
            required
            autocomplete="email"
          />
        </div>

        <!-- Additional fields for student registration -->
        <div v-if="isRegister" class="form-group">
          <label for="name">Nama Lengkap</label>
          <input
            type="text"
            id="name"
            v-model="form.name"
            placeholder="Masukkan nama lengkap Anda"
            required
          />
        </div>

        <div v-if="isRegister" class="form-group">
          <label for="class">Kelas</label>
          <input
            type="text"
            id="class"
            v-model="form.class"
            placeholder="Contoh: 10A, 11B, 12C"
            required
          />
        </div>

        <div class="form-actions">
          <button
            type="submit"
            class="btn btn-primary"
            :disabled="loading"
          >
            {{ isRegister ? 'Daftar' : 'Login' }}
          </button>

          <div class="toggle-mode">
            <span v-if="isRegister">
              Sudah punya akun?
              <a href="#" @click.prevent="isRegister = false">Login</a>
            </span>
            <span v-else>
              Belum punya akun?
              <a href="#" @click.prevent="isRegister = true">Daftar</a>
            </span>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { authService } from '@/services/api';

export default {
  name: 'LoginView',
  data() {
    return {
      isRegister: false,
      loading: false,
      error: null,
      form: {
        username: '',
        password: '',
        email: '',
        name: '',
        class: '',
      }
    };
  },
  created() {
    // Redirect jika sudah login
    if (authService.isLoggedIn()) {
      this.$router.push('/');
    }
  },
  methods: {
    async handleSubmit() {
      this.loading = true;
      this.error = null;

      try {
        let response;

        if (this.isRegister) {
          // Validasi form
          if (!this.validateEmail(this.form.email)) {
            this.error = 'Email tidak valid';
            this.loading = false;
            return;
          }

          if (this.form.password.length < 6) {
            this.error = 'Password minimal 6 karakter';
            this.loading = false;
            return;
          }

          // Register
          response = await authService.register({
            username: this.form.username,
            password: this.form.password,
            email: this.form.email,
            role: 'student', // Default role
            name: this.form.name,
            class: this.form.class
          });
        } else {
          // Login
          response = await authService.login({
            username: this.form.username,
            password: this.form.password
          });
        }

        // Simpan token dan data user
        localStorage.setItem('auth_token', response.data.token);
        localStorage.setItem('user', JSON.stringify(response.data.user));

        // Redirect ke home
        this.$router.push('/');

      } catch (error) {
        console.error('Authentication error:', error);
        this.error = error.response?.data?.error || 'Terjadi kesalahan. Silakan coba lagi.';
      } finally {
        this.loading = false;
      }
    },

    validateEmail(email) {
      const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      return re.test(email);
    }
  }
};
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 150px);
  padding: 20px;
}

.auth-container {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 30px;
  width: 100%;
  max-width: 400px;
}

h1 {
  text-align: center;
  margin-bottom: 24px;
  color: #42b983;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

label {
  font-weight: 600;
  color: #333;
}

input {
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

input:focus {
  outline: none;
  border-color: #42b983;
  box-shadow: 0 0 0 2px rgba(66, 185, 131, 0.2);
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 8px;
}

.btn {
  padding: 12px;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-primary {
  background-color: #42b983;
  color: white;
}

.btn-primary:hover {
  background-color: #3aa876;
}

.btn:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.toggle-mode {
  text-align: center;
  margin-top: 8px;
}

.toggle-mode a {
  color: #42b983;
  text-decoration: none;
  font-weight: 600;
}

.toggle-mode a:hover {
  text-decoration: underline;
}

.error-message {
  background-color: #ffebee;
  color: #d32f2f;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 16px;
  text-align: center;
}
</style>
