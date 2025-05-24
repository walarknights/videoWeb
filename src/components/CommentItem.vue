<template>
  <div>
    <div class="comment-item" :class="{ 'is-reply': isReply }">
      <div class="comment-header" style="height: 60%">
        <div class="row" style="align-items: center; gap: 10px; height: 100%">
          <img
            :src="'http://localhost:8080' + comment.userAvatar"
            alt="..."
            style="border-radius: 50%; height: 60%; width: 25%"
            @click="toPersonal"
          />
          <strong class="user-name">{{ comment.userName }}</strong>
        </div>
        <span class="timestamp">{{ formattedTimestamp }}</span>
      </div>
      <div
        class="comment-content row"
        style="border: 2px solid aquamarine; border-radius: 2%; height: 30%; align-items: center"
      >
        <span v-if="comment.parentUser && isReply" class="reply-to">
          回复 @{{ comment.parentUser.UserName }}:
        </span>
        <div style="padding: 10px">
          {{ comment.content }}
        </div>
      </div>
      <div class="comment-actions" style="height: 10%">
        <button @click="emitReply">回复</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
const props = defineProps({
  comment: {
    type: Object,
    required: true,
  },
  isReply: {
    // 标记是否为子楼回复，用于可能的样式区分
    type: Boolean,
    default: false,
  },
  isClick: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['reply-to'])
const router = useRouter()
const toPersonal = () => {
  console.log(props.comment.userId)

  router.push({
    name: 'Personal',
    params: {
      userId: props.comment.userId,
    },
  })
}
const formattedTimestamp = computed(() => {
  return new Date(props.comment.createdAt).toLocaleString()
})

const emitReply = () => {
  emit('reply-to', props.comment.commentId, props.comment.userName) // 传递被回复评论的ID和用户名
}
</script>

<style scoped>
.comment-item {
  border: 1px solid #eee;
  padding: 10px;
  margin-bottom: 10px;
  border-radius: 4px;
  background-color: #fff;
}
.comment-item.is-reply {
  margin-left: 30px; /* 子楼评论缩进 */
  background-color: #f9f9f9;
  border-left: 3px solid #007bff;
}
.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}
.user-name {
  font-weight: bold;
  color: #333;
}
.timestamp {
  font-size: 0.8em;
  color: #777;
}
.comment-content {
  margin-bottom: 8px;
  word-wrap: break-word;
}
.reply-to {
  color: #007bff;
  margin-right: 5px;
}
.comment-actions button {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 0.9em;
}
.comment-actions button:hover {
  background-color: #0056b3;
}
</style>
