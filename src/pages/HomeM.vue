<template>
  <q-layout view="lHh lpr lFf" style="border: none" class="full-height">
    <q-page-container>
      <CasurolAuto v-model:slide="props.slide" :autoplay="props.autoplay" style="height: 100%" />
      <ContaninerM />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import CasurolAuto from 'src/components/CasurolAuto.vue'
import { reactive } from 'vue'
import ContaninerM from 'src/components/ContaninerM.vue'

const props = reactive({
  autoplay: true,
  slide: 1,
})
</script>
