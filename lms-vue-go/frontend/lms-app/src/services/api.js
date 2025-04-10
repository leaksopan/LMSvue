import axios from "axios";

// Konfigurasi axios untuk panggilan API
const apiClient = axios.create({
  baseURL: "/api", // Menggunakan path relatif untuk memanfaatkan proxy devServer
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    "X-Requested-With": "XMLHttpRequest",
  },
  timeout: 10000, // 10 detik timeout
});

// Fallback client yang menggunakan URL langsung ke backend
const directApiClient = axios.create({
  baseURL: "http://localhost:3001/api", // URL langsung ke backend
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
    "X-Requested-With": "XMLHttpRequest",
  },
  timeout: 10000, // 10 detik timeout
});

// Function to add auth token to request headers
const addAuthToken = (request) => {
  // Ambil token dari localStorage jika ada
  const token = localStorage.getItem("auth_token");
  if (token) {
    request.headers["Authorization"] = `Bearer ${token}`;
  }
  return request;
};

// Tambahkan interceptors untuk menambahkan token otentikasi ke apiClient
apiClient.interceptors.request.use((request) => {
  // Add auth token
  addAuthToken(request);

  // Hanya log request di development mode
  if (process.env.NODE_ENV === "development") {
    console.log("API Request (proxy):", request);
  }

  return request;
});

// Tambahkan interceptors untuk menambahkan token otentikasi ke directApiClient
directApiClient.interceptors.request.use((request) => {
  // Add auth token
  addAuthToken(request);

  // Hanya log request di development mode
  if (process.env.NODE_ENV === "development") {
    console.log("API Request (direct):", request);
  }

  return request;
});

// Response interceptor for apiClient
apiClient.interceptors.response.use(
  (response) => {
    // Hanya log response di development mode
    if (process.env.NODE_ENV === "development") {
      console.log("API Response (proxy):", response);
    }
    return response;
  },
  (error) => {
    // Hanya log error di development mode
    if (process.env.NODE_ENV === "development") {
      console.error("API Error (proxy):", error);
    }

    // Jika error 401 (Unauthorized), redirect ke login
    if (error.response && error.response.status === 401) {
      // Clear token
      localStorage.removeItem("auth_token");
      localStorage.removeItem("user");

      // Redirect ke login jika tidak sedang di halaman login
      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
    }

    return Promise.reject(error);
  }
);

// Response interceptor for directApiClient
directApiClient.interceptors.response.use(
  (response) => {
    // Hanya log response di development mode
    if (process.env.NODE_ENV === "development") {
      console.log("API Response (direct):", response);
    }
    return response;
  },
  (error) => {
    // Hanya log error di development mode
    if (process.env.NODE_ENV === "development") {
      console.error("API Error (direct):", error);
    }

    // Jika error 401 (Unauthorized), redirect ke login
    if (error.response && error.response.status === 401) {
      // Clear token
      localStorage.removeItem("auth_token");
      localStorage.removeItem("user");

      // Redirect ke login jika tidak sedang di halaman login
      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
    }

    return Promise.reject(error);
  }
);

// Service untuk autentikasi
export const authService = {
  // Login user
  async login(credentials) {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.post("/auth/login", credentials);
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  },

  // Register user
  async register(userData) {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.post("/auth/register", userData);
    } catch (error) {
      console.error("Register failed:", error);
      throw error;
    }
  },

  // Get current user
  async getCurrentUser() {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.get("/auth/me");
    } catch (error) {
      console.error("Get current user failed:", error);
      throw error;
    }
  },

  // Logout user (client-side only)
  logout() {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("user");
  },

  // Check if user is logged in
  isLoggedIn() {
    return !!localStorage.getItem("auth_token");
  },

  // Get current user from localStorage
  getUser() {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch (e) {
        return null;
      }
    }
    return null;
  },
};

