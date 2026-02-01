<template>
  <div class="app-container">
    <!-- 导航栏 -->
    <header class="navbar">
      <div class="container">
        <div class="logo">
          <h1>前端文件公益镜像服务</h1>
        </div>
        <nav class="nav-links">
          <a href="#home" class="nav-link">首页</a>
          <a href="#tools" class="nav-link">工具</a>
          <a href="#stats" class="nav-link">数据</a>
          <a href="#about" class="nav-link">关于</a>
        </nav>
      </div>
    </header>

    <!-- 首页/英雄区 -->
    <section id="home" class="hero">
      <div class="container">
        <div class="hero-content">
          <h2>为中国大陆开发者提供前端库镜像加速服务</h2>
          <p>解决前端资源访问慢、不稳定的问题，提供高速、可靠的镜像服务</p>
          <div class="hero-features">
            <div class="feature-item">
              <el-icon><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"></path></svg></el-icon>
              <span>高速镜像</span>
            </div>
            <div class="feature-item">
              <el-icon><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"></path></svg></el-icon>
              <span>稳定可靠</span>
            </div>
            <div class="feature-item">
              <el-icon><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M12 6v6l4 2"></path></svg></el-icon>
              <span>公益免费</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 工具区 -->
    <section id="tools" class="tools-section">
      <div class="container">
        <h3 class="section-title">加速工具</h3>
        <div class="tools-container">
          <!-- URL输入工具 -->
          <div class="tool-card">
            <h4>URL加速</h4>
            <el-input
              v-model="inputUrl"
              type="textarea"
              :rows="3"
              placeholder="请输入需要加速的URL，例如：https://cdn.jsdelivr.net/npm/vue@3/dist/vue.global.js"
              class="url-input"
            ></el-input>
            <div class="tool-buttons">
              <el-button type="primary" @click="processUrl" :loading="processing">
                处理
              </el-button>
              <el-button @click="testLatency" :loading="testing">
                测试延迟
              </el-button>
            </div>
            
            <!-- 处理结果 -->
            <div v-if="resultUrl" class="result-section">
              <h5>加速结果</h5>
              <div class="result-item">
                <span class="result-label">原URL:</span>
                <el-input v-model="inputUrl" readonly class="result-input"></el-input>
              </div>
              <div class="result-item">
                <span class="result-label">加速URL:</span>
                <el-input v-model="resultUrl" readonly class="result-input"></el-input>
                <el-button type="success" @click="copyUrl(resultUrl)" size="small">
                  <el-icon><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg></el-icon>
                  复制
                </el-button>
              </div>
            </div>
            
            <!-- 延迟测试结果 -->
            <div v-if="latencyResult" class="latency-result">
              <h5>延迟测试结果</h5>
              <div class="latency-comparison">
                <div class="latency-item">
                  <span class="latency-label">原地址:</span>
                  <span class="latency-value">{{ latencyResult.original_latency }}ms</span>
                </div>
                <div class="latency-item">
                  <span class="latency-label">加速地址:</span>
                  <span class="latency-value">{{ latencyResult.accelerated_latency }}ms</span>
                </div>
                <div class="latency-item">
                  <span class="latency-label">提升:</span>
                  <span class="latency-improvement">{{ latencyResult.improvement }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 数据统计区 -->
    <section id="stats" class="stats-section">
      <div class="container">
        <h3 class="section-title">运营数据</h3>
        <div class="stats-grid">
          <!-- 统计卡片 -->
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"></path></svg></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ statsData.requests || 0 }}</div>
              <div class="stat-label">总请求数</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path></svg></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ statsData.bandwidth || '0 B' }}</div>
              <div class="stat-label">总流量</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect><rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect><line x1="6" y1="6" x2="6.01" y2="6"></line><line x1="6" y1="18" x2="6.01" y2="18"></line></svg></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ statsData.today_requests || 0 }}</div>
              <div class="stat-label">今日请求</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ statsData.today_traffic || '0 B' }}</div>
              <div class="stat-label">今日流量</div>
            </div>
          </div>
        </div>

        <!-- 热门源站 -->
        <div class="top-sources">
          <h4>热门源站</h4>
          <div class="sources-list">
            <el-tag v-for="(source, index) in statsData.top_sources" :key="index" class="source-tag">
              {{ source }}
            </el-tag>
          </div>
        </div>
      </div>
    </section>

    <!-- 关于/服务介绍区 -->
    <section id="about" class="about-section">
      <div class="container">
        <h3 class="section-title">关于服务</h3>
        <div class="about-content">
          <div class="about-card">
            <h4>服务介绍</h4>
            <p>前端文件公益镜像服务是一个为中国大陆开发者提供的公益项目，旨在解决前端资源访问慢、不稳定的问题。通过反代多个国际知名的前端资源CDN，为开发者提供高速、稳定的镜像加速服务。</p>
          </div>
          
          <div class="about-card">
            <h4>支持的源站</h4>
            <ul class="sources-ul">
              <li><a href="https://cdn.jsdelivr.net" target="_blank">cdn.jsdelivr.net</a></li>
              <li><a href="https://cdnjs.cloudflare.com" target="_blank">cdnjs.cloudflare.com</a></li>
              <li><a href="https://ghcr.io" target="_blank">ghcr.io</a></li>
              <li><a href="https://registry-1.docker.io" target="_blank">registry-1.docker.io</a></li>
              <li><a href="https://unpkg.com" target="_blank">unpkg.com</a></li>
            </ul>
          </div>
          
          <div class="about-card">
            <h4>服务优势</h4>
            <ul class="advantages-ul">
              <li>高速：通过优化网络路由和缓存策略，提供更快的访问速度</li>
              <li>稳定：多源站备份，确保服务可用性</li>
              <li>安全：严格的URL验证和过滤机制，防止滥用</li>
              <li>透明：完整的访问统计和监控，确保服务质量</li>
              <li>公益：完全免费，为开发者社区贡献力量</li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-info">
            <h4>前端文件公益镜像服务</h4>
            <p>为中国大陆开发者提供前端库镜像加速服务</p>
          </div>
          <div class="footer-links">
            <a href="#home">首页</a>
            <a href="#tools">工具</a>
            <a href="#stats">数据</a>
            <a href="#about">关于</a>
          </div>
        </div>
        <div class="footer-bottom">
          <p>© 2026 前端文件公益镜像服务 - 公益项目</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

