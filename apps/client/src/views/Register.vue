<template>
  <div class="register-container">
    <div class="register-box">
      <h1 class="title">CloudIM</h1>
      <p class="subtitle">注册账号</p>

      <el-form :model="form" :rules="rules" ref="formRef" class="register-form">
        <el-form-item prop="phone">
          <el-input
            v-model="form.phone"
            placeholder="请输入手机号"
            prefix-icon="User"
            maxlength="11"
          />
        </el-form-item>

        <el-form-item prop="captcha">
          <el-input
            v-model="form.captcha"
            placeholder="请输入验证码"
            prefix-icon="Key"
            maxlength="6"
          >
            <template #append>
              <el-button @click="handleSendCaptcha" :disabled="captchaDisabled">
                {{ captchaText }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码（至少 8 位，包含字母和数字）"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" class="register-btn" @click="handleRegister">
            注册
          </el-button>
        </el-form-item>

        <div class="links">
          <span>已有账号？</span>
          <router-link to="/login">立即登录</router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import api from '@/api'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)
const captchaDisabled = ref(false)
const captchaText = ref('获取验证码')

const form = reactive({
  phone: '',
  captcha: '',
  password: ''
})

const rules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为 6 位数字', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, message: '密码至少 8 位', trigger: 'blur' },
    { pattern: /^(?=.*[A-Za-z])(?=.*\d)/, message: '密码需包含字母和数字', trigger: 'blur' }
  ]
}

const handleSendCaptcha = () => {
  if (!form.phone || !/^1[3-9]\d{9}$/.test(form.phone)) {
    ElMessage.warning('请输入正确的手机号')
    return
  }

  // MVP 阶段：调用发送验证码接口
  api.post('/auth/captcha', { phone: form.phone })
    .then(() => {
      ElMessage.success('验证码已发送（MVP 阶段固定验证码：123456）')
      captchaDisabled.value = true
      let count = 60
      captchaText.value = `${count}秒后重发`
      const timer = setInterval(() => {
        count--
        if (count <= 0) {
          clearInterval(timer)
          captchaDisabled.value = false
          captchaText.value = '获取验证码'
        } else {
          captchaText.value = `${count}秒后重发`
        }
      }, 1000)
    })
    .catch(() => {
      ElMessage.error('发送验证码失败')
    })
}

const handleRegister = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const res = await api.post('/auth/register', form)
      if (res.code === 0) {
        userStore.login(res.data.token, res.data.user)
        ElMessage.success('注册成功')
        router.push('/')
      } else {
        ElMessage.error(res.message || '注册失败')
      }
    } catch (error) {
      ElMessage.error(error.response?.data?.message || '注册失败')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-box {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

.title {
  text-align: center;
  color: #333;
  margin-bottom: 8px;
}

.subtitle {
  text-align: center;
  color: #999;
  margin-bottom: 30px;
}

.register-form {
  margin-top: 20px;
}

.register-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
}

.links {
  text-align: center;
  margin-top: 16px;
  color: #999;
}

.links a {
  color: #667eea;
  text-decoration: none;
  margin-left: 8px;
}
</style>
