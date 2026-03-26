import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useChatStore = defineStore('chat', () => {
  // 状态
  const conversations = ref([])
  const currentConversation = ref(null)
  const messages = ref({}) // { userId: [messages] }
  const wsConnected = ref(false)

  // Actions
  function setConversations(list) {
    conversations.value = list
  }

  function setCurrentConversation(userId) {
    currentConversation.value = userId
  }

  function addMessage(message) {
    const key = message.sender_id === currentConversation.value
      ? 'current'
      : message.sender_id.toString()

    if (!messages.value[key]) {
      messages.value[key] = []
    }
    messages.value[key].push(message)
  }

  function setMessages(userId, msgList) {
    messages.value[userId.toString()] = msgList
  }

  function setWsConnected(connected) {
    wsConnected.value = connected
  }

  return {
    conversations,
    currentConversation,
    messages,
    wsConnected,
    setConversations,
    setCurrentConversation,
    addMessage,
    setMessages,
    setWsConnected
  }
})