// 响应式数据
const inputUrl = ref('')
const resultUrl = ref('')
const latencyResult = ref(null)
const processing = ref(false)
const testing = ref(false)
const statsData = ref({})

// 处理URL
const processUrl = async () => {
  if (!inputUrl.value) {
    ElMessage.warning('请输入需要加速的URL')
    return
  }

  processing.value = true
  try {
    const response = await axios.post('/api/process-url', {
      url: inputUrl.value
    })
    resultUrl.value = response.data.accelerated_url
    ElMessage.success('URL处理成功')
  } catch (error) {
    ElMessage.error(`处理失败: ${error.response?.data?.error || '未知错误'}`)
  } finally {
    processing.value = false
  }
}

// 测试延迟
const testLatency = async () => {
  if (!inputUrl.value) {
    ElMessage.warning('请输入需要测试的URL')
    return
  }

  testing.value = true
  try {
    const response = await axios.post('/api/test-latency', {
      url: inputUrl.value
    })
    latencyResult.value = response.data
    resultUrl.value = response.data.accelerated_url
    ElMessage.success('延迟测试完成')
  } catch (error) {
    ElMessage.error(`测试失败: ${error.response?.data?.error || '未知错误'}`)
  } finally {
    testing.value = false
  }
}

// 复制URL
const copyUrl = (url) => {
  navigator.clipboard.writeText(url).then(() => {
    ElMessage.success('URL已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败，请手动复制')
  })
}

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await axios.get('/api/stats')
    statsData.value = response.data
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 组件挂载时获取统计数据
onMounted(() => {
  fetchStats()
  // 每30秒刷新一次统计数据
  setInterval(fetchStats, 30000)
})
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* 导航栏 */
.navbar {
  background-color: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.navbar .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
}

.logo h1 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1890ff;
  margin: 0;
}

.nav-links {
  display: flex;
  gap: 2rem;
}

.nav-link {
  color: #333;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s;
}

.nav-link:hover {
  color: #1890ff;
}

