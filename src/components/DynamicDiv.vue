<template>
  <div class="column" style="width: 100%">
    <div style="margin: 0px; padding: 0px; width: 75%" class="row">
      <div style="width: 10%">
        <q-avatar round size="48px">
          <img src="https://cdn.quasar.dev/img/boy-avatar.png" style="width: 100%" />
        </q-avatar>
      </div>
      <div>
        <p>Mary</p>
        <p caption>3 min ago</p>
      </div>
    </div>
    <div class="contanier" style="margin-left: 7.5%">
      <p style="padding: 2% 0; margin: 0px">你好</p>
      <div class="image-container">
        <img :src="props.imageUrl" class="thumbnail" @click="openPreview" alt="Thumbnail" />

        <div v-if="isPreviewOpen" class="preview-overlay" @click.self="closePreview">
          <div class="preview-container">
            <div class="toolbar">
              <button @click="rotateLeft"><i class="fas fa-undo"></i> 向左旋转</button>
              <button @click="rotateRight"><i class="fas fa-redo"></i> 向右旋转</button>
              <button @click="resetRotation"><i class="fas fa-sync"></i> 重置</button>
              <button @click="openFullSize">
                <i class="fas fa-external-link-alt"></i> 查看原图
              </button>
              <button @click="closePreview"><i class="fas fa-times"></i> 关闭</button>
            </div>

            <!-- 预览图片 -->
            <div class="image-wrapper">
              <img
                :src="props.imageUrl"
                class="preview-image"
                :style="{ transform: `rotate(${rotation}deg)` }"
                alt="Preview"
              />
            </div>
          </div>
        </div>
      </div>
      <q-item class="row" style="padding: 0%">
        <q-btn
          icon="sym_o_thumb_up"
          @click="thumbUp"
          :class="isFilled ? 'materialIN' : 'materialOUT'"
        >
          <p style="margin: 0 10%">{{ count }}</p>
        </q-btn>
      </q-item>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  imageUrl: {
    type: String,
    required: true,
  },
})
const count = ref(0)
const isFilled = ref(false)
const isPreviewOpen = ref(false)
const rotation = ref(0)

const openPreview = () => {
  isPreviewOpen.value = true
}

const closePreview = () => {
  isPreviewOpen.value = false
  rotation.value = 0
}

const rotateLeft = () => {
  rotation.value -= 90
}

const rotateRight = () => {
  rotation.value += 90
}

const resetRotation = () => {
  rotation.value = 0
}

const openFullSize = () => {
  window.open(props.imageUrl, '_blank')
}

function thumbUp() {
  isFilled.value = !isFilled.value
  if (isFilled.value == true) {
    count.value++
  } else {
    count.value--
  }
}
</script>

<style scoped>
.image-container {
  position: relative;
  display: inline-block;
}

.thumbnail {
  cursor: pointer;
  max-width: 60%;
  transition: transform 0.3s;
}

.thumbnail:hover {
  transform: scale(1.05);
}

.preview-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.preview-container {
  position: relative;
  display: flex;
  flex-direction: column;
  max-width: 90%;
  max-height: 90%;
}

.toolbar {
  display: flex;
  justify-content: center;
  padding: 10px;
  background-color: rgba(0, 0, 0, 0.7);
  border-radius: 5px 5px 0 0;
  margin-bottom: 10px;
}

.toolbar button {
  margin: 0 5px;
  padding: 8px 12px;
  background-color: #444;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.toolbar button:hover {
  background-color: #666;
}

.image-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.preview-image {
  max-width: 100%;
  max-height: 80vh;
  transition: transform 0.3s ease;
}

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
</style>
