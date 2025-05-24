<template>
  <div
    class="q-pa-md"
    style="display: flex; justify-content: center; width: 100%; margin-top: 2%; height: 200%"
  >
    <q-carousel
      animated
      v-model="currentSlide"
      navigation
      infinite
      :autoplay="isAutoplay"
      arrows
      transition-prev="slide-right"
      transition-next="slide-left"
      @mouseenter="isAutoplay = false"
      @mouseleave="isAutoplay = true"
      style="width: 60%"
    >
      <q-carousel-slide
        :name="1"
        img-src="https://cdn.quasar.dev/img/mountains.jpg"
        style="width: 100%"
      />
      <q-carousel-slide
        :name="2"
        img-src="https://cdn.quasar.dev/img/parallax1.jpg"
        style="width: 100%"
      />
      <q-carousel-slide
        :name="3"
        img-src="https://cdn.quasar.dev/img/parallax2.jpg"
        style="width: 100%"
      />
      <q-carousel-slide
        :name="4"
        img-src="https://cdn.quasar.dev/img/quasar.jpg"
        style="width: 100%"
      />
    </q-carousel>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  autoplay: {
    type: Boolean,
    default: true,
  },
  slide: {
    type: Number,
    default: 1,
  },
})

const isAutoplay = ref(props.autoplay)
const currentSlide = ref(props.slide)
</script>
<style>
.wid {
  width: 1000px;
}
</style>
