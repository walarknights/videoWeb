<template>
  <q-page padding>
    <div class="row" style="justify-content: center">
      <DynamicDiv :imageUrl="image" style="width: 35%" />
    </div>
  </q-page>
</template>

<script setup>
import DynamicDiv from 'src/components/DynamicDiv.vue'
const image = ['https://cdn.quasar.dev/img/boy-avatar.png']
</script>
