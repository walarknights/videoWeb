<template>
  <div class="q-pa-md row justify-center" style="margin-top: 100px">
    <div style="width: 60%" class="row justify-between">
      <div v-for="index in 4" :key="index" style="height: 20%; width: 40%">
        <div class="q-ma-sm col-4" style="height: 70%">
          <img
            src="https://cdn.quasar.dev/img/mountains.jpg"
            style="width: 100%; height: 100%"
            @click="toVedio(1)"
          />

          <div style="height: 20%">
            <div class="text-h6">Card #{{ index }}</div>
            <div class="text-subtitle2">by John Doe</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()

const toVedio = (id) => {
  router.push({ name: 'Video', params: { id } })
}
</script>
