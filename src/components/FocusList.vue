<template>
  <div class="q-pa-md flex justify-center">
    <div style="max-width: 90%; width: 300px">
      <q-intersection
        v-for="(focus, index) in focusList"
        :key="index"
        transition="flip-right"
        class="justify-center"
      >
        <q-item v-if="!isNull" :clickable="false" class="row" style="padding: 0%">
          <q-item-section avatar>
            <q-avatar text-color="white">
              <img :src="'http://localhost:8080' + focus.avatar" alt="..." />
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <q-item-label>{{ focus.userName }}</q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-btn
              @click="addFocus(focus.userId)"
              icon="sym_o_add"
              class="bg-focus text-white"
              style="align-items: center; justify-content: center"
            >
              关注</q-btn
            >
          </q-item-section>
        </q-item>
      </q-intersection>
    </div>
  </div>
  <p v-if="isNull" style="text-align: center">暂无关注</p>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import api from 'src/utils/axios'
import { useQuasar } from 'quasar'
import { useUserStore } from 'src/stores/useStore'
const useStore = useUserStore()
const q = useQuasar()
const route = useRoute()
const userId = route.params.userId
const focusList = ref()
const isNull = ref(false)
const albFocus = ref()
const getFocusList = async () => {
  try {
    const response = await api.get(`/personal/${userId}/focusList`)
    if (response.data) {
      focusList.value = response.data
      isNull.value = false
      console.log(isNull.value)
    } else {
      focusList.value = []
      isNull.value = true
      console.log(isNull.value)
    }
  } catch (error) {
    console.log('获取关注列表失败', error)
    isNull.value = true
    focusList.value = []
  }
}

const addFocus = async (focusedUserId) => {
  try {
    if (useStore.isLoggedIn) {
      const response = await api.post(`/personal/${userId}/addFocus`, {
        FocusId: useStore.userId,
        FocusedId: focusedUserId,
      })
      if (response) {
        albFocus.value = response.data.isFocus
        if (albFocus.value == 0) {
          q.notify({
            message: '不能关注自己',
            color: 'red',
          })
        } else if (albFocus.value == 1) {
          q.notify({
            message: '已经关注过了',
            color: 'blue-4',
          })
        } else {
          // 关注成功后增加following计数
          useStore.incrementFollowing()
          q.notify({
            message: '关注成功',
            color: 'green-10',
          })
        }
      }
    } else {
      q.notify({
        message: '请先登录',
        color: 'red',
      })
    }
  } catch (error) {
    console.log('关注失败', error)
  }
}
onMounted(getFocusList)
</script>
