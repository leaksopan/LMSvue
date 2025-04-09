import axios from 'axios';

// Konfigurasi axios untuk panggilan API
const apiClient = axios.create({
  baseURL: '/api', // Menggunakan path relatif untuk memanfaatkan proxy devServer
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    'X-Requested-With': 'XMLHttpRequest'
  },
  timeout: 10000 // 10 detik timeout
});

// Tambahkan interceptors untuk debugging
apiClient.interceptors.request.use(request => {
  console.log('API Request:', request);
  return request;
});

apiClient.interceptors.response.use(
  response => {
    console.log('API Response:', response);
    return response;
  },
  error => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

// Service untuk data siswa
export const studentService = {
  // Mengambil semua data siswa
  getAll() {
    return apiClient.get('/students');
  },
  
  // Mengambil data siswa berdasarkan ID
  getById(id) {
    return apiClient.get(`/students/${id}`);
  }
};

// Service untuk soal-soal
export const questionService = {
  // Mengambil semua soal
  getAll() {
    return apiClient.get('/questions');
  },
  
  // Mengambil soal berdasarkan ID
  getById(id) {
    return apiClient.get(`/questions/${id}`);
  }
};

// Cek status server
export const serverStatus = {
  check() {
    return apiClient.get('/status');
  }
};

export default {
  studentService,
  questionService,
  serverStatus
}; 