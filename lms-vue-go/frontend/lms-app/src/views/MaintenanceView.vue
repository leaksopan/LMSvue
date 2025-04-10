<template>
  <div class="maintenance-container">
    <h1>Halaman Sedang Dalam Perbaikan</h1>
    
    <div class="maintenance-content">
      <DotLottieVue 
        class="lottie-animation"
        autoplay 
        loop 
        src="/assets/eror404.lottie" 
      />
    </div>
    
    <p>Mohon maaf atas ketidaknyamanannya. Kami sedang bekerja keras untuk segera menyelesaikannya.</p>
    <router-link to="/" class="btn">Kembali ke Beranda</router-link>
  </div>
</template>

<script setup>
import { DotLottieVue } from '@lottiefiles/dotlottie-vue';
</script>

<style scoped>
.maintenance-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 80vh; 
  text-align: center;
  padding: 20px;
}

h1 {
  color: #333;
  margin-bottom: 20px;
}

.maintenance-content {
  margin-bottom: 20px;
  width: 100%;
  max-width: 600px;
  display: flex;
  justify-content: center;
}

.lottie-animation {
  width: 350px;
  height: 350px;
}

p {
  color: #666;
  margin-bottom: 30px;
}

.btn {
  display: inline-block;
  padding: 10px 20px;
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