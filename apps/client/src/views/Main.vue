<template>
  <div class="main-container">
    <!-- 左侧会话列表 -->
    <div class="sidebar">
      <div class="sidebar-header">
        <div class="user-info">
          <el-avatar :size="36" :src="userStore.userInfo?.avatar || ''">
            {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
          </el-avatar>
          <span class="nickname">{{ userStore.userInfo?.nickname || userStore.userInfo?.phone }}</span>
        </div>
        <el-button text @click="handleLogout" title="退出登录">
          <el-icon><SwitchButton /></el-icon>
        </el-button>
      </div>

      <div class="conversation-list">
        <div class="empty-state" v-if="conversations.length === 0">
          <el-empty description="暂无会话" :image-size="80" />
        </div>
      </div>
    </div>

    <!-- 右侧聊天区域 -->
    <div class="chat-area">
      <div class="chat-header" v-if="chatStore.currentConversation">
        <span>聊天窗口</span>
      </div>
      <div class="chat-empty" v-else>
        <el-empty description="选择或开始一个新的会话" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useChatStore } from '@/stores/chat'
import { ElMessageBox } from 'element-plus'
import { SwitchButton } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const chatStore = useChatStore()

// WebSocket 连接
let ws = null

const connectWebSocket = () => {
  const token = userStore.token
  if (!token) return

  const wsUrl = `ws://localhost:8080/ws?token=${token}`
  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    console.log('WebSocket connected')
    chatStore.setWsConnected(true)
  }

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    console.log('WebSocket message:', data)

    if (data.type === 'message') {
      chatStore.addMessage({
        sender_id: data.data.from,
        content: data.data.content,
        timestamp: data.data.timestamp
      })
    }
  }

  ws.onclose = () => {
    console.log('WebSocket closed')
    chatStore.setWsConnected(false)

    // 尝试重连
    setTimeout(() => {
      if (userStore.isLoggedIn) {
        connectWebSocket()
      }
    }, 5000)
  }

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
  }
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    // 关闭 WebSocket 连接
    if (ws) {
      ws.close()
    }
    userStore.logout()
    router.push('/login')
  })
}

onMounted(() => {
  connectWebSocket()
})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
})
</script>

<style scoped>
.main-container {
  display: flex;
  width: 100%;
  height: 100vh;
}

.sidebar {
  width: 280px;
  background: #f5f5f5;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 16px;
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nickname {
  font-size: 14px;
  font-weight: 500;
}

.conversation-list {
  flex: 1;
  overflow-y: auto;
}

.empty-state {
  padding: 40px 20px;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chat-header {
  padding: 16px;
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  font-size: 16px;
  font-weight: 500;
}

.chat-empty {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #fafafa;
}
</style>
