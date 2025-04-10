<template>
  <div class="home">
    <!-- Intro Section dengan GSAP -->
    <section class="intro" v-if="showIntro">
      <div class="grid-container">
        <div class="grid-item slide-up" ref="item1">
          <img src="/assets/images/sekolah1.jpg" alt="Sekolah 1" />
        </div>
        <div class="grid-item slide-up" ref="item2">
          <img src="/assets/images/sekolah2.jpg" alt="Sekolah 2" />
        </div>
        <div class="grid-item slide-down" ref="item3">
          <img src="/assets/images/sekolah3.jpg" alt="Sekolah 3" />
        </div>
        <div class="grid-item slide-down" ref="item4">
          <img src="/assets/images/sekolah4.jpg" alt="Sekolah 4" />
        </div>
        <div class="grid-item slide-up" ref="item5">
          <img src="/assets/images/sekolah5.jpg" alt="Sekolah 5" />
        </div>
      </div>
    </section>

    <!-- Main Content - Menggunakan position absolute dan z-index -->
    <div ref="mainContent" :class="['main-content', { 'show': !showIntro }]">
      <div class="hero">
        <h1>Selamat Datang di LMS Sekolah</h1>
        <p>Sistem Pembelajaran Online untuk Siswa dan Guru</p>

        <div class="status-card" v-if="serverStatus">
          <p><strong>Status Server:</strong> {{ serverStatus.status }}</p>
          <p>
            <small>{{ serverStatus.message }}</small>
          </p>
          <p>
            <small>{{ serverStatus.time }}</small>
          </p>
        </div>
      </div>

      <div class="feature-section">
        <div class="feature-card" ref="featureCard1">
          <h2>Siswa</h2>
          <p>Lihat daftar siswa dan data mereka</p>
          <router-link to="/students" class="btn">Lihat Siswa</router-link>
        </div>

        <div class="feature-card" ref="featureCard2">
          <h2>Soal</h2>
          <p>Akses soal-soal dan jawab kuis</p>
          <router-link to="/questions" class="btn">Lihat Soal</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { serverStatus } from "@/services/api";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";

export default {
  name: "HomeView",
  data() {
    return {
      serverStatus: null,
      showIntro: true,
    };
  },
  created() {
    this.checkServerStatus();

    // Cek apakah intro sudah dilihat sebelumnya
    const introSeen = sessionStorage.getItem("introSeen");
    this.showIntro = !introSeen;
  },
  mounted() {
    gsap.registerPlugin(ScrollTrigger);

    if (this.showIntro) {
      // Mencegah scroll default
      document.body.style.overflow = 'hidden';
      this.setupIntroAnimation();
    } else {
      document.body.style.overflow = 'auto';
      this.setupMainContentAnimation();
    }
  },
  methods: {
    async checkServerStatus() {
      try {
        const response = await serverStatus.check();
        this.serverStatus = response.data;
      } catch (error) {
        console.error("Error connecting to server", error);
      }
    },
    setupIntroAnimation() {
      // Pastikan semua referensi ada sebelum setup animasi
      this.$nextTick(() => {
        // Cek apakah semua ref tersedia
        if (!this.$refs.item1 || !this.$refs.item2 || !this.$refs.item3 ||
            !this.$refs.item4 || !this.$refs.item5 || !this.$refs.mainContent) {
          console.error("Some refs are not available for animation");
          // Fallback ke main content jika ada masalah dengan refs
          this.showIntro = false;
          document.body.style.overflow = 'auto';
          return;
        }

        // Set initial state for all images
        gsap.set([this.$refs.item1, this.$refs.item2, this.$refs.item3, this.$refs.item4, this.$refs.item5], {
          opacity: 1,
          y: 0
        });

        // Variabel untuk menyimpan fungsi wheel handler
        const wheelEventHandler = (e) => {
          e.preventDefault();
          wheelHandler();
        };

        const wheelHandler = () => {
          // Check if refs exist before animating
          const itemsUp = [this.$refs.item1, this.$refs.item2, this.$refs.item5].filter(item => item);
          const itemsDown = [this.$refs.item3, this.$refs.item4].filter(item => item);

          // Only animate if there are elements to animate
          if (itemsUp.length > 0) {
            // Animasi slide up dengan pergerakan yang ekstrim (keluar viewport)
            gsap.to(itemsUp, {
              y: -800,
              opacity: 1,
              duration: 0.8,
              ease: "power2.out"
            });
          }

          if (itemsDown.length > 0) {
            // Animasi slide down dengan pergerakan yang ekstrim (keluar viewport)
            gsap.to(itemsDown, {
              y: 800,
              opacity: 1,
              duration: 0.8,
              ease: "power2.out"
            });
          }

          // Siapkan main content
          setTimeout(() => {
            // Check if mainContent ref exists
            if (this.$refs.mainContent) {
              gsap.to(this.$refs.mainContent, {
                opacity: 1,
                duration: 0.5,
                onComplete: () => {
                  this.showIntro = false;
                  sessionStorage.setItem("introSeen", "true");
                  document.body.style.overflow = 'auto';

                  // Setup animasi main content
                  this.$nextTick(() => {
                    this.setupMainContentAnimation();
                  });
                }
              });
            } else {
              // Fallback if mainContent ref doesn't exist
              this.showIntro = false;
              sessionStorage.setItem("introSeen", "true");
              document.body.style.overflow = 'auto';

              // Setup animasi main content
              this.$nextTick(() => {
                this.setupMainContentAnimation();
              });
            }
          }, 400);

          // Remove event listeners
          window.removeEventListener('wheel', wheelEventHandler);
          window.removeEventListener('touchmove', touchMoveHandler);
          window.removeEventListener('touchstart', touchStartHandler);
        };

        // Touch events
        let touchStarted = false;
        const touchStartHandler = () => {
          touchStarted = true;
        };

        const touchMoveHandler = (e) => {
          if (touchStarted) {
            e.preventDefault();

            // Trigger the same animation as wheel
            wheelHandler();

            touchStarted = false;
          }
        };

        // Add event listeners
        window.addEventListener('wheel', wheelEventHandler, { passive: false });
        window.addEventListener('touchstart', touchStartHandler, { passive: true });
        window.addEventListener('touchmove', touchMoveHandler, { passive: false });

        // Cleanup in beforeUnmount
        this.$options.beforeUnmount = () => {
          window.removeEventListener('wheel', wheelEventHandler);
          window.removeEventListener('touchstart', touchStartHandler);
          window.removeEventListener('touchmove', touchMoveHandler);
        };
      });
    },
    setupMainContentAnimation() {
      this.$nextTick(() => {
        // Get available feature cards
        const featureCards = [];
        if (this.$refs.featureCard1) featureCards.push(this.$refs.featureCard1);
        if (this.$refs.featureCard2) featureCards.push(this.$refs.featureCard2);

        // Only animate if there are feature cards available
        if (featureCards.length === 0) {
          console.warn("No feature card refs are available for animation");
          return;
        }

        // Animasi feature cards pop-out
        gsap.from(featureCards, {
          scale: 0.8,
          opacity: 0,
          duration: 0.8,
          stagger: 0.2,
          ease: "back.out(1.7)"
        });
      });
    }
  }
};
</script>

