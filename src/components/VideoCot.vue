<template>
  <q-page padding style="margin-left: 10%; margin-top: 0%">
    <div class="row">
      <div class="left" style="width: 60%; margin-right: 10%">
        <div v-if="video">
          <h3 style="margin-top: 0%; margin-bottom: 4%">{{ video.title }}</h3>
          <q-chip icon="sym_o_motion_play">{{ formatView }}</q-chip>
          <q-chip icon="sym_o_history">{{ video.update_time }}</q-chip>
          <video
            ref="videoRef"
            @timeupdate="onTimeUpdate"
            :src="'http://localhost:8080' + video.url"
            controls
            style="width: 100%; height: 80%"
          ></video>

          <div class="row" style="justify-content: space-between">
            <q-btn
              @click="likeVideo"
              :class="islike ? 'materialIN' : 'materialOUT'"
              :label="video.likes"
              color="primary"
              icon="sym_o_thumb_up"
            />
            <p v-if="errorMsg" style="color: red; font-size: medium">{{ errorMsg }}</p>
            <q-btn
              @click="favoriteVideo"
              :class="isFavorite ? 'materialIN' : 'materialOUT'"
              :label="video.favorites"
              color="secondary"
              icon="sym_o_star_rate"
            />
          </div>

          <div>
            <ul v-if="isLoggedIn">
              <div>
                <CommentSection />
              </div>
            </ul>
            <p v-else>请登录后查看评论</p>
          </div>
        </div>
        <!-- 可以添加一个加载状态或错误状态 -->
        <div v-else>加载中或视频不存在...</div>
      </div>
      <div v-if="video" class="column">
        <img
          :src="'http://localhost:8080' + video.avatar"
          alt="..."
          style="border-radius: 50%; height: 75px; width: 75px"
        />
        <strong style="font-size: large; justify-items: center">{{ video.userName }}</strong>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useUserStore } from 'src/stores/useStore'
import api from 'src/utils/axios'
import { useRoute } from 'vue-router'
import CommentSection from './CommentSection.vue'
import { useQuasar } from 'quasar'
const q = useQuasar()
const video = ref(null)
const userStore = useUserStore()
const { isLoggedIn } = storeToRefs(userStore)
const userId = computed(() => userStore.userId)

// 添加格式化观看次数的计算属性
const formatView = computed(() => {
  if (!video.value) return '0'
  const views = video.value.view
  if (views >= 10000) {
    return (views / 10000).toFixed(1) + '万'
  }
  return views
})

const errorMsg = ref('')
const route = useRoute()
const videoRef = ref(null)
const progress = ref(0)
const flag = ref(false) // 将flag移到外部
const videoId = parseInt(route.params.id, 10)
async function onTimeUpdate() {
  const vd = videoRef.value
  if (vd && vd.duration > 0) {
    progress.value = Math.floor((vd.currentTime / vd.duration) * 100)
  }

  if (isLoggedIn.value && progress.value >= 30 && !flag.value) {
    try {
      const response = await api.post(`/videos/${videoId}/view`, {
        userId: parseInt(userId.value, 10),
      })
      if (response) {
        flag.value = true
        console.log('视频播放进度已记录')
      }
    } catch (error) {
      console.log('增加播放量失败', error)
    }
  }
}

// 将 videoId 转为 int 类型

const token = localStorage.getItem('authToken')
const isFavorite = ref(false)
const islike = ref(false)
onMounted(async () => {
  try {
    // 确保你的 API 地址和端口正确
    const res1 = await api.get(`/videos/${videoId}`) // 建议使用相对路径或配置 baseURL
    if (userId.value) {
      const res2 = await api.post(`/videos/${videoId}/islike`, {
        userId: parseInt(userId.value, 10),
      })
      if (res2.data.likes == 1) {
        islike.value = true
      } else {
        islike.value = false
      }

      const res3 = await api.post(`/videos/${videoId}/isfavorite`, {
        userId: parseInt(userId.value, 10),
      })
      if (res3.data.favorites == 1) {
        isFavorite.value = true
      } else {
        isFavorite.value = false
      }
    }
    video.value = res1.data
    console.log('视频数据:', video.value) // 添加日志方便调试
  } catch (error) {
    console.error('获取视频失败:', error)
  }
})
const likeVideo = async () => {
  if (!isLoggedIn.value) {
    errorMsg.value = '请先登录'
    return
  } else errorMsg.value = ''

  try {
    const response = await api.post(
      `/videos/${videoId}/like`,
      {
        userId: parseInt(userId.value, 10), // 转换为整数
      },
      {
        headers: {
          Authorization: `Bearer ${token}`, // 在这里添加 Authorization 请求头
        },
      },
    )

    if (response.data) {
      // 更新视频点赞数
      console.log(response.data.message)

      if (response.data.islike == 1) {
        video.value.likes += 1
        islike.value = true
        q.notify({
          message: '点赞成功',
          color: 'brand',
        })
      } else {
        video.value.likes -= 1
        islike.value = false
      }
    }
  } catch (error) {
    console.error('点赞失败', error)
    errorMsg.value = '点赞失败，请稍后重试'
  }
}

const favoriteVideo = async () => {
  if (!isLoggedIn.value) {
    errorMsg.value = '请先登录'
    return
  } else errorMsg.value = ''

  try {
    const response = await api.post(`/videos/${videoId}/favorite`, {
      userId: parseInt(userId.value, 10), // 转换为整数
    })
    if (response.data.isfavorites == 1) {
      // 更新视频收藏数
      video.value.favorites += 1
      isFavorite.value = true
      q.notify({
        message: '收藏成功',
        color: 'favorite',
      })
    } else {
      video.value.favorites -= 1
      isFavorite.value = false
    }
  } catch (error) {
    console.error('收藏失败', error)
    errorMsg.value = '收藏失败，请稍后重试'
  }
}
</script>

<style>
.materialOUT {
  font-variation-settings:
    'FILL' 0,
    'wght' 400,
    'GRAD' 0,
    'opsz' 24;
  border: none;
}

.materialIN {
  font-variation-settings:
    'FILL' 1,
    'wght' 400,
    'GRAD' 0,
    'opsz' 24;
  border: none;
}
.column {
  display: flex;
  align-items: center; /* 垂直居中 */
  gap: 10px; /* 图片与用户名之间的间距 */
}
</style>
