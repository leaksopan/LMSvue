<template>
  <div class="test-api">
    <h1>Test API Connectivity</h1>

    <div class="card">
      <h2>Test Endpoints</h2>
      <button @click="testApiStatus" class="btn">Test API Status</button>
      <button @click="testApiTest" class="btn">Test API Test Endpoint</button>
      <button @click="testApiAuth" class="btn">Test Auth Endpoint</button>
    </div>

    <div v-if="loading" class="loading">
      Loading...
    </div>

    <div v-if="error" class="error">
      <h3>Error:</h3>
      <pre>{{ error }}</pre>
    </div>

    <div v-if="result" class="result">
      <h3>Result:</h3>
      <pre>{{ JSON.stringify(result, null, 2) }}</pre>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'TestApiView',
  data() {
    return {
      loading: false,
      error: null,
      result: null
    };
  },
  methods: {
    async testApiStatus() {
      this.loading = true;
      this.error = null;
      this.result = null;

      try {
        // Use direct URL only
        const response = await axios.get('http://localhost:3001/api/status');
        this.result = response.data;
        console.log('API Status Response:', response);
      } catch (error) {
        console.error('API Status Error:', error);
        this.error = `${error.message}\n${JSON.stringify(error.response?.data || {}, null, 2)}`;
      } finally {
        this.loading = false;
      }
    },

    async testApiTest() {
      this.loading = true;
      this.error = null;
      this.result = null;

      try {
        // Use direct URL only
        const response = await axios.get('http://localhost:3001/api/test');
        this.result = response.data;
        console.log('API Test Response:', response);
      } catch (error) {
        console.error('API Test Error:', error);
        this.error = `${error.message}\n${JSON.stringify(error.response?.data || {}, null, 2)}`;
      } finally {
        this.loading = false;
      }
    },

    async testApiAuth() {
      this.loading = true;
      this.error = null;
      this.result = null;

      try {
        // Use direct URL only
        const response = await axios.post('http://localhost:3001/api/auth/login', {
          username: 'admin',
          password: 'admin123'
        });
        this.result = response.data;
        console.log('API Auth Response:', response);
      } catch (error) {
        console.error('API Auth Error:', error);
        this.error = `${error.message}\n${JSON.stringify(error.response?.data || {}, null, 2)}`;
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.test-api {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

h1 {
  text-align: center;
  margin-bottom: 20px;
}

.card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin-bottom: 20px;
}

.btn {
  background-color: #42b983;
  color: white;
  border: none;
  padding: 10px 15px;
  border-radius: 4px;
  cursor: pointer;
  margin-right: 10px;
  margin-bottom: 10px;
}

.btn:hover {
  background-color: #3aa876;
}

.loading {
  text-align: center;
  margin: 20px 0;
  font-weight: bold;
}

.error {
  background-color: #ffebee;
  color: #d32f2f;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
}

.result {
  background-color: #e8f5e9;
  padding: 15px;
  border-radius: 4px;
  overflow: auto;
}

pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