<style scoped>
/* Main Content Styles */
.home {
  width: 100%;
  overflow-x: hidden;
  position: relative;
  min-height: 100vh;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  opacity: 0;
  transition: opacity 1s ease;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  z-index: -1;
}

.main-content.show {
  opacity: 1;
  z-index: 10;
}

/* Intro Styles */
.intro {
  height: 100vh;
  width: 100%;
  position: relative;
  background-color: #fff;
  overflow: hidden;
  z-index: 5;
}

.grid-container {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 15px;
  height: 100vh;
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
  position: relative;
}

.grid-item {
  overflow: hidden;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  position: relative;
  opacity: 1;
  min-height: 300px;
  background-color: #f8f8f8;
}

.grid-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  transition: transform 0.3s ease;
  opacity: 1;
  display: block;
}

.grid-item:hover img {
  transform: scale(1.05);
}

/* Grid Positioning - Sesuai sketsa */
.grid-item:nth-child(1) {
  grid-column: 1;
  grid-row: 1;
}

.grid-item:nth-child(2) {
  grid-column: 2;
  grid-row: 1;
}

.grid-item:nth-child(3) {
  grid-column: 1;
  grid-row: 2;
}

.grid-item:nth-child(4) {
  grid-column: 2;
  grid-row: 2;
}

.grid-item:nth-child(5) {
  grid-column: 3;
  grid-row: 1 / span 2;
}

/* Existing Styles */
.hero {
  text-align: center;
  margin-bottom: 40px;
  padding: 40px 20px;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.hero h1 {
  font-size: 2.5em;
  margin-bottom: 10px;
  color: #2c3e50;
}

.hero p {
  font-size: 1.2em;
  color: #666;
}

.status-card {
  margin-top: 20px;
  padding: 15px;
  background-color: #e8f5e9;
  border-radius: 8px;
  display: inline-block;
}

.feature-section {
  display: flex;
  gap: 20px;
  justify-content: center;
  flex-wrap: wrap;
}

.feature-card {
  flex: 1;
  min-width: 300px;
  padding: 30px;
  border-radius: 8px;
  background-color: #fff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.feature-card h2 {
  color: #42b983;
  margin-bottom: 15px;
}

.btn {
  display: inline-block;
  margin-top: 15px;
  padding: 10px 20px;
  background-color: #42b983;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  font-weight: bold;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: #3aa876;
}

/* Responsive Design */
@media (max-width: 768px) {
  .grid-container {
    grid-template-columns: 1fr 1fr;
    grid-template-rows: repeat(3, 1fr);
  }

  .grid-item:nth-child(3) {
    grid-column: 1 / span 2;
    grid-row: 2;
  }

  .grid-item:nth-child(4) {
    grid-column: 1;
    grid-row: 3;
  }

  .grid-item:nth-child(5) {
    grid-column: 2;
    grid-row: 3;
  }
}

@media (max-width: 480px) {
  .grid-container {
    grid-template-columns: 1fr;
    grid-template-rows: repeat(5, 1fr);
  }

  .grid-item:nth-child(n) {
    grid-column: 1;
  }

  .grid-item:nth-child(1) { grid-row: 1; }
  .grid-item:nth-child(2) { grid-row: 2; }
  .grid-item:nth-child(3) { grid-row: 3; }
  .grid-item:nth-child(4) { grid-row: 4; }
  .grid-item:nth-child(5) { grid-row: 5; }
}
</style>