/* 英雄区 */
.hero {
  background: linear-gradient(135deg, #1890ff 0%, #36cfc9 100%);
  color: white;
  padding: 4rem 0;
}

.hero-content {
  text-align: center;
}

.hero-content h2 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  font-weight: 600;
}

.hero-content p {
  font-size: 1.2rem;
  margin-bottom: 2rem;
  opacity: 0.9;
}

.hero-features {
  display: flex;
  justify-content: center;
  gap: 3rem;
  margin-top: 3rem;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.feature-item .el-icon {
  font-size: 2rem;
}

/* 通用容器 */
.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

/* 通用部分标题 */
.section-title {
  font-size: 2rem;
  font-weight: 600;
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
}

/* 工具区 */
.tools-section {
  padding: 4rem 0;
  background-color: #ffffff;
}

.tools-container {
  display: flex;
  justify-content: center;
}

.tool-card {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 2rem;
  width: 100%;
  max-width: 800px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.tool-card h4 {
  font-size: 1.2rem;
  margin-bottom: 1rem;
  color: #333;
}

.url-input {
  margin-bottom: 1rem;
}

.tool-buttons {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.result-section {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e8e8e8;
}

.result-section h5 {
  margin-bottom: 1rem;
  color: #333;
}

.result-item {
  margin-bottom: 1rem;
}

.result-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #666;
}

.result-input {
  margin-bottom: 0.5rem;
}

.latency-result {
  margin-top: 1.5rem;
  padding: 1rem;
  background-color: #f0f9eb;
  border-radius: 4px;
  border-left: 4px solid #52c41a;
}

.latency-result h5 {
  margin-bottom: 1rem;
  color: #333;
}

.latency-comparison {
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
}

.latency-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.latency-label {
  font-weight: 500;
  color: #666;
}

.latency-value {
  font-weight: 600;
  color: #333;
}

.latency-improvement {
  font-weight: 600;
  color: #52c41a;
}

/* 数据统计区 */
.stats-section {
  padding: 4rem 0;
  background-color: #f5f7fa;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.stat-card {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 1rem;
}

.stat-icon {
  font-size: 2rem;
  color: #1890ff;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 600;
  color: #333;
}

.stat-label {
  font-size: 0.9rem;
  color: #666;
  margin-top: 0.25rem;
}

.top-sources {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.top-sources h4 {
  margin-bottom: 1rem;
  color: #333;
}

.sources-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.source-tag {
  margin: 0;
}

/* 关于区 */
.about-section {
  padding: 4rem 0;
  background-color: #ffffff;
}

.about-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.about-card {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.about-card h4 {
  font-size: 1.2rem;
  margin-bottom: 1rem;
  color: #333;
}

.about-card p {
  line-height: 1.6;
  color: #666;
}

.sources-ul, .advantages-ul {
  list-style: none;
  padding: 0;
}

.sources-ul li, .advantages-ul li {
  margin-bottom: 0.5rem;
  padding-left: 1.5rem;
  position: relative;
  line-height: 1.5;
  color: #666;
}

.sources-ul li::before, .advantages-ul li::before {
  content: "•";
  color: #1890ff;
  font-weight: bold;
  position: absolute;
  left: 0;
}

.sources-ul a {
  color: #1890ff;
  text-decoration: none;
}

.sources-ul a:hover {
  text-decoration: underline;
}

/* 页脚 */
.footer {
  background-color: #333;
  color: #ffffff;
  padding: 2rem 0;
  margin-top: auto;
}

.footer-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 2rem;
}

.footer-info h4 {
  margin-bottom: 0.5rem;
  font-size: 1.2rem;
}

.footer-info p {
  opacity: 0.8;
}

.footer-links {
  display: flex;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.footer-links a {
  color: #ffffff;
  text-decoration: none;
  opacity: 0.8;
  transition: opacity 0.3s;
}

.footer-links a:hover {
  opacity: 1;
}

.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding-top: 1rem;
  text-align: center;
  opacity: 0.8;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .hero-content h2 {
    font-size: 2rem;
  }

  .hero-features {
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
  }

  .nav-links {
    gap: 1rem;
  }

  .tool-buttons {
    flex-direction: column;
  }

  .latency-comparison {
    flex-direction: column;
    gap: 0.5rem;
  }

  .footer-content {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
}
</style>