// Service untuk data siswa
export const studentService = {
  // Mengambil semua data siswa
  async getAll() {
    try {
      // Use the public students endpoint that doesn't require authentication
      const response = await axios.get(
        "http://localhost:3001/api/public/students"
      );
      return response;
    } catch (error) {
      console.error("Get all students failed:", error);
      throw error;
    }
  },

  // Mengambil data siswa berdasarkan ID
  async getById(id) {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.get(`/students/${id}`);
    } catch (error) {
      console.error(`Get student by ID ${id} failed:`, error);
      throw error;
    }
  },

  // Mengambil profil siswa untuk user yang sedang login
  async getCurrentProfile() {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.get("/students/profile/me");
    } catch (error) {
      console.error("Get current student profile failed:", error);
      throw error;
    }
  },
};

// Service untuk soal-soal
export const questionService = {
  // Mengambil semua soal
  async getAll() {
    try {
      // Use the public questions endpoint that doesn't require authentication
      const response = await axios.get(
        "http://localhost:3001/api/public/questions"
      );
      return response;
    } catch (error) {
      console.error("Get all questions failed:", error);
      throw error;
    }
  },

  // Mengambil soal berdasarkan ID
  async getById(id) {
    try {
      // Create a custom axios instance for this specific request
      const token = localStorage.getItem("auth_token");
      const response = await axios.get(
        `http://localhost:3001/api/questions/${id}`,
        {
          headers: {
            Authorization: token ? `Bearer ${token}` : "",
          },
        }
      );
      return response;
    } catch (error) {
      console.error(`Get question by ID ${id} failed:`, error);
      throw error;
    }
  },

  // Membuat soal baru (admin only)
  async create(questionData) {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.post("/questions", questionData);
    } catch (error) {
      console.error("Create question failed:", error);
      throw error;
    }
  },

  // Mengupdate soal (admin only)
  async update(id, questionData) {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.put(`/questions/${id}`, questionData);
    } catch (error) {
      console.error(`Update question ${id} failed:`, error);
      throw error;
    }
  },

  // Menghapus soal (admin only)
  async delete(id) {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.delete(`/questions/${id}`);
    } catch (error) {
      console.error(`Delete question ${id} failed:`, error);
      throw error;
    }
  },
};

// Service untuk jawaban siswa
export const answerService = {
  // Mengambil semua jawaban siswa yang sedang login
  async getMyAnswers() {
    try {
      // Use the public answers endpoint
      const response = await axios.get(
        "http://localhost:3001/api/public-answers"
      );
      return response;
    } catch (error) {
      console.error("Get my answers failed:", error);
      throw error;
    }
  },

  // Mengambil jawaban siswa untuk soal tertentu
  async getByQuestion(questionId) {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.get(`/answers/my/question/${questionId}`);
    } catch (error) {
      console.error(`Get answer for question ${questionId} failed:`, error);
      throw error;
    }
  },

  // Mengirimkan jawaban siswa
  async submitAnswer(questionId, answer) {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.post("/answers/submit", {
        question_id: questionId,
        answer: answer,
      });
    } catch (error) {
      console.error("Submit answer failed:", error);
      throw error;
    }
  },

  // Mengambil semua jawaban siswa (admin/guru)
  async getAllAnswers() {
    try {
      // Use the public all answers endpoint
      const response = await axios.get(
        "http://localhost:3001/api/public-all-answers"
      );
      return response;
    } catch (error) {
      console.error("Get all answers failed:", error);
      throw error;
    }
  },

  // Memberikan nilai untuk jawaban siswa (admin/guru)
  async gradeAnswer(answerId, score) {
    try {
      // Use the public grade answer endpoint
      const response = await axios.put(
        `http://localhost:3001/api/public-grade-answer/${answerId}`,
        {
          score: score,
        }
      );
      return response;
    } catch (error) {
      console.error(`Grade answer ${answerId} failed:`, error);
      throw error;
    }
  },
};

// Cek status server
export const serverStatus = {
  async check() {
    try {
      // Always use direct URL to avoid proxy issues
      return await directApiClient.get("/status");
    } catch (error) {
      console.error("Server status check failed:", error);
      throw error;
    }
  },
};

export default {
  authService,
  studentService,
  questionService,
  serverStatus,
};
